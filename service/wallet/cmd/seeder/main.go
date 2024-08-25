// Main seeder program to insert wallets data into the database.
package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/service/wallet/entity"
	"github.com/indrasaputra/arjuna/service/wallet/internal/builder"
	"github.com/indrasaputra/arjuna/service/wallet/internal/config"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig(".env")
	checkError(err)
	db, err := builder.BuildBunDB(cfg.Postgres)
	checkError(err)

	val := openJSON("test/fixture/wallets.json")

	insertWallets(ctx, db, val)
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

func insertWallets(ctx context.Context, db *postgres.BunDB, val []byte) {
	var wallets []*entity.Wallet
	_ = json.Unmarshal(val, &wallets)

	query := "INSERT INTO wallets (id, user_id, balance, created_at, updated_at) VALUES (?, ?, ?, NOW(), NOW())"
	for _, wallet := range wallets {
		_, err := db.Exec(ctx, query, wallet.ID, wallet.UserID, wallet.Balance)
		checkError(err)
	}
	log.Printf("Successfully insert %d wallets\n", len(wallets))
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
