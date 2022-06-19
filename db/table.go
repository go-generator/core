package db

import "sort"

func IsSingleTable(db Database) (bool, string) {
	if len(db.DatabaseChangeLog) == 1 {
		if len(db.DatabaseChangeLog[0].Changes) == 1 {
			return true, db.DatabaseChangeLog[0].Changes[0].CreateTable.TableName
		}
	}
	return false, ""
}
func GetTable(name string, tables []ChangeSetList) *CreateTable {
	i, _ := SearchTable(name, tables)
	if i >= 0 {
		return &tables[i].Changes[0].CreateTable
	}
	return nil
}
func SearchTable(search string, a []ChangeSetList) (result int, searchCount int) {
	mid := len(a) / 2
	switch {
	case len(a) == 0:
		result = -1 // not found
	case a[mid].Changes[0].CreateTable.TableName > search:
		result, searchCount = SearchTable(search, a[:mid])
	case a[mid].Changes[0].CreateTable.TableName < search:
		result, searchCount = SearchTable(search, a[mid+1:])
		if result >= 0 { // if anything but the -1 "not found" result
			result += mid + 1
		}
	default: // a[mid] == search
		result = mid // found
	}
	searchCount++
	return
}

func Sort(data Database) {
	sort.Sort(ByTableName(data.DatabaseChangeLog))
	for _, change := range data.DatabaseChangeLog {
		if len(change.Changes) > 0 {
			sort.Sort(ByColumnName(change.Changes[0].CreateTable.Columns))
		}
	}
}

type ByTableName []ChangeSetList

func (a ByTableName) Len() int      { return len(a) }
func (a ByTableName) Less(i, j int) bool {
	return a[i].Changes[0].CreateTable.TableName < a[j].Changes[0].CreateTable.TableName
}
func (a ByTableName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type DatabaseDiff struct {
	Add    []CreateTable `yaml:"add" mapstructure:"add"`
	Drop   []CreateTable `yaml:"drop" mapstructure:"drop"`
	Modify []TableDiff   `yaml:"modify" mapstructure:"modify"`
}

func DiffDatabase(t1 Database, t2 Database) DatabaseDiff {
	add := make([]CreateTable, 0)
	drop := make([]CreateTable, 0)
	modify := make([]TableDiff, 0)
	for _, table := range t2.DatabaseChangeLog {
		table1 := GetTable(table.Changes[0].CreateTable.TableName, t1.DatabaseChangeLog)
		if table1 == nil {
			add = append(add, table.Changes[0].CreateTable)
		} else {
			if !CompareTables(*table1, table.Changes[0].CreateTable) {
				diff := DiffTable(*table1, table.Changes[0].CreateTable)
				modify = append(modify, diff)
			}
		}
	}
	for _, table := range t1.DatabaseChangeLog {
		table2 := GetTable(table.Changes[0].CreateTable.TableName, t2.DatabaseChangeLog)
		if table2 == nil {
			drop = append(drop, table.Changes[0].CreateTable)
		}
	}
	return DatabaseDiff{Add: add, Drop: drop, Modify: modify}
}

func CompareTables(table1 CreateTable, table2 CreateTable) bool {
	l := len(table1.Columns)
	if l != len(table2.Columns) {
		return false
	}
	for i := 0; i < l; i++ {
		e := CompareColumns(table1.Columns[i], table2.Columns[i])
		if e == false {
			return e
		}
	}
	return true
}
