package metadata

type Project struct {
	Root       string            `mapstructure:"root" json:"root,omitempty" gorm:"column:root" bson:"root,omitempty" dynamodbav:"root,omitempty" firestore:"root,omitempty"`
	Language   string            `mapstructure:"language" json:"language,omitempty" gorm:"column:language" bson:"language,omitempty" dynamodbav:"language,omitempty" firestore:"language,omitempty"`
	Env        map[string]string `mapstructure:"env" json:"env,omitempty" gorm:"column:env" bson:"env,omitempty" dynamodbav:"env,omitempty" firestore:"env,omitempty"`
	Statics    []Template        `mapstructure:"statics" json:"statics,omitempty" gorm:"column:statics" bson:"statics,omitempty" dynamodbav:"statics,omitempty" firestore:"statics,omitempty"`
	Collection []string          `mapstructure:"collection" json:"collection,omitempty" gorm:"column:collection" bson:"collection,omitempty" dynamodbav:"collection,omitempty" firestore:"collection,omitempty"`
	Arrays     []Template        `mapstructure:"arrays" json:"arrays,omitempty" gorm:"column:arrays" bson:"arrays,omitempty" dynamodbav:"arrays,omitempty" firestore:"arrays,omitempty"`
	Entities   []Template        `mapstructure:"entities" json:"entities,omitempty" gorm:"column:entities" bson:"entities,omitempty" dynamodbav:"entities,omitempty" firestore:"entities,omitempty"`
	TypesFile  string            `mapstructure:"types_file" json:"typesFile,omitempty" gorm:"column:typesfile" bson:"typesFile,omitempty" dynamodbav:"typesFile,omitempty" firestore:"typesFile,omitempty"`
	Types      map[string]string `mapstructure:"types" json:"types,omitempty" gorm:"column:types" bson:"types,omitempty" dynamodbav:"types,omitempty" firestore:"types,omitempty"`
	ModelsFile string            `mapstructure:"models_file" json:"modelsFile,omitempty" gorm:"column:modelsfile" bson:"modelsFile,omitempty" dynamodbav:"modelsFile,omitempty" firestore:"modelsFile,omitempty"`
	Models     []Model           `mapstructure:"models" json:"models,omitempty" gorm:"column:models" bson:"models,omitempty" dynamodbav:"models,omitempty" firestore:"models,omitempty"`
}
