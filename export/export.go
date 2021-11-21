package export

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-generator/core"
	"github.com/go-generator/core/build"
	d "github.com/go-generator/core/driver"
	gdb "github.com/go-generator/core/export/db"
	"github.com/go-generator/core/export/query"
	"github.com/go-generator/core/export/relationship"
	s "github.com/go-generator/core/export/sql"
	st "github.com/go-generator/core/strings"
)

var (
	sqliteNotNullIndex,
	tableFieldsIndex,
	primaryKeysIndex map[string]int
)

func init() {
	var err error
	sqliteNotNullIndex, err = s.GetColumnIndexes(reflect.TypeOf(relationship.SqliteNotNull{}))
	if err != nil {
		panic(err)
	}
	tableFieldsIndex, err = s.GetColumnIndexes(reflect.TypeOf(gdb.TableFields{}))
	if err != nil {
		panic(err)
	}
	primaryKeysIndex, err = s.GetColumnIndexes(reflect.TypeOf(relationship.PrimaryKey{}))
	if err != nil {
		panic(err)
	}
}

func ToModel(types map[string]string, table string, rt []relationship.RelTables, hasCompositeKey bool, sqlTable []gdb.TableFields) (*metadata.Model, error) { //s *TableInfo, conn *gorm.DB, tables []string, packageName, output string) {
	table = strings.ToLower(table)
	var m metadata.Model
	var raw string
	if !strings.Contains(table, "_") {
		raw = st.BuildSnakeName(table)
	} else {
		raw = st.UnBuildSnakeName(strings.ToLower(table))
	}
	n := st.ToSingular(raw)
	tableNames := build.BuildNames(n)
	m.Name = tableNames["Name"]
	m.Table = table
	for _, v := range sqlTable {
		colNames := build.BuildNames(strings.ToLower(v.Column))
		var f metadata.Field
		f.Column = v.Column
		f.Name = colNames["Name"]
		f.Type = types[v.DataType]
		if v.Length.Valid {
			l, err := strconv.Atoi(v.Length.String)
			if err != nil {
				return nil, err
			}
			f.Length = l
		}
		if v.ColumnKey == "PRI" {
			f.Key = true
		}
		m.Fields = append(m.Fields, f)
	}
	for _, ref := range rt {
		if ref.Table == table {
			refNames := build.BuildNames(ref.ReferencedTable)
			var relModel metadata.Relationship
			relModel.Ref = ref.ReferencedTable
			relModel.Model = st.ToSingular(refNames["Name"])
			relModel.Fields = append(relModel.Fields, metadata.Link{
				Column: ref.Column,
				To:     ref.ReferencedColumn,
			})
			if ref.Relationship == relationship.OneToMany {
				m.Arrays = append(m.Arrays, relModel)
			}
			if ref.Relationship == relationship.ManyToOne {
				m.Models = append(m.Models, relModel)
			}
			if ref.Relationship == relationship.OneToOne {
				m.Ones = append(m.Ones, relModel)
			}
		}
	}
	return &m, nil
}

func ToModels(ctx context.Context, db *sql.DB, database string, tables []string, rt []relationship.RelTables, types map[string]string, primaryKeys map[string][]string) ([]metadata.Model, error) {
	var projectModels []metadata.Model
	for _, t := range tables {
		var tablesData gdb.TableInfo
		err := InitTables(ctx, db, database, t, &tablesData, primaryKeys)
		if err != nil {
			return nil, err
		}
		m, err := ToModel(types, t, rt, tablesData.HasCompositeKey, tablesData.Fields)
		if err != nil {
			return nil, err
		}
		projectModels = append(projectModels, *m)
	}
	return projectModels, nil
}

func InitTables(ctx context.Context, db *sql.DB, database, table string, st *gdb.TableInfo, primaryKeys map[string][]string) error {
	query := ""
	switch s.GetDriver(db) {
	case d.Mysql:
		query = `
			SELECT 
				TABLE_NAME AS 'table',
				COLUMN_NAME AS 'column_name',
				DATA_TYPE AS 'type',
				IS_NULLABLE AS 'is_nullable',
				COLUMN_KEY AS 'column_key',
				CHARACTER_MAXIMUM_LENGTH AS 'length'
			FROM
				information_schema.columns
			WHERE
				TABLE_SCHEMA = '%v'
					AND TABLE_NAME = '%v'`
		query = fmt.Sprintf(query, database, table)
	case d.Postgres:
		query := `
			SELECT TABLE_NAME AS TABLE,
				COLUMN_NAME,
				IS_NULLABLE,
				CHARACTER_MAXIMUM_LENGTH AS LENGTH,
				UDT_NAME AS TYPE
			FROM INFORMATION_SCHEMA.COLUMNS
			WHERE TABLE_NAME = '%v';`
		query = fmt.Sprintf(query, table)
	case d.Mssql:
		query = `
			SELECT 
    			TABLE_NAME AS 'table',
    			COLUMN_NAME AS 'column_name',
				DATA_TYPE AS 'type',
				IS_NULLABLE AS 'is_nullable',
				CHARACTER_MAXIMUM_LENGTH AS 'length'
			FROM
				information_schema.columns
			WHERE
				TABLE_NAME = '%v'`
		query = fmt.Sprintf(query, table)
	case d.Sqlite3:
		query := `
		select name as 'column_name', type, pk as 'column_key' from pragma_table_info('%v');`
		query = fmt.Sprintf(query, table)
		var notNull []relationship.SqliteNotNull
		err := s.Query(ctx, db, tableFieldsIndex, &st.Fields, query)
		if err != nil {
			return err
		}
		query = `
		select * from pragma_table_info('%v');`
		query = fmt.Sprintf(query, table)
		err = s.Query(ctx, db, sqliteNotNullIndex, &notNull, query)
		if err != nil {
			return err
		}
		sqlitePKMap(st, notNull)
		return nil
	case d.Oracle:
		query = `
			SELECT
				col.owner AS "schema_name",
				col.table_name AS "table",
				col.column_name AS "column_name",
				col.data_type AS "type",
				col.data_length AS "length",
				col.nullable AS "is_nullable"
			FROM
				sys.all_tab_columns col
			INNER JOIN sys.all_tables t ON
				col.owner = t.owner
				AND col.table_name = t.table_name
			WHERE
				col.owner = '%v'
				AND col.table_name = '%v'
			ORDER BY
				col.column_id`
		query = fmt.Sprintf(query, database, table)
	}
	err := s.Query(ctx, db, tableFieldsIndex, &st.Fields, query)
	if err != nil {
		return err
	}
	for i := range st.Fields {
		if relationship.IsPrimaryKey(st.Fields[i].Column, table, primaryKeys) {
			st.Fields[i].ColumnKey = "PRI"
		}
	}
	st.HasCompositeKey = HasCKey(table, primaryKeys)
	for i := range st.Fields {
		if relationship.IsPrimaryKey(st.Fields[i].Column, table, primaryKeys) {
			st.Fields[i].ColumnKey = "PRI"
		}
	}
	st.HasCompositeKey = HasCKey(table, primaryKeys)
	return nil
}

func HasCKey(table string, pks map[string][]string) bool {
	if len(pks[table]) > 1 {
		return true
	}
	return false
}

func GetAllPrimaryKeys(ctx context.Context, db *sql.DB, dbName, driver string, tables []string) (map[string][]string, error) {
	primaryKeys := make(map[string][]string)
	for _, t := range tables {
		var keys []relationship.PrimaryKey
		q, err := query.ListAllPrimaryKeys(dbName, driver, t)
		err = s.Query(ctx, db, primaryKeysIndex, &keys, q)
		if err != nil {
			return nil, err
		}
		var pks []string
		for _, k := range keys {
			pks = append(pks, k.Column)
		}
		primaryKeys[t] = pks
	}
	return primaryKeys, nil
}

func sqlitePKMap(st *gdb.TableInfo, notNull []relationship.SqliteNotNull) {
	for i := range st.Fields {
		if st.Fields[i].ColumnKey == "1" {
			st.Fields[i].ColumnKey = "PRI"
		}
		re := regexp.MustCompile(`\d+`)
		if strings.Contains(st.Fields[i].DataType, "NUMERIC") {
			st.Fields[i].DataType = "NUMERIC"
		}
		if strings.Contains(st.Fields[i].DataType, "VARCHAR") {
			st.Fields[i].Length = sql.NullString{
				String: re.FindString(st.Fields[i].DataType),
				Valid:  true,
			}
			st.Fields[i].DataType = "VARCHAR"
		}
		if strings.Contains(st.Fields[i].DataType, "CHARACTER") {
			st.Fields[i].Length = sql.NullString{
				String: re.FindString(st.Fields[i].DataType),
				Valid:  true,
			}
			st.Fields[i].DataType = "CHARACTER"
		}
		if strings.Contains(st.Fields[i].DataType, "NVARCHAR") {
			st.Fields[i].Length = sql.NullString{
				String: re.FindString(st.Fields[i].DataType),
				Valid:  true,
			}
			st.Fields[i].DataType = "NVARCHAR"
		}
		for j := range notNull {
			if st.Fields[i].Column == notNull[j].Name {
				if notNull[j].NotNull {
					st.Fields[i].IsNullable = "0"
				} else {
					st.Fields[i].IsNullable = "1"
				}
			}
		}
	}
}
