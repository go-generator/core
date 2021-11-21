package metadata

type Relationship struct {
	Ref    string `mapstructure:"table" json:"table,omitempty" gorm:"column:table" bson:"table,omitempty" dynamodbav:"table,omitempty" firestore:"table,omitempty"`
	Model  string `mapstructure:"model" json:"model,omitempty" gorm:"column:model" bson:"model,omitempty" dynamodbav:"model,omitempty" firestore:"model,omitempty"`
	Fields []Link `mapstructure:"fields" json:"fields,omitempty" gorm:"column:fields" bson:"fields,omitempty" dynamodbav:"fields,omitempty" firestore:"fields,omitempty"`
}
