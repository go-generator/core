package core

type DatabaseConfig struct {
	DSN      string `yaml:"dsn" mapstructure:"dsn" json:"dsn,omitempty" gorm:"column:dsn" bson:"dsn,omitempty" dynamodbav:"dsn,omitempty" firestore:"dsn,omitempty"`
	Driver   string `yaml:"driver" mapstructure:"driver" json:"driver,omitempty" gorm:"column:driver" bson:"driver,omitempty" dynamodbav:"driver,omitempty" firestore:"driver,omitempty"`
	Host     string `yaml:"host" mapstructure:"host" json:"host,omitempty" gorm:"column:host" bson:"host,omitempty" dynamodbav:"host,omitempty" firestore:"host,omitempty"`
	Port     int64  `yaml:"port" mapstructure:"port" json:"port,omitempty" gorm:"column:port" bson:"port,omitempty" dynamodbav:"port,omitempty" firestore:"port,omitempty"`
	Database string `yaml:"database" mapstructure:"database" json:"database,omitempty" gorm:"column:database" bson:"database,omitempty" dynamodbav:"database,omitempty" firestore:"database,omitempty"`
	User     string `yaml:"user" mapstructure:"user" json:"user,omitempty" gorm:"column:user" bson:"user,omitempty" dynamodbav:"user,omitempty" firestore:"user,omitempty"`
	Password string `yaml:"password" mapstructure:"password" json:"password,omitempty" gorm:"column:password" bson:"password,omitempty" dynamodbav:"password,omitempty" firestore:"password,omitempty"`
}

type Database struct {
	MySql    string `yaml:"mysql" mapstructure:"mysql" json:"mysql,omitempty" gorm:"column:mysql" bson:"mysql,omitempty" dynamodbav:"mysql,omitempty" firestore:"mysql,omitempty"`
	Postgres string `yaml:"postgres" mapstructure:"postgres" json:"postgres,omitempty" gorm:"column:postgres" bson:"postgres,omitempty" dynamodbav:"postgres,omitempty" firestore:"postgres,omitempty"`
	Sqlite3  string `yaml:"sqlite3" mapstructure:"sqlite3" json:"sqlite3,omitempty" gorm:"column:sqlite3" bson:"sqlite3,omitempty" dynamodbav:"sqlite3,omitempty" firestore:"sqlite3,omitempty"`
	Mssql    string `yaml:"mssql" mapstructure:"mssql" json:"mssql,omitempty" gorm:"column:mssql" bson:"mssql,omitempty" dynamodbav:"mssql,omitempty" firestore:"mssql,omitempty"`
	Oracle   string `yaml:"oracle" mapstructure:"oracle" json:"oracle,omitempty" gorm:"column:oracle" bson:"oracle,omitempty" dynamodbav:"oracle,omitempty" firestore:"oracle,omitempty"`
}
