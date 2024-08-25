// Main seeder program to insert accounts data into the database.
package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"

	"golang.org/x/crypto/bcrypt"

	"github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/service/auth/entity"
	"github.com/indrasaputra/arjuna/service/auth/internal/builder"
	"github.com/indrasaputra/arjuna/service/auth/internal/config"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig(".env")
	checkError(err)
	db, err := builder.BuildBunDB(cfg.Postgres)
	checkError(err)

	val := openJSON("test/fixture/accounts.json")

	insertAccounts(ctx, db, val)
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

func insertAccounts(ctx context.Context, db *postgres.BunDB, val []byte) {
	var accounts []*entity.Account
	_ = json.Unmarshal(val, &accounts)

	query := "INSERT INTO accounts (id, user_id, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, NOW(), NOW())"
	for _, account := range accounts {
		password, _ := bcrypt.GenerateFromPassword([]byte(account.Password), bcrypt.MinCost)
		_, err := db.Exec(ctx, query, account.ID, account.UserID, account.Email, string(password))
		checkError(err)
	}
	log.Printf("Successfully insert %d accounts\n", len(accounts))
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
