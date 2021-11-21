package project

import (
	"database/sql"
	"errors"
	s "github.com/core-go/sql"
	metadata "github.com/go-generator/core"
	"path/filepath"
	"strings"
)

func ConnectDB(dbCache metadata.Database, driver string) (*sql.DB, error) {
	switch driver {
	case s.DriverMysql:
		return sql.Open(driver, dbCache.MySql)
	case s.DriverPostgres:
		return sql.Open(driver, dbCache.Postgres)
	case s.DriverMssql:
		return sql.Open(driver, dbCache.Mssql)
	case s.DriverSqlite3:
		return sql.Open(driver, dbCache.Sqlite3)
	case s.DriverOracle:
		return sql.Open("godror", dbCache.Oracle)
	default:
		return nil, errors.New(s.DriverNotSupport)
	}
}

func SelectDSN(dbCache metadata.Database, driver string) string {
	switch driver {
	case s.DriverMysql:
		return dbCache.MySql
	case s.DriverPostgres:
		return dbCache.Postgres
	case s.DriverMssql:
		return dbCache.Mssql
	case s.DriverSqlite3:
		return dbCache.Sqlite3
	case s.DriverOracle:
		return dbCache.Oracle
	default:
		return ""
	}
}

func UpdateDBCache(dbCache *metadata.Database, driver, dsn string) {
	switch driver {
	case s.DriverMysql:
		dbCache.MySql = dsn
	case s.DriverPostgres:
		dbCache.Postgres = dsn
	case s.DriverMssql:
		dbCache.Mssql = dsn
	case s.DriverSqlite3:
		dbCache.Sqlite3 = dsn
	case s.DriverOracle:
		dbCache.Oracle = dsn
	}
}

func GetDatabaseName(dbCache metadata.Database, driver string) (string, error) {
	switch driver {
	case s.DriverMysql:
		s1 := strings.Split(dbCache.MySql, "/")
		if len(s1) < 2 {
			return "", errors.New("invalid datasource")
		}
		s2 := strings.Split(s1[1], "?")
		return s2[0], nil
	case s.DriverPostgres:
		s1 := strings.Split(dbCache.Postgres, "dbname=")
		if len(s1) < 2 {
			return "", errors.New("invalid datasource")
		}
		s2 := strings.Split(s1[1], " ")
		return s2[0], nil
	case s.DriverMssql:
		s1 := strings.Split(dbCache.Mssql, "database=")
		if len(s1) < 2 {
			return "", errors.New("invalid datasource")
		}
		s2 := strings.Split(s1[1], "&")
		return s2[0], nil
	case s.DriverSqlite3:
		return filepath.Base(dbCache.Sqlite3), nil
	default:
		return "", errors.New(s.DriverNotSupport)
	}
}
