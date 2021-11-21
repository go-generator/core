package db

import (
	"errors"
	"path/filepath"
	"strings"
)

const (
	DriverPostgres   = "postgres"
	DriverMysql      = "mysql"
	DriverMssql      = "mssql"
	DriverOracle     = "oracle"
	DriverSqlite3    = "sqlite3"
	DriverNotSupport = "no support"
)

func DetectDriver(s string) string {
	if strings.Index(s, "sqlserver:") == 0 {
		return "mssql"
	} else if strings.Index(s, "postgres:") == 0 {
		return "postgres"
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

func ExtractDBName(dsn string, driver string) (string, error) {
	switch driver {
	case DriverMysql:
		s1 := strings.Split(dsn, "/")
		if len(s1) < 2 {
			return "", errors.New("invalid datasource")
		}
		s2 := strings.Split(s1[1], "?")
		return s2[0], nil
	case DriverPostgres:
		s1 := strings.Split(dsn, "dbname=")
		if len(s1) < 2 {
			return "", errors.New("invalid datasource")
		}
		s2 := strings.Split(s1[1], " ")
		return s2[0], nil
	case DriverMssql:
		s1 := strings.Split(dsn, "database=")
		if len(s1) < 2 {
			return "", errors.New("invalid datasource")
		}
		s2 := strings.Split(s1[1], "&")
		return s2[0], nil
	case DriverSqlite3:
		return filepath.Base(dsn), nil
	default:
		return "", errors.New(DriverNotSupport)
	}
}
