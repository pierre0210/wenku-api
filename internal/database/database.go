package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pierre0210/wenku-api/internal/database/table/chapter"
	"github.com/pierre0210/wenku-api/internal/database/table/user"
)

var DB *sql.DB

func fatalError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func InitDatabase() {
	dbPath, err := filepath.Abs(os.Getenv("DBPATH"))
	fatalError(err)

	DB, err = sql.Open("sqlite3", dbPath)
	fatalError(err)
	log.Println(fmt.Sprintf("%s connected.", dbPath))

	_, err = chapter.Chapter(DB)
	fatalError(err)
	log.Println("chapter_info table created.")
	_, err = user.User(DB)
	fatalError(err)
	log.Println("user_info table created.")
}
