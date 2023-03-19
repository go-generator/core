package core

type Relationship struct {
	Ref    string `yaml:"table" mapstructure:"table" json:"table,omitempty" gorm:"column:table" bson:"table,omitempty" dynamodbav:"table,omitempty" firestore:"table,omitempty"`
	Model  string `yaml:"model" mapstructure:"model" json:"model,omitempty" gorm:"column:model" bson:"model,omitempty" dynamodbav:"model,omitempty" firestore:"model,omitempty"`
	Fields []Link `yaml:"fields" mapstructure:"fields" json:"fields,omitempty" gorm:"column:fields" bson:"fields,omitempty" dynamodbav:"fields,omitempty" firestore:"fields,omitempty"`
}
