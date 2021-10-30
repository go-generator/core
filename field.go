package metadata

type Field struct {
	Name   string `mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Source string `mapstructure:"source" json:"source,omitempty" gorm:"column:source" bson:"source,omitempty" dynamodbav:"source,omitempty" firestore:"source,omitempty"`
	Column string `mapstructure:"column" json:"column,omitempty" gorm:"column:column" bson:"column,omitempty" dynamodbav:"column,omitempty" firestore:"column,omitempty"`
	Type   string `mapstructure:"type" json:"type,omitempty" gorm:"column:type" bson:"type,omitempty" dynamodbav:"type,omitempty" firestore:"type,omitempty"`
	Length int    `mapstructure:"length" json:"length,omitempty" gorm:"column:length" bson:"length,omitempty" dynamodbav:"length,omitempty" firestore:"length,omitempty"`
	Key    bool   `mapstructure:"key" json:"key,omitempty" gorm:"column:key" bson:"key,omitempty" dynamodbav:"key,omitempty" firestore:"key,omitempty"`
}
