package metadata

type Field struct {
	Name       string `mapstructure:"name" json:"name,omitempty" gorm:"column:name" bson:"name,omitempty" dynamodbav:"name,omitempty" firestore:"name,omitempty"`
	Source     string `mapstructure:"source" json:"source,omitempty" gorm:"column:source" bson:"source,omitempty" dynamodbav:"source,omitempty" firestore:"source,omitempty"`
	Column     string `mapstructure:"column" json:"column,omitempty" gorm:"column:column" bson:"column,omitempty" dynamodbav:"column,omitempty" firestore:"column,omitempty"`
	Type       string `mapstructure:"type" json:"type,omitempty" gorm:"column:type" bson:"type,omitempty" dynamodbav:"type,omitempty" firestore:"type,omitempty"`
	DbType     string `mapstructure:"db_type" json:"dbType,omitempty" gorm:"column:dbtype" bson:"dbType,omitempty" dynamodbav:"dbType,omitempty" firestore:"dbType,omitempty"`
	FullDbType string `mapstructure:"full_db_type" json:"fullDbType,omitempty" gorm:"column:fulldbtype" bson:"fullDbType,omitempty" dynamodbav:"fullDbType,omitempty" firestore:"fullDbType,omitempty"`
	Required   bool   `mapstructure:"required" json:"required,omitempty" gorm:"column:required" bson:"required,omitempty" dynamodbav:"required,omitempty" firestore:"required,omitempty"`
	Length     int    `mapstructure:"length" json:"length,omitempty" gorm:"column:length" bson:"length,omitempty" dynamodbav:"length,omitempty" firestore:"length,omitempty"`
	Key        bool   `mapstructure:"key" json:"key,omitempty" gorm:"column:key" bson:"key,omitempty" dynamodbav:"key,omitempty" firestore:"key,omitempty"`
	KeyName    string `mapstructure:"key_name" json:"keyName,omitempty" gorm:"column:keyname" bson:"keyName,omitempty" dynamodbav:"keyName,omitempty" firestore:"keyName,omitempty"`
	Precision  *int   `mapstructure:"precision" json:"precision,omitempty" gorm:"column:precision" bson:"precision,omitempty" dynamodbav:"precision,omitempty" firestore:"precision,omitempty"`
	Scale      *int   `mapstructure:"scale" json:"scale,omitempty" gorm:"column:scale" bson:"scale,omitempty" dynamodbav:"scale,omitempty" firestore:"scale,omitempty"`
}
