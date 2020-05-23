package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

const (
	dsn string = "postgres://postgres:postgres@localhost:5432/postgres?connect_timeout=10&sslmode=disable&pool_max_conns=10"

	createSQL string = `
		CREATE TABLE IF NOT EXISTS parent(
			parent_id UUID PRIMARY KEY,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP,
			deleted BOOLEAN NOT NULL
		);
		CREATE TABLE IF NOT EXISTS child (
			child_id UUID PRIMARY KEY,
			FOREIGN KEY (child_id) REFERENCES parent (parent_id),
			name VARCHAR (255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			deleted_at TIMESTAMP,
			deleted BOOLEAN NOT NULL
		);`

	truncateSQL string = `
		TRUNCATE parent CASCADE;
	`

	insertParentSQL string = `
		INSERT INTO parent (
			parent_id,
			created_at,
			updated_at,
			deleted
		) VALUES ($1, $2, $3, $4);`

	insertChildSQL string = `
		INSERT INTO child (
			child_id,
			name,
			created_at,
			updated_at,
			deleted
		) VALUES ($1, $2, $3, $4, $5);`
)

func getUUIDv4() string {
	u4, err := uuid.NewV4()
	if err != nil {
		log.Fatalf("failed to generate UUID: %v", err)
	}
	return u4.String()
}

func insertChild(ctx context.Context, wg *sync.WaitGroup, p *pgxpool.Pool) {
	defer wg.Done()

	tx, err := p.Begin(ctx)
	if err != nil {
		log.Panicf("Unable to start transaction: %v\n", err)
	}

	log.Printf("Transaction started (%v)", tx)

	now := time.Now()

	parentID := getUUIDv4()

	if _, err := tx.Exec(ctx, insertParentSQL, parentID, now, now, false); err != nil {
		tx.Rollback(ctx)
		log.Panicf("Unable to create entity: %v\n", err)
	}

	if _, err := tx.Exec(ctx, insertChildSQL, parentID, fmt.Sprintf("name-%s", parentID), now, now, false); err != nil {
		tx.Rollback(ctx)
		log.Panicf("Unable to create entity: %v\n", err)
	}

	log.Printf("Transaction finished (%v)[entityID=%v, now=%v]", tx, parentID, now)

	tx.Commit(ctx)
}

func main() {
	ctx := context.Background()

	dbpool, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer dbpool.Close()

	var t string
	err = dbpool.QueryRow(ctx, "select 'TestOK!'").Scan(&t)
	if err != nil {
		log.Panicf("QueryRow failed: %v\n", err)
	}

	log.Printf("Connected [%v](%v)\n", dsn, t)

	for _, sql := range []string{createSQL, truncateSQL} {
		if _, err := dbpool.Exec(ctx, sql); err != nil {
			log.Panicf("Failed to create table: %v", err)
		}
	}

	wg := sync.WaitGroup{}

	for n := 0; n < 10; n++ {
		wg.Add(1)
		go insertChild(ctx, &wg, dbpool)
	}

	wg.Wait()
}
