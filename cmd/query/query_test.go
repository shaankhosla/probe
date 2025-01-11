package query

import (
	"encoding/csv"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReformatSQL(t *testing.T) {

	filename := "test.csv"
	query := "select * from data"

	reformattedQuery := reformatSQL(query, filename)
	contains := strings.Contains(reformattedQuery, "select * from 'test.csv' limit")
	assert.Equal(t, true, contains, "")

	query = "select * from dataaa"
	reformattedQuery = reformatSQL(query, filename)
	contains = strings.Contains(reformattedQuery, "select * from 'test.csv'aa")
	assert.Equal(t, true, contains, "")

	query = "select * from data limit 3"
	reformattedQuery = reformatSQL(query, filename)
	contains = strings.Contains(reformattedQuery, "select * from 'test.csv' limit 3")
	assert.Equal(t, true, contains, "")
}

func TestExecuteQuery(t *testing.T) {
	InitializeDB()

	file, err := os.CreateTemp("", "temp-*.csv")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer func() {
		file.Close()
		os.Remove(file.Name())
	}()

	writer := csv.NewWriter(file)
	err = writer.WriteAll([][]string{
		{"a", "b", "c"},
		{"1", "2", "3"},
		{"3", "4", "5"},
	})
	if err != nil {
		t.Fatalf("Failed to write to CSV: %v", err)
	}
	writer.Flush()

	if err := writer.Error(); err != nil {
		t.Fatalf("Failed to flush CSV writer: %v", err)
	}

	results, err := ExecuteSQL("select sum(a) as a, max(c) as c from data", file.Name())
	if err != nil {
		t.Fatalf("Failed to query: %v", err)
	}
	data := [][]string{
		{"a", "c"},
		{"4", "5"},
	}
	assert.Equal(t, results, data, "Results didn't match query")
}
