package relationship

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"regexp"
	"strings"

	d "github.com/go-generator/core/driver"
	"github.com/go-generator/core/export/query"
	s "github.com/go-generator/core/export/sql"
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

type TableRelation struct {
	Table    string
	RefTable string
}

type ColumnRelation struct {
	Col    string
	RefCol string
}

// GetRelationshipTable
// 1-1 -> both fields are unique
// 1-n -> only one field is unique
// n-n -> both fields are not unique
// self reference will be in the same table with the same datatype
func GetRelationshipTable(ctx context.Context, db *sql.DB, database string, tables []string, primaryKeys map[string][]string) ([]RelTables, error) {
	driver := s.GetDriver(db)
	var (
		relations      []RelTables
		sqliteRels     []SqliteRel
		tableRelations []TableRelation
		colRelations   []ColumnRelation
	)
	switch s.GetDriver(db) {
	case d.Sqlite3:
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
			if len(tables) < 1 {
				return nil, errors.New("error getting table relation")
			}
			parentTb := tables[0]
			for i := 1; i < len(tables); i++ {
				tableRelations = append(tableRelations, TableRelation{
					Table:    parentTb,
					RefTable: tables[i],
				})
			}
			columns := fk.FindAllString(sqliteRels[i].Sql, -1)
			for i := range columns {
				columns[i] = strings.ReplaceAll(columns[i], `([`, ``)
				columns[i] = strings.ReplaceAll(columns[i], `])`, ``)
			}
			for i := 0; i < len(columns)-1; i += 2 {
				colRelations = append(colRelations, ColumnRelation{
					Col:    columns[i],
					RefCol: columns[i+1],
				})
			}
			for i := 0; i < len(tableRelations); i++ {
				var rel RelTables
				rel.Table = tableRelations[i].RefTable
				rel.ReferencedTable = tableRelations[i].Table
				rel.Column = colRelations[i].RefCol
				rel.ReferencedColumn = colRelations[i].Col
				relations = append(relations, rel)
			}
		}
	case d.Oracle:
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
		}
		if isP1 && !isP2 {
			relations[i].Relationship = OneToMany
		}
		if !isP1 && isP2 {
			relations[i].Relationship = ManyToOne
		}
	}
	for i := range relations {
		if relations[i].Relationship == OneToMany {
			skip := false
			reverse := RelTables{
				Table:            relations[i].ReferencedTable,
				Column:           relations[i].ReferencedColumn,
				ReferencedTable:  relations[i].Table,
				ReferencedColumn: relations[i].Column,
				Relationship:     ManyToOne,
			}
			for j := range relations {
				if reverse == relations[j] {
					skip = true
					break
				}
			} // skip duplicate
			if !skip {
				relations = append(relations, reverse)
			}
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
