package metadata

type Config struct {
	ProjectName  string `mapstructure:"projectName" json:"projectName,omitempty" gorm:"column:project_name" bson:"projectName,omitempty" dynamodbav:"projectName,omitempty" firestore:"projectName,omitempty"`
	TemplatePath string `mapstructure:"template" json:"template,omitempty" gorm:"column:template" bson:"template,omitempty" dynamodbav:"template,omitempty" firestore:"template,omitempty"`
	PrjTmplPath  string `mapstructure:"prjTmplPath" json:"prjTmplPath,omitempty" gorm:"column:project_tmpl_path" bson:"prjTmplPath,omitempty" dynamodbav:"prjTmplPath,omitempty" firestore:"prjTmplPath,omitempty"`
	PrjTmplName  string `mapstructure:"prjTmplName" json:"prjTmplName,omitempty" gorm:"column:project_tmpl_name" bson:"prjTmplName,omitempty" dynamodbav:"prjTmplName,omitempty" firestore:"prjTmplName,omitempty"`
	Cache        string `mapstructure:"cache" json:"cache,omitempty" gorm:"column:cache" bson:"cache,omitempty" dynamodbav:"cache,omitempty" firestore:"cache,omitempty"`
	DBCache      string `mapstructure:"dbCache" json:"dbCache,omitempty" gorm:"column:db_cache" bson:"dbCache,omitempty" dynamodbav:"dbCache,omitempty" firestore:"dbCache,omitempty"`
}
