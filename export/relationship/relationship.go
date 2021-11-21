package relationship

import (
	"context"
	"database/sql"
	"reflect"
	"regexp"
	"strings"

	s "github.com/core-go/sql"
	"github.com/go-generator/core/export/query"
)

var (
	sqliteRelIndex, relTableIndex map[string]int
)

func init() {
	var errIndex error
	sqliteRelIndex, errIndex = s.GetColumnIndexes(reflect.TypeOf(SqliteRel{}))
	if errIndex != nil {
		panic(errIndex)
	}
	relTableIndex, errIndex = s.GetColumnIndexes(reflect.TypeOf(RelTables{}))
	if errIndex != nil {
		panic(errIndex)
	}
}

func pop(slice []string) []string {
	return append(slice[:0], slice[1:]...)
}

// GetRelationshipTable
// 1-1 -> both fields are unique
// 1-n -> only one field is unique
// n-n -> both fields are not unique
// self reference will be in the same table with the same datatype
func GetRelationshipTable(ctx context.Context, db *sql.DB, database string, tables []string, primaryKeys map[string][]string) ([]RelTables, error) {
	driver := s.GetDriver(db)
	var relations []RelTables
	var sqliteRels []SqliteRel
	switch s.GetDriver(db) {
	case s.DriverSqlite3:
		listReferenceQuery, err := query.ListReferenceQuery(database, driver, "")
		if err != nil {
			return nil, err
		}
		err = s.Query(ctx, db, sqliteRelIndex, &sqliteRels, listReferenceQuery)
		if err != nil {
			return nil, err
		}

		tb := regexp.MustCompile(`(?s)\".*?\"`)
		fk := regexp.MustCompile(`(?s)\(\[.*?\]\)`)
		for i := range sqliteRels {
			tables := tb.FindAllString(sqliteRels[i].Sql, -1)
			for i := range tables {
				tables[i] = strings.ReplaceAll(tables[i], `"`, ``)
			}
			var tbNames []string
			parentTb := tables[0]
			tbNames = append(tbNames, parentTb)
			tables = pop(tables)
			for i := range tables {
				if i == len(tables)-1 {
					tbNames = append(tbNames, tables[i])
				} else {
					tbNames = append(tbNames, parentTb)
					tbNames = append(tbNames, tables[i])
				}
			}
			columns := fk.FindAllString(sqliteRels[i].Sql, -1)
			for i := range columns {
				columns[i] = strings.ReplaceAll(columns[i], `([`, ``)
				columns[i] = strings.ReplaceAll(columns[i], `])`, ``)
			}
			if len(tables) > 0 {
				for i := 0; i < len(tbNames)-1; i++ {
					var rel RelTables
					rel.Table = tbNames[i]
					rel.ReferencedTable = tbNames[i+1]
					rel.Column = columns[i]
					rel.ReferencedColumn = columns[i+1]
					relations = append(relations, rel)
				}
			}
		}
	case s.DriverOracle:
		var relTable []RelTables
		for i := range tables {
			q, err := query.ListReferenceQuery(database, driver, tables[i])
			if err != nil {
				return nil, err
			}
			err = s.Query(ctx, db, relTableIndex, &relTable, q)
			if err != nil {
				return nil, err
			}
			relations = append(relations, relTable...)
		}
	default:
		listReferenceQuery, err := query.ListReferenceQuery(database, driver, "")
		if err != nil {
			return nil, err
		}
		err = s.Query(ctx, db, relTableIndex, &relations, listReferenceQuery)
		if err != nil {
			return nil, err
		}
	}
	for i := range relations {
		isP1 := IsPrimaryKey(relations[i].Column, relations[i].Table, primaryKeys)
		isP2 := IsPrimaryKey(relations[i].ReferencedColumn, relations[i].ReferencedTable, primaryKeys)
		if isP1 && isP2 {
			if len(primaryKeys[relations[i].Table]) != len(primaryKeys[relations[i].ReferencedTable]) {
				relations[i].Relationship = OneToMany
			} else {
				relations[i].Relationship = OneToOne
			}
		} else {
			relations[i].Relationship = OneToMany
		}
	}
	return relations, nil
} // Find all columns, table and its referenced columns, tables

func IsPrimaryKey(key, table string, pks map[string][]string) bool {
	for i := range pks[table] {
		if key == pks[table][i] {
			return true
		}
	}
	return false
}
