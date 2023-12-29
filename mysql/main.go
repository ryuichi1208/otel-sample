package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"go.nhat.io/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func add(ctx context.Context, a, b int) int {
	_, span := tracer.Start(ctx, "add")
	defer span.End()
	return a + b
}

func Do(ctx context.Context) {
	ctx, span := tracer.Start(ctx, "Do")
	defer span.End()
	add(ctx, 1, 2)
	time.Sleep(1 * time.Second)
}

func Run(ctx context.Context) {
	// データベース接続の設定
	dsn := "root:mysql@tcp(localhost:3306)/db"
	db, err := openDB(dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// データベース接続の確認
	err = db.PingContext(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// SELECTクエリの実行
	var (
		version string
	)
	rows, err := db.QueryContext(ctx, "SELECT version()")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&version)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(version)
	}

	// エラー処理
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("sleep")
	time.Sleep(5 * time.Second)
}

func f(ctx context.Context) {
	_, span := tracer.Start(ctx, "f")
	defer span.End()
	Do(ctx)
	Run(ctx)
	time.Sleep(1 * time.Second)
}

func main() {
	ctx := context.Background()
	shutdown := initProvider(ctx)
	defer shutdown()
	ctx, span := tracer.Start(ctx, "main")
	f(ctx)
	span.End()
	time.Sleep(3 * time.Second)
}

func openDB(dsn string) (*sql.DB, error) {
	// Register the otelsql wrapper for the provided postgres driver.
	driverName, err := otelsql.Register("mysql",
		otelsql.AllowRoot(),
		otelsql.TraceQueryWithoutArgs(),
		otelsql.TraceRowsClose(),
		otelsql.TraceRowsAffected(),
		otelsql.WithDatabaseName("my_database"),        // Optional.
		otelsql.WithSystem(semconv.DBSystemPostgreSQL), // Optional.
	)
	if err != nil {
		return nil, err
	}

	// Connect to a Postgres database using the postgres driver wrapper.
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}

	if err := otelsql.RecordStats(db); err != nil {
		return nil, err
	}

	return db, nil
}
