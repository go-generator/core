package relationship

import (
	"context"
	"database/sql"
	"errors"
	s "github.com/core-go/sql"
	"github.com/go-generator/core/export/query"
	"regexp"
	"strings"
)

// FindRelationships
// 1-1 -> both fields are unique
// 1-n -> only one field is unique
// n-n -> both fields are not unique
// self reference will be in the same table with the same datatype
func FindRelationships(ctx context.Context, db *sql.DB, database string) ([]RelTables, []string, error) {
	relTables, err := initRelationshipTable(ctx, db, database)
	if err != nil {
		return nil, nil, err
	}
	joinTables1, err := listJoinTables(ctx, db, database, relTables)
	for i := range relTables {
		relTables[i].Relationship, err = findRelationShip(ctx, db, database, joinTables1, &relTables[i])
		if err != nil {
			return nil, nil, err
		}
	}
	joinTables2, err := listJoinTables(ctx, db, database, relTables)
	if err != nil {
		return nil, nil, err
	}
	return relTables, joinTables2, err
}

func setRelation(ctx context.Context, db *sql.DB, unique, refUnique bool, database, driver string, rt *RelTables) (string, error) {
	// Already cover the ManyToMany case where a joined table consists of two or more primary key tags that are all foreign keys
	isPrimaryTag, err := CheckPrimaryTag(ctx, db, database, driver, rt.Table, rt.Column)
	if err != nil {
		return "", err
	}
	isReferencedPrimaryTag, err := CheckPrimaryTag(ctx, db, database, driver, rt.ReferencedTable, rt.ReferencedColumn)
	if err != nil {
		return "", err
	}
	//if !refUnique {
	//	return Unsupported, err
	//}
	if unique {
		var keys []PrimaryKey
		query, err := query.ListAllPrimaryKeys(database, driver, rt.Table)
		err = s.Query(ctx, db, nil, &keys, query)
		if err != nil {
			return "", err
		}
		if len(keys) == 1 { // Only one column has Primary Tag
			if isPrimaryTag && isReferencedPrimaryTag { // Both are Primary key
				return OneToOne, err
			}
			if !isPrimaryTag && isReferencedPrimaryTag { // Column is only a foreign key referenced to other primary key
				return ManyToOne, err
			}
		}
		if len(keys) > 1 { // Consist of at least one column that has primary key tag and not referenced to other table
			return OneToMany, err
		}
	}
	if !unique {
		return ManyToOne, err
	}
	if !unique && !refUnique {
		return Unsupported, err
	}
	return Unknown, err
}

func findRelationShip(ctx context.Context, db *sql.DB, database string, joinedTable []string, rt *RelTables) (string, error) { //TODO: switch gorm to core-go/sql
	driver := s.GetDriver(db)
	unique, err := isUnique(ctx, db, database, driver, rt.Table, rt.Column)
	if err != nil {
		return "", err
	}
	refUnique, err := isUnique(ctx, db, database, driver, rt.ReferencedTable, rt.ReferencedColumn)
	if err != nil {
		return "", err
	}
	for _, v := range joinedTable {
		if rt.Table == v {
			return ManyToMany, err
		}
	}
	return setRelation(ctx, db, unique, refUnique, database, driver, rt)
}

func isForeignKey(table, column string, relTables []RelTables) bool {
	for _, v := range relTables {
		if v.Table == table && v.Column == column {
			return true
		}
	}
	return false
} // Check if the column of the table is a foreign key

func isJoinTable(table string, columns []string, rt []RelTables) bool {
	for _, v := range columns {
		if isForeignKey(table, v, rt) == false {
			return false
		}
	}
	return true
}

func initRelationshipTable(ctx context.Context, db *sql.DB, database string) ([]RelTables, error) {
	driver := s.GetDriver(db)
	var relTables []RelTables
	var sqliteRels []SqliteRel
	query, err := query.ListReferenceQuery(database, driver, "")
	if err != nil {
		return nil, err
	}
	switch s.GetDriver(db) {
	case s.DriverSqlite3:
		err = s.Query(ctx, db, nil, &sqliteRels, query)
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
					relTables = append(relTables, rel)
				}
			}
		}
		return relTables, err
	default:
		err = s.Query(ctx, db, nil, &relTables, query)
		if err != nil {
			return nil, err
		}
		return relTables, err
	}
} // Find all columns, table and its referenced columns, tables

func listJoinTables(ctx context.Context, db *sql.DB, database string, rt []RelTables) ([]string, error) {
	var joinTable []string
	tables, err := ListTables(ctx, db, database)
	if err != nil {
		return nil, err
	}
	for _, v := range tables {
		primaryKeys, err := listPrimaryKeys(ctx, db, database, v.Name)
		if err != nil {
			return nil, err
		}
		var res []string
		for _, v := range primaryKeys {
			res = append(res, v.Column)
		}
		if len(primaryKeys) > 1 && isJoinTable(v.Name, res, rt) {
			joinTable = append(joinTable, v.Name)
		}
	}
	return joinTable, err
}

func pop(slice []string) []string {
	return append(slice[:0], slice[1:]...)
}

func isUnique(ctx context.Context, db *sql.DB, database, driver, table, column string) (bool, error) {
	var (
		mySqlIndex  []MySqlUnique
		pgIndex     []PostgresUnique
		mssqlIndex  []MssqlUnique
		sqliteIndex []SqliteUnique
	)
	query, err := query.ListUniqueQuery(database, driver, table)
	if err != nil {
		return false, err
	}
	switch driver {
	case s.DriverMysql:
		err = s.Query(ctx, db, nil, &mySqlIndex, query)
		if err != nil {
			return false, err
		}
		for _, v := range mySqlIndex {
			if v.Column == column {
				if v.NonUnique == false {
					return true, err
				}
			}
		}
	case s.DriverPostgres:
		err = s.Query(ctx, db, nil, &pgIndex, query)
		if err != nil {
			return false, err
		}
		for _, v := range pgIndex {
			if strings.Contains(v.Index, "unq") {
				tokens := strings.Split(v.Index, "_")
				for i := range tokens {
					if tokens[i] == "unq" {
						columnName := strings.Join(tokens[i:], "_")
						if column == columnName {
							return true, err
						}
					}
				}
			}
		}
	case s.DriverMssql:
		err = s.Query(ctx, db, nil, &mssqlIndex, query)
		if err != nil {
			return false, err
		}
		for _, v := range mssqlIndex {
			if v.Column == column {
				if strings.Contains(v.Constraint, "UQ") && column == v.Column {
					return true, err
				}
			}
		}
	case s.DriverSqlite3:
		err = s.Query(ctx, db, nil, &sqliteIndex, query)
		if err != nil {
			return false, err
		}
		for _, v := range sqliteIndex {
			if strings.Contains(v.Name, column) {
				if v.Unique == "1" {
					return true, err
				}
			}
		}
	default:
		return false, errors.New(s.DriverNotSupport)
	}
	return false, err
} // Check if a column is unique

func CheckPrimaryTag(ctx context.Context, db *sql.DB, database, driver, table, column string) (bool, error) {
	//TODO: Add check primary tag for other relationship
	var (
		mySqlIndex    []MySqlUnique
		mssqlIndex    []MssqlUnique
		postgresIndex []PostgresUnique
		sqliteIndex   []SqliteUnique
	)
	query, err := query.ListUniqueQuery(database, driver, table)
	if err != nil {
		return false, err
	}
	switch driver {
	case s.DriverMysql:
		err := s.Query(ctx, db, nil, &mySqlIndex, query)
		if err != nil {
			return false, err
		}
		for _, v := range mySqlIndex {
			if v.Column == column {
				if v.Key == "PRIMARY" {
					return true, err
				}
			}
		}
	case s.DriverPostgres:
		err := s.Query(ctx, db, nil, &postgresIndex, query)
		if err != nil {
			return false, err
		}
		for _, v := range postgresIndex {
			if v.Index == column {
				if strings.Contains(v.Index, "pkey") && strings.Contains(v.Index, column) {
					return true, err
				}
			}
		}
	case s.DriverMssql:
		err = s.Query(ctx, db, nil, &mssqlIndex, query)
		if err != nil {
			return false, err
		}
		for _, v := range mssqlIndex {
			if v.Column == column {
				if strings.Contains(v.Constraint, "PK") {
					return true, err
				}
			}
		}
	case s.DriverSqlite3:
		err = s.Query(ctx, db, nil, &sqliteIndex, query)
		if err != nil {
			return false, err
		}
		for _, v := range sqliteIndex {
			if strings.Contains(v.Name, column) {
				if v.Origin == "pk" {
					return true, err
				}
			}
		}
	case s.DriverOracle:

	}
	return false, err
} // Check if a column has primary tag

func listPrimaryKeys(ctx context.Context, db *sql.DB, database, table string) ([]PrimaryKey, error) { // Return a slice of Column of the composite key
	driver := s.GetDriver(db)
	var res []PrimaryKey
	query, err := query.ListAllPrimaryKeys(database, driver, table)
	if err != nil {
		return nil, err
	}
	err = s.Query(ctx, db, nil, &res, query)
	if err != nil {
		return nil, err
	}
	return res, err
}

func ListTables(ctx context.Context, db *sql.DB, database string) ([]Tables, error) {
	driver := s.GetDriver(db)
	var tables []Tables
	query, err := query.ListTablesQuery(database, "", driver)
	if err != nil {
		return nil, err
	}
	err = s.Query(ctx, db, nil, &tables, query)
	if err != nil {
		return nil, err
	}
	return tables, err
}

func GetRelationshipTable(ctx context.Context, db *sql.DB, database string, tables []string) ([]RelTables, error) {
	driver := s.GetDriver(db)
	var relTables []RelTables
	var sqliteRels []SqliteRel
	switch s.GetDriver(db) {
	case s.DriverSqlite3:
		listReferenceQuery, err := query.ListReferenceQuery(database, driver, "")
		if err != nil {
			return nil, err
		}
		err = s.Query(ctx, db, nil, &sqliteRels, listReferenceQuery)
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
					relTables = append(relTables, rel)
				}
			}
		}
		return relTables, err
	case s.DriverOracle:
		var relTable []RelTables
		for i := range tables {
			q, err := query.ListReferenceQuery(database, driver, tables[i])
			if err != nil {
				return nil, err
			}
			err = s.Query(ctx, db, nil, &relTable, q)
			if err != nil {
				return nil, err
			}
			relTables = append(relTables, relTable...)
		}
		return relTables, nil
	default:
		listReferenceQuery, err := query.ListReferenceQuery(database, driver, "")
		if err != nil {
			return nil, err
		}
		err = s.Query(ctx, db, nil, &relTables, listReferenceQuery)
		if err != nil {
			return nil, err
		}
		return relTables, err
	}
} // Find all columns, table and its referenced columns, tables
