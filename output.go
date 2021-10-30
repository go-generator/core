package metadata

type FileInfo struct {
	Name       string
	StructName string
	Fields     []FieldInfo
	IDFields   []FieldInfo
}

type FieldInfo struct {
	Name string
	Type string
}

type Output struct {
	Directory string `mapstructure:"directory" json:"directory,omitempty" gorm:"column:directory" bson:"directory,omitempty" dynamodbav:"directory,omitempty" firestore:"directory,omitempty"`
	Path      string `mapstructure:"path" json:"path,omitempty" gorm:"column:path" bson:"path,omitempty" dynamodbav:"path,omitempty" firestore:"path,omitempty"`
	Files     []File `mapstructure:"files" json:"files,omitempty" gorm:"column:files" bson:"files,omitempty" dynamodbav:"files,omitempty" firestore:"files,omitempty"`
	OutFile   []FileInfo
}
