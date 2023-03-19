package metadata

type Link struct {
	Column string `yaml:"column" mapstructure:"column" json:"column,omitempty" gorm:"column:column" bson:"column,omitempty" dynamodbav:"column,omitempty" firestore:"column,omitempty"`
	To     string `yaml:"to" mapstructure:"to" json:"to,omitempty" gorm:"column:to" bson:"to,omitempty" dynamodbav:"to,omitempty" firestore:"to,omitempty"`
}
