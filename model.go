package metadata

import "strings"

type Model struct {
	Name      string          `mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Source    string          `mapstructure:"source" json:"source,omitempty" gorm:"column:source" bson:"source,omitempty" dynamodbav:"source,omitempty" firestore:"source,omitempty"`
	Table     string          `mapstructure:"table" json:"table,omitempty" gorm:"column:table" bson:"table,omitempty" dynamodbav:"table,omitempty" firestore:"table,omitempty"`
	Alias     []TypeAlias     `mapstructure:"alias" json:"alias,omitempty" gorm:"column:alias" bson:"alias,omitempty" dynamodbav:"alias,omitempty" firestore:"alias,omitempty"`
	Ones      []Relationship  `mapstructure:"ones" json:"ones,omitempty" gorm:"column:ones" bson:"ones,omitempty" dynamodbav:"ones,omitempty" firestore:"ones,omitempty"`
	Models    []Relationship  `mapstructure:"models" json:"models,omitempty" gorm:"column:models" bson:"models,omitempty" dynamodbav:"models,omitempty" firestore:"models,omitempty"`
	Arrays    []Relationship  `mapstructure:"arrays" json:"arrays,omitempty" gorm:"column:arrays" bson:"arrays,omitempty" dynamodbav:"arrays,omitempty" firestore:"arrays,omitempty"`
	Fields    []Field         `mapstructure:"fields" json:"fields,omitempty" gorm:"column:fields" bson:"fields,omitempty" dynamodbav:"fields,omitempty" firestore:"fields,omitempty"`
	WriteFile strings.Builder `json:"-" bson:"-" dynamodbav:"-" firestore:"-"`
}
