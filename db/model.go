package db

type SingleDatabase struct {
	Filename string   `yaml:"filename" mapstructure:"filename"`
	Table    string   `yaml:"table" mapstructure:"table"`
	Db       Database `yaml:"db" mapstructure:"db"`
}
type Database struct {
	DatabaseChangeLog []ChangeSetList `yaml:"databaseChangeLog" mapstructure:"databaseChangeLog"`
}

type ChangeSetList struct {
	ChangeSet `yaml:"changeSet"`
}
type ChangeSet struct {
	ID      string    `yaml:"id" mapstructure:"id"`
	Author  string    `yaml:"author" mapstructure:"author"`
	Labels  string    `yaml:"labels" mapstructure:"labels"`
	Changes []Changes `yaml:"changes" mapstructure:"changes"`
}

type Changes struct {
	CreateTable CreateTable `yaml:"createTable" mapstructure:"createTable"`
}

type CreateTable struct {
	Columns   []ColumnType `yaml:"columns" mapstructure:"columns"`
	TableName string       `yaml:"tableName" mapstructure:"tableName"`
}

type ColumnType struct {
	Column `yaml:"column"`
}
type Column struct {
	Constraints *Constraints `yaml:"constraints,omitempty" mapstructure:"constraints"`
	Name        string       `yaml:"name" mapstructure:"name"`
	Type        string       `yaml:"type" mapstructure:"type"`
	No          int          `yaml:"no" mapstructure:"no"`
}

type Constraints struct {
	Nullable       bool   `yaml:"nullable,omitempty" mapstructure:"nullable"`
	PrimaryKey     bool   `yaml:"primaryKey,omitempty" mapstructure:"primaryKey"`
	PrimaryKeyName string `yaml:"primaryKeyName,omitempty" mapstructure:"primaryKeyName" `
}
