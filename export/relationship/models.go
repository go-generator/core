package relationship

const (
	OneToOne    = "one to one"
	OneToMany   = "one to many"
	ManyToOne   = "many to one"
	ManyToMany  = "many to many"
	Unknown     = "unknown"
	Unsupported = "unsupported"
)

type Tables struct {
	Name string `gorm:"column:table"`
}

type RelTables struct {
	Table            string `gorm:"column:table"`
	Column           string `gorm:"column:column"`
	ReferencedTable  string `gorm:"column:referenced_table"`
	ReferencedColumn string `gorm:"column:referenced_column"`
	Relationship     string
}

type SqliteRel struct {
	Sql string `gorm:"column:sql"`
}

type MySqlUnique struct {
	Column    string `gorm:"column:Column_name"`
	NonUnique bool   `gorm:"column:Non_unique"` // False mean it's unique, True means it can contain duplicate
	Key       string `gorm:"column:Key_name"`
}

type PostgresUnique struct {
	Table string `gorm:"column:table"`
	Index string `gorm:"column:index"`
}

type MssqlUnique struct {
	Column     string `gorm:"column:column"`
	Constraint string `gorm:"column:constraint"`
}

type SqliteUnique struct {
	Name   string `gorm:"column:name"`
	Unique string `gorm:"column:unique"`
	Origin string `gorm:"column:origin"`
}

type SqliteNotNull struct {
	Name    string `gorm:"column:name"`
	NotNull bool   `gorm:"column:notnull"`
}

type PrimaryKey struct {
	Column string `gorm:"column:column"`
}
