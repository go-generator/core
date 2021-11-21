package strings

import "strings"

const (
	DriverPostgres = "postgres"
	DriverMysql    = "mysql"
	DriverMssql    = "mssql"
	DriverOracle   = "oracle"
	DriverSqlite3  = "sqlite3"
)

var ends = []string{"ies", "ees", "aes", "ues", "ves", "ays", "eys", "iys", "oys", "uys", "aos", "eos", "ios", "oos", "uos", "oes", "ses", "xes", "zes", "shes", "ches"}
var replaces = []string{"y", "ee", "ae", "ue", "fe", "ay", "ey", "iy", "oy", "uy", "ao", "eo", "io", "oo", "uo", "o", "s", "x", "z", "sh", "ch"}
var plural = []string{"people", "children", "women", "men", "fungus", "feet", "teeth", "geese", "mice", "gasses", "phenomena", "criteria", "sheep", "series", "species", "deer", "fish", "roofs", "beliefs", "chefs", "chiefs", "photos", "pianos", "halos", "volcanos", "volcanoes", "fezzes"}
var singular = []string{"person", "child", "woman", "man", "fungi", "foot", "tooth", "goose", "mouse", "gas", "phenomenon", "criterion", "sheep", "series", "species", "deer", "fish", "roof", "belief", "chef", "chief", "photo", "piano", "halo", "volcano", "volcano", "fez"}

func ToCamelCase(s string) string {
	if len(s) <= 2 {
		return strings.ToLower(s)
	} else {
		return strings.ToLower(string(s[0])) + strings.ToLower(s[1:])
	}
}
func ToPascalCase(s string) string {
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}
func ToSingular(s string) string {
	if len(s) <= 1 {
		return s
	}
	for i, si := range plural {
		if strings.HasSuffix(s, si) {
			return s[0:len(s)-len(si)] + singular[i]
		}
	}
	for i, si := range ends {
		if strings.HasSuffix(s, si) {
			return s[0:len(s)-len(si)] + replaces[i]
		}
	}
	if strings.HasSuffix(s, "s") {
		return s[0 : len(s)-1]
	}
	return s
}
func ToPlural(s string) string {
	if len(s) <= 1 {
		return s
	}
	for i, si := range singular {
		if strings.HasSuffix(s, si) {
			return s[0:len(s)-len(si)] + plural[i]
		}
	}
	x := s[len(s)-1:]
	if x == "y" {
		return s[0:len(s)-1] + "ies"
	}
	if x == "s" || x == "x" || strings.HasSuffix(s, "sh") || strings.HasSuffix(s, "ch") {
		return s[0:] + "es"
	}
	return s[0:] + "s"
}
func BuildSnakeName(s string) string {
	s2 := strings.ToLower(s)
	s3 := ""
	for i := range s {
		if strings.ToLower(string(s2[i])) != strings.ToLower(string(s[i])) {
			s3 += "_" + string(s2[i])
		} else {
			s3 += string(s2[i])
		}
	}
	if string(s3[0]) == "_" {
		return s3[1:]
	}
	return s3
}
func UnBuildSnakeName(s string) string {
	s2 := strings.ToUpper(s)
	s1 := string(s[0])
	for i := 1; i < len(s); i++ {
		if string(s[i-1]) == "_" {
			s1 = s1[:len(s1)-1]
			s1 += string(s2[i])
		} else {
			s1 += string(s[i])
		}
	}
	return s1
}

func ImportDriver(s string) string {
	switch s {
	case DriverMysql:
		return `_ "github.com/go-sql-driver/mysql"`
	case DriverMssql:
		return `_ "github.com/denisenkom/go-mssqldb"`
	case DriverPostgres:
		return `_ "github.com/lib/pq"`
	case DriverSqlite3:
		return `_ "github.com/mattn/go-sqlite3"`
	case DriverOracle:
		return `_ "github.com/godror/godror"`
	case "godror":
		return `_ "github.com/godror/godror"`
	default:
		return ""
	}
}

func ImportGoMod(s string) string {
	switch s {
	case DriverMysql:
		return `github.com/go-sql-driver/mysql v1.6.0`
	case DriverMssql:
		return `github.com/denisenkom/go-mssqldb v0.11.0`
	case DriverPostgres:
		return `github.com/lib/pq v1.10.3`
	case DriverSqlite3:
		return `github.com/mattn/go-sqlite3 v1.14.9`
	case DriverOracle:
		return `github.com/godror/godror v0.29.0`
	case "godror":
		return `github.com/godror/godror v0.29.0`
	default:
		return ""
	}
}
