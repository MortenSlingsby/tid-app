package main

import (
	"database/sql"
	"fmt"

	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	_ "github.com/mattn/go-sqlite3"
)

type AO struct {
	tag         string
	description string
}

type Log struct {
	id          int
	tag         string
	description string
	startTime   string
	endTime     string
	duration    string
	active      string
}

func createTable(db *sql.DB, week int) {
	codes := getCodes(db, week)
	t := table.NewWriter()
	for _, code := range codes {
		row := calcRow(db, code, week)
		t.AppendRow(row)
	}
	t.AppendHeader(calcHeader(week))
	t.AppendFooter(calcFooter(db, week))
	fmt.Println(t.Render())
}

func getFirstDay(week int) time.Time {
	now := time.Now()
	weekDay := int(now.Weekday())
	return now.AddDate(0, 0, -weekDay+week)
}

func secondString(seconds int) string {
	hours := float64(seconds) / 3600

	// hours := seconds / 3600
	// minutes := (seconds % 3600) / 60

	return fmt.Sprintf("%.2f", hours)
}

func calcHeader(week int) []any {
	row := make([]any, 8)
	row[0] = "AO"
	date := getFirstDay(week)
	for i := range 7 {
		row[i+1] = fmt.Sprintf("%s\n%s", date.Format("02.01"), date.Weekday().String())
		date = date.AddDate(0, 0, 1)
	}
	return row
}

func calcFooter(db *sql.DB, week int) []any {
	row := make([]any, 8)
	row[0] = "Total"
	date := getFirstDay(week)
	for i := range 7 {
		row[i+1] = secondString(calcVal(db, "", date.Format("2006-01-02"), true))
		date = date.AddDate(0, 0, 1)
	}
	return row
}

func calcRow(db *sql.DB, code string, week int) []any {
	row := make([]any, 8)
	row[0] = fullName(db, code)
	date := getFirstDay(week)
	for i := range 7 {
		row[i+1] = secondString(calcVal(db, code, date.Format("2006-01-02"), false))
		date = date.AddDate(0, 0, 1)
	}
	return row
}

func calcVal(db *sql.DB, code string, date string, total bool) int {
	var query string
	var result sql.NullInt64
	var err error
	if total {
		query = "SELECT sum(duration) FROM log WHERE DATE(start_time) = ?"
		err = db.QueryRow(query, date).Scan(&result)
	} else {
		query = "SELECT sum(duration) FROM log WHERE code = ? AND DATE(start_time) = ?"
		err = db.QueryRow(query, code, date).Scan(&result)
	}

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No rows were returned!")
			return 0
		} else {
			panic(err)
		}
	}
	if result.Valid {
		return int(result.Int64)
	}
	return 0
}

func fullName(db *sql.DB, code string) string {
	query := "SELECT name FROM AO WHERE code = ?"
	var result string
	err := db.QueryRow(query, code).Scan(&result)
	if err != nil {
		panic(err)
	}
	return result
}

func getCodes(db *sql.DB, week int) []string {
	var result []string
	firstDate := getFirstDay(week)
	lastDate := firstDate.AddDate(0, 0, 7)
	query := "SELECT DISTINCT code FROM log WHERE start_time > ? and (end_time < ? OR end_time IS 'fix')"
	rows, err := db.Query(query, firstDate.Format("2006-01-02"), lastDate.Format("2006-01-02"))
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var value string
		if err := rows.Scan(&value); err != nil {
			panic(err)
		}
		result = append(result, value)
	}
	return result
}

func showAO(db *sql.DB) {
	rows, err := db.Query("SELECT * FROM AO")
	if err != nil {
		panic(err)
	}
	t := table.NewWriter()
	for rows.Next() {
		var aO AO
		if err := rows.Scan(&aO.tag, &aO.description); err != nil {
			panic(err)
		}
		t.AppendRow(table.Row{aO.tag, aO.description})
	}
	t.AppendHeader(table.Row{"Tag", "Description"})
	fmt.Println(t.Render())
}

func showLog(db *sql.DB) {
	now := time.Now().Format("2006-01-02")
	rows, err := db.Query(`
    select
      l.id,
      l.code,
      a.name,
      l.start_time,
      l.end_time,
      CASE
        WHEN l.active = 0 THEN
          printf('%s%02d:%02d', (CASE WHEN l.duration < 0 THEN '-' ELSE '' END), abs(l.duration) / 3600, (abs(l.duration) % 3600) / 60)
        ELSE
          printf('%02d:%02d', (strftime('%s', 'now', 'utc') - strftime('%s', l.start_time, 'utc')) / 3600, ((strftime('%s', 'now', 'utc') - strftime('%s', l.start_time, 'utc')) % 3600) / 60)
      END as duration,
      l.active
    from log as l
    inner join AO as a ON l.code = a.code
    WHERE start_time > ?
  `, now)
	if err != nil {
		panic(err)
	}
	t := table.NewWriter()
	for rows.Next() {
		var log Log
		if err := rows.Scan(&log.id, &log.tag, &log.description, &log.startTime, &log.endTime, &log.duration, &log.active); err != nil {
			panic(err)
		}
		t.AppendRow(table.Row{log.id, log.tag, log.description, log.startTime, log.endTime, log.duration, log.active})
	}
	t.AppendHeader(table.Row{"ID", "Tag", "Description", "Start Time", "End time", "Duration", "Active"})
	fmt.Println(t.Render())
}

func todayTotalOut(db *sql.DB) {
	total := secondString(todayTotalNonActive(db) + todayTotalActive(db))
	expected := secondString(todayTotal(db))
	out := fmt.Sprintf("You have logged %s out of a possible %s today", total, expected)
	fmt.Println(out)
}

func todayTotalNonActive(db *sql.DB) int {
	var sum int
	now := time.Now().Format("2006-01-02")
	err := db.QueryRow(`
    select
      sum(duration)
    FROM log
    WHERE start_time > ? AND active = 0
  `, now).Scan(&sum)
	if err != nil {
		panic(err)
	}

	return sum
}

func todayTotalActive(db *sql.DB) int {
	var sum sql.NullInt64
	now := time.Now().Format("2006-01-02")
	err := db.QueryRow(`
    select
        sum(strftime('%s', 'now', 'utc') - strftime('%s', start_time, 'utc'))
    FROM log
    WHERE start_time > ? AND active = 1
  `, now).Scan(&sum)
	if err != nil {
		panic(err)
	}

	if sum.Valid {
		return int(sum.Int64)
	}
	return 0
}

func todayTotal(db *sql.DB) int {
	var sum sql.NullInt64
	now := time.Now().Format("2006-01-02")
	err := db.QueryRow(`
    select
        strftime('%s', 'now', 'utc') - strftime('%s', start_time, 'utc')
    FROM log
    WHERE start_time > ?
    ORDER BY start_time ASC
    LIMIT 1
  `, now).Scan(&sum)
	if err != nil {
		panic(err)
	}
	if sum.Valid {
		return int(sum.Int64)
	}
	return 0
}

func dropLog(db *sql.DB, id int) string {
	query := "DELETE FROM log WHERE id = ?"
	var result string
	_, err := db.Exec(query, id)
	if err != nil {
		panic(err)
	}
	return result
}

func dropAO(db *sql.DB, id string) string {
	query := "DELETE FROM AO WHERE code = ?"
	var result string
	_, err := db.Exec(query, id)
	if err != nil {
		panic(err)
	}
	return result
}
