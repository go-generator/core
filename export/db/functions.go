package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
	"unicode"

	d "github.com/go-generator/core/driver"
	"github.com/go-generator/core/export/query"
	s "github.com/go-generator/core/export/sql"
)

var (
	tablesIndex map[string]int
)

func init() {
	var err error
	tablesIndex, err = s.GetColumnIndexes(reflect.TypeOf(Tables{}))
	if err != nil {
		panic(err)
	}
}

type Tables struct {
	Table string `gorm:"column:table"`
}

func ToLower(s string) string {
	if len(s) < 0 {
		return ""
	}
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}

func ListTables(ctx context.Context, db *sql.DB, database string) ([]string, error) {
	driver := s.GetDriver(db)
	var (
		tables []Tables
		res    []string
	)
	query, err := query.ListTablesQuery(database, database, driver)
	if err != nil {
		return nil, err
	}
	err = s.Query(ctx, db, tablesIndex, &tables, query)
	if err != nil {
		return nil, err
	}
	for i := range tables {
		res = append(res, tables[i].Table)
	}
	return res, err
}

func buildTableQuery(database, driver string) (string, error) {
	switch driver {
	case d.Mysql:
		query := `
		SELECT 
    		TABLE_NAME AS 'table'
		FROM
    		information_schema.tables
		WHERE
    		table_schema = '%v'`
		return fmt.Sprintf(query, database), nil
	case d.Postgres:
		return `
		SELECT 
    		table_name as table
		FROM
    		information_schema.tables
		WHERE
    		table_schema='public' AND table_type='BASE TABLE'`, nil
	default:
		return "", errors.New("unsupported driver")
	}
}

func ReformatGoName(s string) string {
	var field strings.Builder
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Println(err)
	}
	tokens := strings.Split(s, "_")
	for _, t := range tokens {
		alphanumericString := reg.ReplaceAllString(t, "")
		field.WriteString(strings.Title(alphanumericString))
	}
	return field.String()
}
