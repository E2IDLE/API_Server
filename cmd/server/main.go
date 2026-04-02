package main

import (
	"API_Server/internal/config"
	"API_Server/internal/handler"
	"API_Server/internal/middleware"
	"API_Server/internal/repository"
	"API_Server/internal/router"
	"API_Server/internal/service"
	"API_Server/internal/ws"
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	sqliteDB := initSQLite()
	defer sqliteDB.Close()
	createTables(sqliteDB)

	cfg := config.Load()

	// ── DB 연결 ──
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("DB 연결 실패: %v", err)
	}

	if err := pool.Ping(context.Background()); err != nil {
		log.Fatalf("DB Ping 실패: %v", err)
	}

	log.Println("PostgreSQL 연결 성공")
	defer pool.Close()

	// ── Repository ──
	userRepo := repository.NewUserRepository(pool)
	tokenRepo := repository.NewTokenRepository(pool)
	agentRepo := repository.NewAgentRepository(pool)
	sessionRepo := repository.NewSessionRepository(pool)
	candidateRepo := repository.NewCandidateRepository(pool)

	_ = repository.NewChatLogRepository(sqliteDB)

	// ── Service ──
	authSvc := service.NewAuthService(userRepo, tokenRepo)
	userSvc := service.NewUserService(userRepo)
	agentSvc := service.NewAgentService(agentRepo)
	sessionSvc := service.NewSessionService(sessionRepo)
	candidateSvc := service.NewCandidateService(candidateRepo)
	turnSvc := service.NewTurnService(cfg)

	// ── WebSocket Hub ──
	hub := ws.NewHub()
	go hub.Run()

	// ── Handler ──
	authH := handler.NewAuthHandler(authSvc)
	userH := handler.NewUserHandler(userSvc)
	agentH := handler.NewAgentHandler(agentSvc)
	sessionH := handler.NewSessionHandler(sessionSvc, hub)
	candidateH := handler.NewCandidateHandler(candidateSvc, hub)
	turnH := handler.NewTurnHandler(turnSvc)
	wsH := handler.NewWSHandler(hub, tokenRepo)

	// ── Middleware ──
	authMw := middleware.NewAuthMiddleware(tokenRepo)

	// ── Router ──
	r := router.Setup(authH, userH, agentH, sessionH, candidateH, turnH, wsH, authMw)

	// ── HTTP 서버 ──
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: r,
	}

	go func() {
		log.Printf("서버 시작: :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("서버 오류: %v", err)
		}
	}()

	// ── Graceful Shutdown ──
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("서버 종료 중...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("서버 강제 종료: %v", err)
	}
	log.Println("서버 종료 완료")
}

func initSQLite() *sql.DB {
	_ = os.MkdirAll("data", os.ModePerm)

	db, err := sql.Open("sqlite3", "./data/app.db")
	if err != nil {
		log.Fatal("sqlite open error:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("sqlite ping error:", err)
	}

	log.Println("SQLite connected")
	return db
}

func createTables(db *sql.DB) {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS chat_logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			session_id TEXT NOT NULL,
			role TEXT NOT NULL,
			content TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);`,
		// 필요한 테이블 여기에 추가
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			log.Fatalf("테이블 생성 실패: %v", err)
		}
	}
	log.Println("SQLite 테이블 생성 완료")
}
