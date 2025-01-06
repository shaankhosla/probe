package query

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestReformatSQL(t *testing.T) {

	filename := "test.csv"
	query := "select * from data"

	reformattedQuery := reformatSQL(query, filename)
	contains := strings.Contains(reformattedQuery, "select * from 'test.csv' limit")
	assert.Equal(t, true, contains, "Should contain")
}
