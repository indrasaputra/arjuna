package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"

	sdkpg "github.com/indrasaputra/arjuna/pkg/sdk/database/postgres"
	"github.com/indrasaputra/arjuna/service/user/entity"
	"github.com/indrasaputra/arjuna/service/user/internal/config"
	"github.com/uptrace/bun"
)

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig(".env")
	checkError(err)
	bunDB, err := sdkpg.NewDBWithPgx(cfg.Postgres)
	checkError(err)

	val := openJson("test/fixture/users.json")

	insertUsers(ctx, bunDB, val)
}

func openJson(file string) []byte {
	jsonFile, err := os.Open(file)
	checkError(err)
	defer func() {
		_ = jsonFile.Close()
	}()

	val, err := io.ReadAll(jsonFile)
	checkError(err)

	return val
}

func insertUsers(ctx context.Context, db *bun.DB, val []byte) {
	var users []*entity.User
	_ = json.Unmarshal(val, &users)

	query := "INSERT INTO users (id, name, created_at, updated_at) VALUES (?, ?, NOW(), NOW())"
	for _, user := range users {
		// password, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
		_, err := db.ExecContext(ctx, query, user.ID, user.Name)
		checkError(err)
	}
	log.Printf("Successfully insert %d users\n", len(users))
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
