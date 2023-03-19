package template

import "time"

type Template struct {
	Id        string `yaml:"id" mapstructure:"id" json:"id" gorm:"column:id;primary_key" bson:"_id" dynamodbav:"id" firestore:"-" avro:"id" validate:"max=40" match:"equal"`
	Content   string `yaml:"content" mapstructure:"content" json:"content" gorm:"column:content" bson:"content" dynamodbav:"content" firestore:"content" avro:"content"`
	Data      map[string]string
	UpdatedAt *time.Time `yaml:"updated_at" mapstructure:"updated_at" json:"updatedAt,omitempty" gorm:"column:updated_at" bson:"updatedAt" dynamodbav:"updateDate" firestore:"updatedAt" avro:"updatedAt"`
}

func (t *Template) SetData(templates []Template) {
	var m = make(map[string]string, 0)
	for _, v := range templates {
		m[v.Id] = v.Content
	}
	t.Data = m
}

type FileTemplate struct {
	FileNames []string `yaml:"id" mapstructure:"id" json:"id"`
}

type FileContentTemplate struct {
	FileName string `yaml:"filename" mapstructure:"filename" json:"fileName"`
	Content  string `yaml:"content" mapstructure:"content" json:"content"`
}
