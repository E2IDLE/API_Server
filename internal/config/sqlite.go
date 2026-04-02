package config

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var SqliteDB *sql.DB

func InitSQLite(dbPath string) {
	var err error
	SqliteDB, err = sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("SQLite 연결 실패: %v", err)
	}

	// 연결 확인
	if err = SqliteDB.Ping(); err != nil {
		log.Fatalf("SQLite Ping 실패: %v", err)
	}

	// WAL 모드 활성화 (동시성 향상)
	_, err = SqliteDB.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		log.Fatalf("PRAGMA 설정 실패: %v", err)
	}

	log.Println("SQLite 연결 성공:", dbPath)
}
