package metadata

type Template struct {
	Name    string `mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	File    string `mapstructure:"file" json:"file,omitempty" gorm:"column:file" bson:"file,omitempty" dynamodbav:"file,omitempty" firestore:"file,omitempty"`
	Model   bool   `mapstructure:"model" json:"model,omitempty" gorm:"column:model" bson:"model,omitempty" dynamodbav:"model,omitempty" firestore:"model,omitempty"`
	Replace bool   `mapstructure:"replace" json:"replace,omitempty" gorm:"column:replace" bson:"replace,omitempty" dynamodbav:"replace,omitempty" firestore:"replace,omitempty"`
}
