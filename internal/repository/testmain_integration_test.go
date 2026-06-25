//go:build integration

package repository

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var testRepo *PostgresRepository
var testDB *sql.DB

func TestMain(m *testing.M) {
	ctx := context.Background()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:17-alpine",
			ExposedPorts: []string{"5432/tcp"},
			Env: map[string]string{
				"POSTGRES_USER":     "test",
				"POSTGRES_PASSWORD": "test",
				"POSTGRES_DB":       "testdb",
			},
			WaitingFor: wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		},
		Started: true,
	})
	if err != nil {
		log.Fatalf("start postgres container: %v", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		container.Terminate(ctx)
		log.Fatalf("get container host: %v", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		container.Terminate(ctx)
		log.Fatalf("get container port: %v", err)
	}

	dsn := fmt.Sprintf("postgres://test:test@%s:%s/testdb?sslmode=disable", host, port.Port())

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		container.Terminate(ctx)
		log.Fatalf("open db: %v", err)
	}

	if err := runMigrations(db); err != nil {
		db.Close()
		container.Terminate(ctx)
		log.Fatalf("run migrations: %v", err)
	}

	testDB = db
	testRepo = &PostgresRepository{QBQueries: New(db)}

	code := m.Run()

	db.Close()
	container.Terminate(ctx)
	os.Exit(code)
}

func runMigrations(db *sql.DB) error {
	entries, err := os.ReadDir("../../db/migrations")
	if err != nil {
		return fmt.Errorf("read migrations dir: %w", err)
	}

	var files []string
	for _, e := range entries {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".up.sql") {
			files = append(files, e.Name())
		}
	}
	sort.Strings(files)

	for _, name := range files {
		content, err := os.ReadFile("../../db/migrations/" + name)
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}
		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("exec %s: %w", name, err)
		}
	}

	return nil
}
