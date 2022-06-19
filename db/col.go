package db

func GetColumn(name string, cols []ColumnType) *ColumnType {
	i, _ := SearchColumn(name, cols)
	if i >= 0 {
		return &cols[i]
	}
	return nil
}
func SearchColumn(search string, a []ColumnType) (result int, searchCount int) {
	mid := len(a) / 2
	switch {
	case len(a) == 0:
		result = -1 // not found
	case a[mid].Name > search:
		result, searchCount = SearchColumn(search, a[:mid])
	case a[mid].Name < search:
		result, searchCount = SearchColumn(search, a[mid+1:])
		if result >= 0 { // if anything but the -1 "not found" result
			result += mid + 1
		}
	default: // a[mid] == search
		result = mid // found
	}
	searchCount++
	return
}

type ByColumnName []ColumnType

func (a ByColumnName) Len() int           { return len(a) }
func (a ByColumnName) Less(i, j int) bool { return a[i].Name < a[j].Name }
func (a ByColumnName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type ByColumnNo []ColumnType

func (a ByColumnNo) Len() int           { return len(a) }
func (a ByColumnNo) Less(i, j int) bool { return a[i].No < a[j].No }
func (a ByColumnNo) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type TableDiff struct {
	TableName string
	Add    []ColumnType `yaml:"add" mapstructure:"add"`
	Drop   []ColumnType `yaml:"drop" mapstructure:"drop"`
	Modify []ColumnType `yaml:"modify" mapstructure:"modify"`
}

func DiffTable(t1 CreateTable, t2 CreateTable) TableDiff {
	add := make([]ColumnType, 0)
	drop := make([]ColumnType, 0)
	modify := make([]ColumnType, 0)
	for _, col := range t2.Columns {
		col1 := GetColumn(col.Name, t1.Columns)
		if col1 == nil {
			add = append(add, col)
		} else {
			if !CompareColumns(*col1, col) {
				modify = append(modify, col)
			}
		}
	}
	for _, col := range t1.Columns {
		col2 := GetColumn(col.Name, t2.Columns)
		if col2 == nil {
			drop = append(drop, col)
		}
	}
	return TableDiff{TableName: t2.TableName, Add: add, Drop: drop, Modify: modify}
}

func CompareColumns(col1 ColumnType, col2 ColumnType) bool {
	if col1.Name != col2.Name {
		return false
	}
	if col1.Type != col2.Type {
		return false
	}
	if col1.Constraints != nil && col2.Constraints != nil {
		if col1.Constraints.Nullable != col2.Constraints.Nullable {
			return false
		}
		if col1.Constraints.PrimaryKey != col2.Constraints.PrimaryKey {
			return false
		}
		if col1.Constraints.PrimaryKeyName != col2.Constraints.PrimaryKeyName {
			return false
		}
	}
	return true
}
