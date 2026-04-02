package main

import (
	"API_Server/internal/config"
	"API_Server/internal/database"
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

	_ "modernc.org/sqlite"
)

func main() {
	// ── SQLite 연결 ──
	db := initSQLite()
	defer db.Close()
	database.RunMigrations(db)

	cfg := config.Load()

	// ── Repository ──
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	agentRepo := repository.NewAgentRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	candidateRepo := repository.NewCandidateRepository(db)

	// ── Service ──
	authSvc := service.NewAuthService(userRepo, tokenRepo)
	userSvc := service.NewUserService(userRepo)
	agentSvc := service.NewAgentService(agentRepo)
	sessionSvc := service.NewSessionService(sessionRepo, tokenRepo)
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
	wsH := handler.NewWSHandler(hub, tokenRepo, sessionSvc)

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

	db, err := sql.Open("sqlite", "./data/app.db")
	if err != nil {
		log.Fatal("sqlite open error:", err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal("sqlite ping error:", err)
	}

	db.SetMaxOpenConns(1) // SQLite 동시 쓰기 방지

	log.Println("SQLite connected")
	return db
}
