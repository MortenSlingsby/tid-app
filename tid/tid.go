package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/urfave/cli/v2"
	"log"
	"os"
  "time"
	"path/filepath"
  "fmt"

)

func db_path() string {
  homeDir, _ := os.UserHomeDir()
  db_path := filepath.Join(homeDir, ".tid", "tid.db")
  return db_path
}

func main() {
	app := &cli.App{
		Usage: "App for tidsregistrering",
		Commands: []*cli.Command{
			{
        Name:    "add",
        Aliases: []string{"a"},
        Usage:   "Add a new AO to the list",
        Action: func(cCtx *cli.Context) error {
          sqliteDatabase, _ := sql.Open("sqlite3", db_path())
          insertAO := `INSERT INTO AO (code, name) VALUES (?, ?)`
          defer sqliteDatabase.Close()
          _ , err := sqliteDatabase.Exec(insertAO, cCtx.Args().First(), cCtx.String("full-name"))
          if err != nil {
            panic(err)
          }
          fmt.Println("Lagt til ny AO med code:", cCtx.Args().First(), "og beskrivelse: ", cCtx.String("full-name"))
          return nil
        },
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name:  "full-name",
                Usage:  "Full name of AO",
                Aliases: []string{"f"},
                Value: "",
                Required: false,
            },
        },
        Before: func(cCtx *cli.Context) error {
          if cCtx.Args().Len() != 1 {
            return fmt.Errorf("Expecting exactly 1 argument")
          }
          return nil
        },
      },
			{
        Name:    "start",
        Aliases: []string{"s"},
        Usage:   "Start time for an AO",
        Action: func(cCtx *cli.Context) error {
          sqliteDatabase, _ := sql.Open("sqlite3", db_path())
          now := time.Now().Format("2006-01-02 15:04:05")
          updateQuery := `
            UPDATE log
            SET active = 0, end_time = ?
            WHERE active = 1;
          `
          _, err := sqliteDatabase.Exec(updateQuery, now)
          if err != nil {
              log.Fatal(err)
          }

          insertAO := `INSERT INTO log (code, start_time, end_time, active) VALUES (?, ?, ?, ?)`
          exists, err := valueExists(sqliteDatabase, cCtx.Args().First()) 
          if !exists {
            log.Fatal("Code does not exist")
          }
          if err != nil {
            log.Fatal(err)
          }
          defer sqliteDatabase.Close()
          _ , err = sqliteDatabase.Exec(insertAO, cCtx.Args().First(), now, "", 1)
          if err != nil {
            panic(err)
          }
          fmt.Println("Starting log for:", cCtx.Args().First())
          return nil
        },
        Before: func(cCtx *cli.Context) error {
          if cCtx.Args().Len() != 1 {
            return fmt.Errorf("Expecting exactly 1 argument")
          }
          return nil
        },
      },
			{
        Name:    "stop",
        Usage:   "Stop log for today",
        Action: func(cCtx *cli.Context) error {
          sqliteDatabase, _ := sql.Open("sqlite3", db_path())
          now := time.Now().Format("2006-01-02 15:04:05")
          updateQuery := `
            UPDATE log
            SET active = 0, end_time = ?
            WHERE active = 1;
          `
          _, err := sqliteDatabase.Exec(updateQuery, now)
          if err != nil {
              log.Fatal(err)
          }
          defer sqliteDatabase.Close()
          fmt.Println("End log for today")
          return nil
        },
      },
			{
        Name:    "view",
        Usage:   "Show aggregated data",
        Action: func(cCtx *cli.Context) error {
          sqliteDatabase, _ := sql.Open("sqlite3", db_path())
          createTable(sqliteDatabase)
          defer sqliteDatabase.Close()

          return nil
        },
      },
			{
        Name:    "list",
        Usage:   "List all AOs in the table",
        Action: func(cCtx *cli.Context) error {
          sqliteDatabase, _ := sql.Open("sqlite3", db_path())
          showAO(sqliteDatabase)
          defer sqliteDatabase.Close()

          return nil
        },
      },
		},
	}
	if fileExists(db_path()) {
		os.MkdirAll(filepath.Dir(db_path()), 0700)
		create_database()
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func create_database() {
	sqliteDatabase, _ := sql.Open("sqlite3", db_path())
	defer sqliteDatabase.Close()
	initAO(sqliteDatabase)
	initLog(sqliteDatabase)
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return os.IsNotExist(err)
}

func initAO(db *sql.DB) {
	sqlAO := `CREATE TABLE IF NOT EXISTS AO (
		"code" TEXT NOT NULL PRIMARY KEY,		
		"name" TEXT
  );`

	_, err := db.Exec(sqlAO)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func initLog(db *sql.DB) {
	sqlLog := `CREATE TABLE IF NOT EXISTS log (
		"code" TEXT NOT NULL,		
		"start_time" TEXT,
    "end_time" TEXT,
    "active" INTEGER
  );`

	_, err := db.Exec(sqlLog)
	if err != nil {
		log.Fatal(err.Error())
	}
}

func valueExists(db *sql.DB, value string) (bool, error) {
    var exists bool

    query := `SELECT EXISTS(SELECT 1 FROM AO WHERE code = ? LIMIT 1)`
    err := db.QueryRow(query, value).Scan(&exists)
    if err != nil {
        return false, err
    }

    return exists, nil
}
// SELECT code, sum(strftime('%s', end_time) - strftime('%s', start_time)) from log group by code;
