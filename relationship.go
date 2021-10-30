package metadata

type Relationship struct {
	Ref    string `mapstructure:"table" json:"table,omitempty" gorm:"column:table" bson:"table,omitempty" dynamodbav:"table,omitempty" firestore:"table,omitempty"`
	Fields []Link `mapstructure:"fields" json:"fields,omitempty" gorm:"column:fields" bson:"fields,omitempty" dynamodbav:"fields,omitempty" firestore:"fields,omitempty"`
}
