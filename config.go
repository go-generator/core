package metadata

type Config struct {
	ProjectName string `mapstructure:"project_name" json:"projectName,omitempty" gorm:"column:projectname" bson:"projectName,omitempty" dynamodbav:"projectName,omitempty" firestore:"projectName,omitempty"`
	Template    string `mapstructure:"template" json:"template,omitempty" gorm:"column:template" bson:"template,omitempty" dynamodbav:"template,omitempty" firestore:"template,omitempty"`
	Types       string `mapstructure:"types" json:"types,omitempty" gorm:"column:types" bson:"types,omitempty" dynamodbav:"types,omitempty" firestore:"types,omitempty"`
	Projects    string `mapstructure:"projects" json:"projects,omitempty" gorm:"column:projects" bson:"projects,omitempty" dynamodbav:"projects,omitempty" firestore:"projects,omitempty"`
	Project     string `mapstructure:"project" json:"project,omitempty" gorm:"column:project" bson:"project,omitempty" dynamodbav:"project,omitempty" firestore:"project,omitempty"`
}
