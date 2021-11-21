package export

import (
	"context"
	"database/sql"
	"fmt"
	s "github.com/core-go/sql"
	"github.com/go-generator/core"
	"github.com/go-generator/core/build"
	edb "github.com/go-generator/core/export/db"
	"github.com/go-generator/core/export/query"
	"github.com/go-generator/core/export/relationship"
	st "github.com/go-generator/core/strings"
	"regexp"
	"strconv"
	"strings"
)

func ToModel(types map[string]string, table string, rt []relationship.RelTables, hasCompositeKey bool, sqlTable []edb.TableFields) (*metadata.Model, error) { //s *TableInfo, conn *gorm.DB, tables []string, packageName, output string) {
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
	m.Source = table
	for _, v := range sqlTable {
		colNames := build.BuildNames(strings.ToLower(v.Column))
		var f metadata.Field
		if hasCompositeKey {
			f.Source = colNames["name"]
		} else {
			if v.ColumnKey == "PRI" {
				f.Source = "_id"
			} else {
				f.Source = colNames["name"]
			}
		}
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
		//rl := getRelationship(v.Column, rt)
		//if rl != nil {
		//	var rls metadata.Relationship
		//	var foreign metadata.Field
		//	tmpMap := generator.BuildNames(rl.Table)
		//	foreign.Name = tmpMap["Name"]
		//	foreign.Source = tmpMap["name"]
		//	foreign.Type = "*[]" + tmpMap["Name"]
		//	if rl.Relationship == relationship.ManyToOne && table == rl.ReferencedTable { // have Many to One relation, add a field to the current struct
		//		rls.Ref = rl.Table
		//		rls.Fields = append(rls.Fields, metadata.Link{
		//			Column: rl.Column,
		//			To:     rl.ReferencedColumn,
		//		})
		//		if m.Arrays == nil {
		//			m.Arrays = append(m.Arrays, rls)
		//		} else {
		//			for j := range m.Arrays {
		//				if m.Arrays[j].Ref == rls.Ref {
		//					m.Arrays[j].Fields = append(m.Arrays[j].Fields, rls.Fields...)
		//					break
		//				}
		//				if j == len(m.Arrays)-1 {
		//					m.Arrays = append(m.Arrays, rls)
		//				}
		//			}
		//		}
		//		for i := range m.Fields {
		//			if m.Fields[i] == foreign {
		//				break
		//			}
		//			if i == len(m.Fields)-1 {
		//				m.Fields = append(m.Fields, foreign)
		//			}
		//		}
		//	}
		//}
		m.Fields = append(m.Fields, f)
	}
	return &m, nil
}

func ToModels(ctx context.Context, db *sql.DB, database string, tables []string, rt []relationship.RelTables, types map[string]string, primaryKeys map[string][]string) ([]metadata.Model, error) {
	var projectModels []metadata.Model
	for _, t := range tables {
		var tablesData edb.TableInfo
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

func InitTables(ctx context.Context, db *sql.DB, database, table string, st *edb.TableInfo, primaryKeys map[string][]string) error {
	query := ""
	switch s.GetDriver(db) {
	case s.DriverMysql:
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
	case s.DriverPostgres:
		query := `
			SELECT TABLE_NAME AS TABLE,
				COLUMN_NAME,
				IS_NULLABLE,
				CHARACTER_MAXIMUM_LENGTH AS LENGTH,
				UDT_NAME AS TYPE
			FROM INFORMATION_SCHEMA.COLUMNS
			WHERE TABLE_NAME = '%v';`
		query = fmt.Sprintf(query, table)
	case s.DriverMssql:
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
	case s.DriverSqlite3:
		query := `
		select name as 'column_name', type, pk as 'column_key' from pragma_table_info('%v');`
		query = fmt.Sprintf(query, table)
		var notNull []relationship.SqliteNotNull
		err := s.Query(ctx, db, nil, &st.Fields, query)
		if err != nil {
			return err
		}
		query = `
		select * from pragma_table_info('%v');`
		query = fmt.Sprintf(query, table)
		err = s.Query(ctx, db, nil, &notNull, query)
		if err != nil {
			return err
		}
		sqlitePKMap(st, notNull)
		return nil
	case s.DriverOracle:
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
	err := s.Query(ctx, db, nil, &st.Fields, query)
	if err != nil {
		return err
	}
	for i := range st.Fields {
		if IsPrimaryKey(st.Fields[i].Column, table, primaryKeys) {
			st.Fields[i].ColumnKey = "PRI"
		}
	}
	st.HasCompositeKey = HasCKey(table, primaryKeys)
	for i := range st.Fields {
		if IsPrimaryKey(st.Fields[i].Column, table, primaryKeys) {
			st.Fields[i].ColumnKey = "PRI"
		}
	}
	st.HasCompositeKey = HasCKey(table, primaryKeys)
	return nil
}

func IsPrimaryKey(key, table string, pks map[string][]string) bool {
	for i := range pks[table] {
		if key == pks[table][i] {
			return true
		}
	}
	return false
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
		err = s.Query(ctx, db, nil, &keys, q)
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

func sqlitePKMap(st *edb.TableInfo, notNull []relationship.SqliteNotNull) {
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
