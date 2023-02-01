package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pierre0210/wenku-api/internal/database/table/chapter"
)

var DB *sql.DB

func fatalError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func InitDatabase() {
	err := godotenv.Load()
	fatalError(err)

	dbPath, err := filepath.Abs(os.Getenv("DBPATH"))
	fatalError(err)

	DB, err = sql.Open("sqlite3", dbPath)
	fatalError(err)
	log.Println(fmt.Sprintf("%s connected.", dbPath))

	_, err = chapter.Chapter(DB)
	fatalError(err)
	log.Println("chapter_info table created.")
}
