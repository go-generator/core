package export

import (
	"context"
	"database/sql"
	"fmt"
	s "github.com/core-go/sql"
	"github.com/go-generator/core"
	edb "github.com/go-generator/core/export/db"
	"github.com/go-generator/core/export/relationship"
	"github.com/go-generator/core/generator"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func ToModel(types map[string]string, table string, rt []relationship.RelTables, hasCompositeKey bool, sqlTable []edb.TableFields) (*metadata.Model, error) { //s *TableInfo, conn *gorm.DB, tables []string, packageName, output string) {
	var m metadata.Model
	tableNames := generator.BuildNames(table)
	m.Name = tableNames["Name"]
	m.Table = tableNames["name"]
	m.Source = tableNames["name"]
	for _, v := range sqlTable {
		colNames := generator.BuildNames(v.Column)
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
		f.Name = colNames["name"]
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
		rl := getRelationship(v.Column, rt)
		if rl != nil {
			var rls metadata.Relationship
			var foreign metadata.Field
			tmpMap := generator.BuildNames(rl.Table)
			foreign.Name = tmpMap["Name"]
			foreign.Source = tmpMap["name"]
			foreign.Type = "*[]" + tmpMap["Names"]                                        // for many to many relationship
			if rl.Relationship == relationship.ManyToOne && table == rl.ReferencedTable { // have Many to One relation, add a field to the current struct
				rls.Ref = rl.Table
				rls.Fields = append(rls.Fields, metadata.Link{
					Column: rl.Column,
					To:     rl.ReferencedColumn,
				})
				if m.Arrays == nil {
					m.Arrays = append(m.Arrays, rls)
				} else {
					for j := range m.Arrays {
						if m.Arrays[j].Ref == rls.Ref {
							m.Arrays[j].Fields = append(m.Arrays[j].Fields, rls.Fields...)
							break
						}
						if j == len(m.Arrays)-1 {
							m.Arrays = append(m.Arrays, rls)
						}
					}
				}
				for i := range m.Fields {
					if m.Fields[i] == foreign {
						break
					}
					if i == len(m.Fields)-1 {
						m.Fields = append(m.Fields, foreign)
					}
				}
			}
		}
		m.Fields = append(m.Fields, f)
	}
	return &m, nil
}

func getRelationship(column string, rt []relationship.RelTables) *relationship.RelTables {
	for _, v := range rt {
		if column == v.ReferencedColumn {
			return &v
		}
	}
	return nil
}

func ToModels(ctx context.Context, db *sql.DB, database string, tables []string, rt []relationship.RelTables, types map[string]string) ([]metadata.Model, error) {
	var projectModels []metadata.Model
	for _, t := range tables {
		var tablesData edb.TableInfo
		err := InitTables(ctx, db, database, t, &tablesData)
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

func InitTables(ctx context.Context, db *sql.DB, database, table string, st *edb.TableInfo) error {
	switch s.GetDriver(db) {
	case s.DriverMysql:
		query := `
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
		err := s.Query(ctx, db, nil, &st.Fields, query)
		if err != nil {
			return err
		}
		st.HasCompositeKey = HasCompositeKey(st.Fields)
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
		err := s.Query(ctx, db, nil, &st.Fields, query)
		if err != nil {
			return err
		}
		count := 0
		for i := range st.Fields {
			c, err := relationship.CheckPrimaryTag(ctx, db, database, s.GetDriver(db), table, st.Fields[i].Column)
			if err != nil {
				return err
			}
			if c {
				st.Fields[i].ColumnKey = "PRI"
				count++
			}
		}
		if count < 2 {
			st.HasCompositeKey = true
		} else {
			st.HasCompositeKey = false
		}
	case s.DriverMssql:
		query := `
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
		err := s.Query(ctx, db, nil, &st.Fields, query)
		if err != nil {
			return err
		}
		count := 0
		for i := range st.Fields {
			c, err := relationship.CheckPrimaryTag(ctx, db, database, s.DriverMssql, table, st.Fields[i].Column)
			if err != nil {
				return err
			}
			if c {
				st.Fields[i].ColumnKey = "PRI"
				count++
			}
		}
		if count < 2 {
			st.HasCompositeKey = true
		} else {
			st.HasCompositeKey = false
		}
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
		for i := range st.Fields {
			if st.Fields[i].ColumnKey == "1" {
				st.Fields[i].ColumnKey = "PRI"
			}
			re := regexp.MustCompile(`\d+`)
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
	return nil
}

func HasCompositeKey(st []edb.TableFields) bool {
	var count int
	for _, v := range st {
		if v.ColumnKey == "PRI" {
			count++
		}
	}
	return count < 2
}

func DetectDriver(s string) string {
	if strings.Index(s, "sqlserver:") == 0 {
		return "mssql"
	} else {
		if strings.Index(s, "user=") >= 0 && strings.Index(s, "password=") >= 0 {
			if strings.Index(s, "dbname=") >= 0 || strings.Index(s, "host=") >= 0 || strings.Index(s, "port=") >= 0 {
				return "postgres"
			} else {
				return "godror"
			}
		} else {
			_, err := filepath.Abs(s)
			if (strings.Index(s, "@tcp(") >= 0 || strings.Index(s, "charset=") > 0 || strings.Index(s, "parseTime=") > 0 || strings.Index(s, "loc=") > 0 || strings.Index(s, "@") >= 0 || strings.Index(s, ":") >= 0) && err != nil {
				return "mysql"
			} else {
				return "sqlite3"
			}
		}
	}
}
