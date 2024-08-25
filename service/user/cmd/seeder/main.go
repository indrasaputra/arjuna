// Main seeder program to insert users data into the database.
package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/builder"
	"github.com/indrasaputra/arjuna/service/user/internal/config"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig(".env")
	checkError(err)
	db, err := builder.BuildBunDB(cfg.Postgres)
	checkError(err)

	val := openJSON("test/fixture/users.json")

	insertUsers(ctx, db, val)
}

func openJSON(file string) []byte {
	jsonFile, err := os.Open(filepath.Clean(file))
	checkError(err)
	defer func() {
		_ = jsonFile.Close()
	}()

	val, err := io.ReadAll(jsonFile)
	checkError(err)

	return val
}

func insertUsers(ctx context.Context, db *postgres.BunDB, val []byte) {
	var users []*entity.User
	_ = json.Unmarshal(val, &users)

	query := "INSERT INTO users (id, name, created_at, updated_at) VALUES (?, ?, NOW(), NOW())"
	for _, user := range users {
		_, err := db.Exec(ctx, query, user.ID, user.Name)
		checkError(err)
	}
	log.Printf("Successfully insert %d users\n", len(users))
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
