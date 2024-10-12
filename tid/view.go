package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	 "time"
	"github.com/jedib0t/go-pretty/v6/table"
	// "github.com/jedib0t/go-pretty/v6/text"
)

type AO struct {
  tag string
  description string
}

type Log struct {
  tag string
  description string
  startTime string
  endTime string
  duration string
  active string
}

func createTable(db *sql.DB, week int) {
  codes := get_codes(db, week)
  t := table.NewWriter()
  for _, code := range codes {
    row := calcRow(db, code, week)
    t.AppendRow(row)
  }
  t.AppendHeader(calcHeader(week))
  fmt.Println(t.Render())
}

func get_first_day(week int) time.Time {
  now := time.Now()
  weekDay := int(now.Weekday())
  return now.AddDate(0, 0, -weekDay+week)
}

func secondString(seconds int) string {
    hours := seconds / 3600
    minutes := (seconds % 3600) / 60
    return fmt.Sprintf("%02d:%02d", hours, minutes)
}

func calcHeader(week int) []interface{} {
  row := make([]interface{}, 8)
  row[0] = "AO"
  date := get_first_day(week)
  for i := 0; i < 7; i++ {
    row[i+1] = fmt.Sprintf("%s\n%s", date.Format("02.01"), date.Weekday().String())
    date = date.AddDate(0,0,1)
  }
  return row
}

func calcRow(db *sql.DB, code string, week int) []interface{} {
  row := make([]interface{}, 8)
  row[0] = fullName(db, code)
  date := get_first_day(week)
  for i := 0; i < 7; i++ {
    row[i+1] = secondString(calcVal(db, code, date.Format("2006-01-02")))
    date = date.AddDate(0,0,1)
  }
  return row
}

func calcVal(db *sql.DB, code string, date string) int {
  query := "SELECT sum(strftime('%s', end_time) - strftime('%s', start_time)) FROM log WHERE code = ? AND DATE(start_time) = ?"

  var result sql.NullInt64
  err := db.QueryRow(query, code, date).Scan(&result)
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

func get_codes(db *sql.DB, week int) []string {
  var result []string
  first_date := get_first_day(week)
  last_date := first_date.AddDate(0,0,7)
  query := "SELECT DISTINCT code FROM log WHERE start_time > ? and end_time < ?"
  rows, err := db.Query(query, first_date.Format("2006-01-02"), last_date.Format("2006-01-02"))
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
      l.code,
      a.name,
      l.start_time,
      l.end_time,
      printf('%02d:%02d', l.duration / 3600, (l.duration % 3600) / 60) as duration,
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
    if err := rows.Scan(&log.tag, &log.description, &log.startTime, &log.endTime, &log.duration, &log.active); err != nil {
        panic(err)
    }
    t.AppendRow(table.Row{log.tag, log.description, log.startTime, log.endTime, log.duration, log.active})
  }
  t.AppendHeader(table.Row{"Tag", "Description", "Start Time", "End time", "Duration", "Active"})
  fmt.Println(t.Render())
}
