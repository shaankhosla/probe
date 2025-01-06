package query

import (
	"fmt"
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

	filename := os.TempDir() + "temp.csv"

	fmt.Println(filename)
	// os.WriteFile(filename, "a, b, c\n1, 2, 3")
	// assert.Equal(t, false, true, "")
}
