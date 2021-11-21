package metadata

type Config struct {
	Project     string `mapstructure:"project" json:"project,omitempty" gorm:"column:project" bson:"project,omitempty" dynamodbav:"project,omitempty" firestore:"project,omitempty"`
	ProjectName string `mapstructure:"project_name" json:"projectName,omitempty" gorm:"column:projectname" bson:"projectName,omitempty" dynamodbav:"projectName,omitempty" firestore:"projectName,omitempty"`
	ProjectPath string `mapstructure:"project_path" json:"projectPath,omitempty" gorm:"column:projectpath" bson:"projectPath,omitempty" dynamodbav:"projectPath,omitempty" firestore:"projectPath,omitempty"`
	Template    string `mapstructure:"template" json:"template,omitempty" gorm:"column:template" bson:"template,omitempty" dynamodbav:"template,omitempty" firestore:"template,omitempty"` // Template Path
	Cache       string `mapstructure:"cache" json:"cache,omitempty" gorm:"column:cache" bson:"cache,omitempty" dynamodbav:"cache,omitempty" firestore:"cache,omitempty"`
	DBCache     string `mapstructure:"dbcache" json:"dbCache,omitempty" gorm:"column:dbcache" bson:"dbCache,omitempty" dynamodbav:"dbCache,omitempty" firestore:"dbCache,omitempty"`
}
