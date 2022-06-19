package export

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/gertd/go-pluralize"
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
func Trim(s string, prefix, suffix string) string {
	if len(prefix) > 0 {
		s = s[len(prefix):]
	}
	if len(suffix) > 0 {
		s = s[:len(s) - len(suffix)]
	}
	return s
}
func ToModel(types map[string]string, table string, rt []relationship.RelTables, sqlTable []gdb.TableFields, options...string) (*metadata.Model, error) { //s *TableInfo, conn *gorm.DB, tables []string, packageName, output string) {
	pluralize := pluralize.NewClient()
	origin := table
	table = strings.ToLower(table)
	prefix := ""
	suffix := ""
	if len(options) > 0 {
		prefix = options[0]
	}
	if len(options) > 1 {
		suffix = options[1]
	}
	table = Trim(table, prefix, suffix)
	var m metadata.Model
	var raw string
	if !strings.Contains(table, "_") {
		raw = st.BuildSnakeName(table)
	} else {
		raw = st.UnBuildSnakeName(strings.ToLower(table))
	}
	n := pluralize.Singular(raw)
	tableNames := build.BuildNames(n)
	if tableNames["Name"] == origin || tableNames["name"] == origin {
		m.Name = origin
	} else {
		m.Name = tableNames["Name"]
		m.Table = origin
	}
	for _, v := range sqlTable {
		var f metadata.Field
		org := v.Column
		x := v.Column
		if strings.ToUpper(x) == v.Column {
			x = strings.ToLower(x)
		}
		colNames := build.BuildNames(x)
		if colNames["Name"] == org || colNames["name"] == org {
			f.Name = org
		} else {
			f.Name = colNames["Name"]
		}
		if v.Column == colNames["Name"] || x == colNames["name"] {
			f.Name = x
		} else {
			f.Column = v.Column
			f.Name = colNames["Name"]
		}
		f.Type = types[v.DataType]
		f.DbType = v.DbType
		f.FullDbType = v.FullDataType
		if v.Length != nil {
			l, err := strconv.Atoi(*v.Length)
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
		if ref.Table == origin {
			refNames := build.BuildNames(ref.ReferencedTable)
			var relModel metadata.Relationship
			relModel.Ref = ref.ReferencedTable
			relModel.Model = pluralize.Singular(refNames["Name"])
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

func ToModels(ctx context.Context, db *sql.DB, database string, tables []string, rt []relationship.RelTables, types map[string]string, primaryKeys map[string][]string, options...string) ([]metadata.Model, error) {
	var projectModels []metadata.Model
	for _, t := range tables {
		var tablesData gdb.TableInfo
		err := InitTables(ctx, db, database, t, &tablesData, primaryKeys)
		if err != nil {
			return nil, err
		}
		m, err := ToModel(types, t, rt, tablesData.Fields, options...)
		if err != nil {
			return nil, err
		}
		projectModels = append(projectModels, *m)
	}
	return projectModels, nil
}

func InitTables(ctx context.Context, db *sql.DB, database, table string, st *gdb.TableInfo, primaryKeys map[string][]string) error {
	query := ""
	driver := s.GetDriver(db)
	switch driver {
	case d.Mysql:
		query = `
			select
				table_name as 'table',
				column_name as 'column_name',
				data_type as 'type',
				is_nullable as 'is_nullable',
				column_key as 'column_key',
				character_maximum_length as 'length',
				numeric_precision as 'precision',
				numeric_scale as 'scale'
			from
				information_schema.columns
			where
				table_schema = '%v'
				and table_name = '%v'`
		query = fmt.Sprintf(query, database, table)
		err := s.Query(ctx, db, tableFieldsIndex, &st.Fields, query)
		if err != nil {
			return err
		}
		for i := range st.Fields {
			if st.Fields[i].Length != nil && strings.Contains(st.Fields[i].DataType, "char") {
				st.Fields[i].FullDataType = fmt.Sprintf("%s(%s)", st.Fields[i].DataType, *st.Fields[i].Length)
			}
		}
	case d.Postgres:
		query = `
			select
				table_name as "table",
				column_name,
				is_nullable,
				character_maximum_length as "length",
				udt_name as "type",
				numeric_scale as "scale",
				numeric_precision as "precision"
			from
				information_schema.columns
			where
				table_name = '%v';`
		query = fmt.Sprintf(query, table)
		err := s.Query(ctx, db, tableFieldsIndex, &st.Fields, query)
		if err != nil {
			return err
		}
		for i := range st.Fields {
			if st.Fields[i].Length != nil && strings.Contains(st.Fields[i].DataType, "char") && strings.Index(st.Fields[i].DataType, "_") < 0 {
				st.Fields[i].FullDataType = fmt.Sprintf("%s(%s)", st.Fields[i].DataType, *st.Fields[i].Length)
			}
		}
	case d.Mssql:
		query = `
			select
				table_name as 'table',
				column_name as 'column_name',
				data_type as 'type',
				is_nullable as 'is_nullable',
				character_maximum_length as 'length'
			from
				information_schema.columns
			where
				table_name = '%v'`
		query = fmt.Sprintf(query, table)
		err := s.Query(ctx, db, tableFieldsIndex, &st.Fields, query)
		if err != nil {
			return err
		}
		for i := range st.Fields {
			if st.Fields[i].Length != nil && strings.Contains(st.Fields[i].DataType, "char") {
				st.Fields[i].FullDataType = fmt.Sprintf("%s(%s)", st.Fields[i].DataType, *st.Fields[i].Length)
			}
		}
	case d.Sqlite3:
		query = `
			select
				name as 'column_name',
				type,
				pk as 'column_key'
			from
				pragma_table_info('%v');`
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
		for i := range st.Fields {
			if st.Fields[i].Length != nil && (st.Fields[i].DataType == "TEXT" || strings.Contains(st.Fields[i].DataType, "CHAR")) {
				st.Fields[i].FullDataType = fmt.Sprintf("%s(%s)", st.Fields[i].DataType, *st.Fields[i].Length)
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
	case d.Oracle:
		query = `
			SELECT
				col.owner AS "schema_name",
				col.table_name AS "table",
				col.column_name AS "column_name",
				col.data_type AS "type",
				col.data_length AS "length",
				col.nullable AS "is_nullable",
				col.data_precision AS "precision",
				col.data_scale AS "scale"
			FROM
				sys.all_tab_columns col
			INNER JOIN sys.all_tables t ON
				col.table_name = t.table_name
			WHERE col.table_name = '%v'
			ORDER BY
				col.column_id`
		query = fmt.Sprintf(query, table)
		err := s.Query(ctx, db, tableFieldsIndex, &st.Fields, query)
		if err != nil {
			return err
		}
		for i := range st.Fields {
			if st.Fields[i].Length != nil && strings.Contains(st.Fields[i].DataType, "CHAR") {
				st.Fields[i].FullDataType = fmt.Sprintf("%s(%s BYTE)", st.Fields[i].DataType, *st.Fields[i].Length)
			}
		}
	}
	for i := range st.Fields {
		st.Fields[i].DbType = st.Fields[i].DataType
		if relationship.IsPrimaryKey(st.Fields[i].Column, table, primaryKeys) {
			st.Fields[i].ColumnKey = "PRI"
		}
	}
	mapDataType(driver, st)
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

func mapDataType(driver string, st *gdb.TableInfo) {
	re := regexp.MustCompile(`\d+`)
	switch driver {
	case d.Mysql:
		for i := range st.Fields {
			if strings.Contains(st.Fields[i].DataType, "decimal") {
				scale := 0
				precision := 38 // default precision for oracle
				if st.Fields[i].Scale != nil {
					scale = *st.Fields[i].Scale
				}
				if st.Fields[i].Precision != nil {
					precision = *st.Fields[i].Precision
				}
				if st.Fields[i].Scale == nil && st.Fields[i].Precision == nil {
					st.Fields[i].DataType = "decimal(7,0)"
				} else {
					if scale == 0 {
						if precision <= 2 {
							st.Fields[i].DataType = "decimal(2,0)"
						} else {
							if precision <= 4 {
								st.Fields[i].DataType = "decimal(4,0)"
							} else {
								if precision <= 6 {
									st.Fields[i].DataType = "decimal(6,0)"
								} else {
									st.Fields[i].DataType = "decimal(7,0)"
								}
							}
						}
					} else {
						st.Fields[i].DataType = "decimal"
					}
				}
			}
		}
	case d.Sqlite3:
		for i := range st.Fields {
			l := re.FindString(st.Fields[i].DataType)
			if st.Fields[i].ColumnKey == "1" {
				st.Fields[i].ColumnKey = "PRI"
			}
			if strings.Contains(st.Fields[i].DataType, "NUMERIC") {
				st.Fields[i].DataType = "NUMERIC"
			}
			if strings.Contains(st.Fields[i].DataType, "VARCHAR") {
				st.Fields[i].Length = &l
				st.Fields[i].DataType = "VARCHAR"
			}
			if strings.Contains(st.Fields[i].DataType, "CHARACTER") {
				st.Fields[i].Length = &l
				st.Fields[i].DataType = "CHARACTER"
			}
			if strings.Contains(st.Fields[i].DataType, "NVARCHAR") {
				st.Fields[i].Length = &l
				st.Fields[i].DataType = "NVARCHAR"
			}
		}
	case d.Oracle:
		for i := range st.Fields {
			if strings.Contains(st.Fields[i].DataType, "NUMERIC") {
				st.Fields[i].DataType = "NUMERIC"
			}
			if strings.Contains(st.Fields[i].DataType, "NUMBER") {
				scale := 0
				precision := 38 // default precision for oracle
				if st.Fields[i].Scale != nil {
					scale = *st.Fields[i].Scale
				}
				if st.Fields[i].Precision != nil {
					precision = *st.Fields[i].Precision
				}
				if st.Fields[i].Scale == nil && st.Fields[i].Precision == nil {
					st.Fields[i].DataType = "NUMBER(7,0)"
				} else {
					if scale == 0 {
						if precision <= 2 {
							st.Fields[i].DataType = "NUMBER(2,0)"
						} else {
							if precision <= 4 {
								st.Fields[i].DataType = "NUMBER(4,0)"
							} else {
								if precision <= 6 {
									st.Fields[i].DataType = "NUMBER(6,0)"
								} else {
									st.Fields[i].DataType = "NUMBER(7,0)"
								}
							}
						}
					} else {
						st.Fields[i].DataType = "NUMBER"
					}
				}
			}
		}
	}
}
