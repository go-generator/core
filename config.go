package metadata

type Config struct {
	Prefix      string `yaml:"prefix" mapstructure:"prefix" json:"prefix,omitempty" gorm:"column:prefix" bson:"prefix,omitempty" dynamodbav:"prefix,omitempty" firestore:"prefix,omitempty"`
	Suffix      string `yaml:"suffix" mapstructure:"suffix" json:"suffix,omitempty" gorm:"column:suffix" bson:"suffix,omitempty" dynamodbav:"suffix,omitempty" firestore:"suffix,omitempty"`
	Project     string `yaml:"project" mapstructure:"project" json:"project,omitempty" gorm:"column:project" bson:"project,omitempty" dynamodbav:"project,omitempty" firestore:"project,omitempty"`
	ProjectName string `yaml:"project_name" mapstructure:"project_name" json:"projectName,omitempty" gorm:"column:projectname" bson:"projectName,omitempty" dynamodbav:"projectName,omitempty" firestore:"projectName,omitempty"`
	ProjectPath string `yaml:"project_path" mapstructure:"project_path" json:"projectPath,omitempty" gorm:"column:projectpath" bson:"projectPath,omitempty" dynamodbav:"projectPath,omitempty" firestore:"projectPath,omitempty"`
	Template    string `yaml:"template" mapstructure:"template" json:"template,omitempty" gorm:"column:template" bson:"template,omitempty" dynamodbav:"template,omitempty" firestore:"template,omitempty"` // Template Path
	DBCache     string `yaml:"dbcache" mapstructure:"dbcache" json:"dbCache,omitempty" gorm:"column:dbcache" bson:"dbCache,omitempty" dynamodbav:"dbCache,omitempty" firestore:"dbCache,omitempty"`
	DB          string `yaml:"db" mapstructure:"db" json:"db,omitempty" gorm:"column:db" bson:"db,omitempty" dynamodbav:"db,omitempty" firestore:"db,omitempty"`
}
