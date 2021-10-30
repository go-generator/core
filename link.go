package metadata

type Link struct {
	Column string `mapstructure:"column" json:"column,omitempty" gorm:"column:column" bson:"column,omitempty" dynamodbav:"column,omitempty" firestore:"column,omitempty"`
	To     string `mapstructure:"to" json:"to,omitempty" gorm:"column:to" bson:"to,omitempty" dynamodbav:"to,omitempty" firestore:"to,omitempty"`
}
