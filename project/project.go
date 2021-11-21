package project

import (
	"database/sql"
	"errors"
	metadata "github.com/go-generator/core"
	d "github.com/go-generator/core/driver"
	"os"
	"path/filepath"
	"strings"
)

const (
	TypesJsonEnv   = "G_TYPES_JSON"
	WindowsIconEnv = "G_WINDOWS_ICON"
	AppIconEnv     = "G_APP_ICON"
	TemplatePath   = "G_TEMPLATE_PATH"
	ConfigEnv      = "G_CONFIG_PATH"
)

func SetPathEnv(key, value string) error {
	path, err := filepath.Abs(value)
	if err != nil {
		return err
	}
	if os.Getenv(key) == "" {
		err = os.Setenv(key, path)
		if err != nil {
			return err
		}
	}
	return err
}

func ConnectDB(dbCache metadata.Database, driver string) (*sql.DB, error) {
	switch driver {
	case d.Mysql:
		return sql.Open(driver, dbCache.MySql)
	case d.Postgres:
		return sql.Open(driver, dbCache.Postgres)
	case d.Mssql:
		return sql.Open(driver, dbCache.Mssql)
	case d.Sqlite3:
		return sql.Open(driver, dbCache.Sqlite3)
	case d.Oracle:
		return sql.Open("godror", dbCache.Oracle)
	default:
		return nil, errors.New(d.NotSupport)
	}
}

func SelectDSN(dbCache metadata.Database, driver string) string {
	switch driver {
	case d.Mysql:
		return dbCache.MySql
	case d.Postgres:
		return dbCache.Postgres
	case d.Mssql:
		return dbCache.Mssql
	case d.Sqlite3:
		return dbCache.Sqlite3
	case d.Oracle:
		return dbCache.Oracle
	default:
		return ""
	}
}

func UpdateDBCache(dbCache *metadata.Database, driver, dsn string) {
	switch driver {
	case d.Mysql:
		dbCache.MySql = dsn
	case d.Postgres:
		dbCache.Postgres = dsn
	case d.Mssql:
		dbCache.Mssql = dsn
	case d.Sqlite3:
		dbCache.Sqlite3 = dsn
	case d.Oracle:
		dbCache.Oracle = dsn
	}
}

func GetDatabaseName(dbCache metadata.Database, driver string) (string, error) {
	switch driver {
	case d.Mysql:
		return GetName(dbCache.MySql)
	case d.Postgres:
		s1 := strings.Split(dbCache.Postgres, "dbname=")
		if len(s1) < 2 {
			return GetName(dbCache.Postgres)
		}
		s2 := strings.Split(s1[1], " ")
		return s2[0], nil
	case d.Mssql:
		s1 := strings.Split(dbCache.Mssql, "database=")
		if len(s1) < 2 {
			return "", errors.New("invalid datasource")
		}
		s2 := strings.Split(s1[1], "&")
		return s2[0], nil
	case d.Sqlite3:
		return filepath.Base(dbCache.Sqlite3), nil
	case d.Oracle:
		s1 := strings.Split(dbCache.Oracle, "//")
		if len(s1) < 2 {
			return "", errors.New("invalid datasource")
		}
		s2 := strings.Split(s1[1], ":")
		return s2[0], nil
	default:
		return "", errors.New(d.NotSupport)
	}
}
func GetName(s string) (string, error) {
	i := strings.LastIndex(s, "/")
	if i >= 0 {
		j := strings.LastIndex(s, "?")
		if j < 0 {
			return s[i+1:], nil
		} else {
			return s[i+1:j], nil
		}
	}
	return "", errors.New("invalid datasource")
}
