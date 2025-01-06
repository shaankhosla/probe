package query

import (
	"database/sql"
	"fmt"
	_ "github.com/marcboeker/go-duckdb"
	"log"
	"strings"
)

const (
	defaultRows    = 100
	insertionOrder = "SET preserve_insertion_order = false;\n"
)

var db *sql.DB

func InitializeDB() {
	var err error
	db, err = sql.Open("duckdb", "")
	if err != nil {
		log.Fatal(err)
	}

}

func reformatSQL(query string, filename string) string {
	fromStatement := fmt.Sprintf("from '%s'", filename)
	query = strings.Replace(query, "from data", fromStatement, 1)
	if !strings.Contains(query, "limit") {
		query += fmt.Sprintf(" limit %d", defaultRows)
	}
	query = insertionOrder + query
	return query

}

func GetAllColumns(filename string) string {
	query := "select * from data"
	query = reformatSQL(query, filename)

	rows, err := db.Query(query)
	if err != nil {
		log.Fatalf("Failed to execute query: %s", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		log.Fatalf("Failed to fetch column names: %s", err)
	}

	columnText := ""
	for _, col := range columns {
		columnText += col + "\n"
	}
	return columnText

}

func ExecuteSQL(query string, filename string) ([][]string, error) {
	query = reformatSQL(query, filename)

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch column names: %w", err)
	}

	var result [][]string
	result = append(result, columns)
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePointers := make([]interface{}, len(columns))
		for i := range values {
			valuePointers[i] = &values[i]
		}

		if err := rows.Scan(valuePointers...); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		row := make([]string, len(columns))
		for i, v := range values {
			if v != nil {
				row[i] = fmt.Sprintf("%v", v)
			} else {
				row[i] = ""
			}
		}

		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("row iteration error: %w", err)
	}

	return result, nil
}
