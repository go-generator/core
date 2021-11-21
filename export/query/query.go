package query

import (
	"errors"
	"fmt"

	d "github.com/go-generator/core/driver"
)

func ListTablesQuery(database, dbname, driver string) (string, error) {
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
	case d.Mssql:
		return `
		SELECT name 'table'
		FROM     sys.sysobjects
		WHERE  (xtype = 'U')`, nil
	case d.Sqlite3:
		return `
		SELECT 
			name as 'table'
		FROM 
			sqlite_schema
		WHERE 
			type ='table' AND 
			name NOT LIKE 'sqlite_%';`, nil
	case d.Oracle:
		query := `
		SELECT
			DISTINCT OWNER,
			OBJECT_NAME "table"
		FROM
			DBA_OBJECTS
		WHERE
			OBJECT_TYPE = 'TABLE'
			AND OWNER = '%v'`
		return fmt.Sprintf(query, dbname), nil
	default:
		return "", errors.New("unsupported driver")
	}
}

func ListUniqueQuery(database, driver, table string) (string, error) {
	switch driver {
	case d.Mysql:
		query := `show indexes from %v.%v`
		return fmt.Sprintf(query, database, table), nil
	case d.Postgres:
		query := `
		SELECT TABLENAME AS TABLE,
			INDEXNAME AS INDEX
		FROM PG_INDEXES
		WHERE TABLENAME = '%v'`
		return fmt.Sprintf(query, table), nil
	case d.Mssql:
		query := `
			SELECT 
    			COLUMN_NAME AS 'column', CONSTRAINT_NAME AS 'constraint'
			FROM
    			INFORMATION_SCHEMA.CONSTRAINT_COLUMN_USAGE
			WHERE
    			(TABLE_NAME = '%v')`
		return fmt.Sprintf(query, table), nil
	case d.Sqlite3:
		query := `
		PRAGMA INDEX_LIST('%v');`
		return fmt.Sprintf(query, table), nil
	default:
		return "", errors.New(d.NotSupport)
	}
}

func ListCompositeKeyQuery(database, driver, table string) (string, error) { //TODO: get composite keys for other databases
	switch driver {
	case d.Mysql:
		query := `
			SELECT 
    			COLUMN_NAME as 'column'
			FROM
				information_schema.KEY_COLUMN_USAGE
			WHERE
				table_schema = '%v'
					AND table_name = '%v'
					AND constraint_name = 'PRIMARY';`
		return fmt.Sprintf(query, database, table), nil
	case d.Mssql:
		query := `
			SELECT K.COLUMN_NAME 'column'
			FROM     INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS K INNER JOIN
                     INFORMATION_SCHEMA.TABLE_CONSTRAINTS AS TC ON K.TABLE_CATALOG = TC.TABLE_CATALOG AND K.TABLE_SCHEMA = TC.TABLE_SCHEMA AND K.CONSTRAINT_NAME = TC.CONSTRAINT_NAME
			WHERE  (TC.CONSTRAINT_TYPE = 'PRIMARY KEY') AND (K.TABLE_NAME = '%v')`
		return fmt.Sprintf(query, table), nil
	default:
		return "", errors.New(d.NotSupport)
	}
}

func ListAllPrimaryKeys(database, driver, table string) (string, error) {
	query := ""
	switch driver {
	case d.Mysql:
		query = `
			SELECT 
				K.COLUMN_NAME as 'column'
			FROM
				INFORMATION_SCHEMA.TABLE_CONSTRAINTS AS C
					JOIN
				INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS K ON C.TABLE_NAME = K.TABLE_NAME
					AND C.CONSTRAINT_CATALOG = K.CONSTRAINT_CATALOG
					AND C.CONSTRAINT_SCHEMA = K.CONSTRAINT_SCHEMA
					AND C.CONSTRAINT_NAME = K.CONSTRAINT_NAME
			WHERE
				C.TABLE_SCHEMA = '%v'
					AND K.TABLE_NAME = '%v'
					AND C.CONSTRAINT_TYPE = 'PRIMARY KEY'`
		return fmt.Sprintf(query, database, table), nil
	case d.Postgres:
		query := `
		SELECT 
			kc.column_name as column
		FROM
			information_schema.table_constraints tc
				JOIN
			information_schema.key_column_usage kc ON kc.table_name = tc.table_name
				AND kc.table_schema = tc.table_schema
				AND kc.constraint_name = tc.constraint_name
		WHERE
			tc.constraint_type = 'PRIMARY KEY'
				AND kc.ordinal_position IS NOT NULL
				AND tc.table_name = '%v'
		ORDER BY tc.table_schema , tc.table_name, kc.position_in_unique_constraint;`
		return fmt.Sprintf(query, table), nil
	case d.Mssql:
		query = `
			SELECT K.COLUMN_NAME 'column'
			FROM     INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS K INNER JOIN
                     INFORMATION_SCHEMA.TABLE_CONSTRAINTS AS TC ON K.TABLE_CATALOG = TC.TABLE_CATALOG AND K.TABLE_SCHEMA = TC.TABLE_SCHEMA AND K.CONSTRAINT_NAME = TC.CONSTRAINT_NAME
			WHERE  (TC.CONSTRAINT_TYPE = 'PRIMARY KEY') AND (K.TABLE_NAME = '%v')`
		return fmt.Sprintf(query, table), nil
	case d.Sqlite3:
		query := `
		select name as 'column' from pragma_table_info('%v') as 'p' where p.pk = TRUE`
		return fmt.Sprintf(query, table), nil
	case d.Oracle:
		query := `
		SELECT
		all_cons_columns.table_name AS "table",
		all_cons_columns.column_name AS "column"
		FROM
			all_constraints,
			all_cons_columns
		WHERE
			all_constraints.constraint_type = 'P'
			AND all_constraints.constraint_name = all_cons_columns.constraint_name
			AND all_constraints.owner = all_cons_columns.owner
			AND all_cons_columns.table_name = '%v'
			AND all_cons_columns.owner = '%v'
		ORDER BY
			all_cons_columns.owner,
			all_cons_columns.table_name,
			all_cons_columns.position`
		return fmt.Sprintf(query, table, database), nil
	}
	return query, nil
}

func ListReferenceQuery(database, driver, table string) (string, error) {
	switch driver {
	case d.Mysql:
		query := `
		SELECT
			TABLE_NAME as 'table',
			COLUMN_NAME as 'column',
			REFERENCED_TABLE_NAME as 'referenced_table',
			REFERENCED_COLUMN_NAME as 'referenced_column'
		FROM
    		information_schema.key_column_usage
		WHERE
    		constraint_schema = '%v'
        	AND referenced_table_schema IS NOT NULL
        	AND referenced_table_name IS NOT NULL
        	AND referenced_column_name IS NOT NULL`
		return fmt.Sprintf(query, database), nil
	case d.Postgres:
		return `
		SELECT
			TC.TABLE_NAME AS table,
			KCU.COLUMN_NAME AS column,
			CCU.TABLE_NAME AS referenced_table,
			CCU.COLUMN_NAME AS referenced_column
		FROM INFORMATION_SCHEMA.TABLE_CONSTRAINTS AS TC
		JOIN INFORMATION_SCHEMA.KEY_COLUMN_USAGE AS KCU ON TC.CONSTRAINT_NAME = KCU.CONSTRAINT_NAME
		AND TC.TABLE_SCHEMA = KCU.TABLE_SCHEMA
		JOIN INFORMATION_SCHEMA.CONSTRAINT_COLUMN_USAGE AS CCU ON CCU.CONSTRAINT_NAME = TC.CONSTRAINT_NAME
		AND CCU.TABLE_SCHEMA = TC.TABLE_SCHEMA
		WHERE TC.CONSTRAINT_TYPE = 'FOREIGN KEY';`, nil
	case d.Mssql:
		return `
		SELECT tp.name AS 'table', cp.name AS 'column', tr.name AS 'referenced_table', cr.name AS 'referenced_column'
		FROM     sys.foreign_keys AS fk INNER JOIN
                  sys.tables AS tp ON fk.parent_object_id = tp.object_id INNER JOIN
                  sys.tables AS tr ON fk.referenced_object_id = tr.object_id INNER JOIN
                  sys.foreign_key_columns AS fkc ON fkc.constraint_object_id = fk.object_id INNER JOIN
                  sys.columns AS cp ON fkc.parent_column_id = cp.column_id AND fkc.parent_object_id = cp.object_id INNER JOIN
                  sys.columns AS cr ON fkc.referenced_column_id = cr.column_id AND fkc.referenced_object_id = cr.object_id
		ORDER BY 'table'`, nil
	case d.Sqlite3:
		return `
		SELECT * FROM sqlite_master WHERE type = 'table' AND sql LIKE '%FOREIGN KEY%'`, nil
	case d.Oracle:
		query := `
		SELECT
			uc_r.table_name AS "table",
			ucc_r.column_name AS "column",
			uc_p.table_name AS "referenced_table",
			ucc_p.column_name AS "referenced_column"
		FROM
			user_constraints uc_r
		JOIN user_cons_columns ucc_r ON
			ucc_r.constraint_name = uc_r.constraint_name
		JOIN user_constraints uc_p ON
			uc_p.constraint_name = uc_r.r_constraint_name
		JOIN user_cons_columns ucc_p ON
			ucc_p.constraint_name = uc_p.constraint_name
			AND ucc_p.position = ucc_r.position
		WHERE
			uc_r.constraint_type = 'R'
			AND UC_R.owner = '%v'`
		return fmt.Sprintf(query, database), nil
	default:
		return "", errors.New(d.NotSupport)
	}
}
