package db

import (
	"io/ioutil"
	"path/filepath"
	"sort"
)

func List(path string) ([]string, error) {
	var names []string
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	files, err := ioutil.ReadDir(absPath)
	if err != nil {
		return names, err
	}
	for _, file := range files {
		if file.IsDir() == false {
			names = append(names, path + "/" + file.Name())
		}
	}
	return names, nil
}

func LoadDatabase(path string) ([]SingleDatabase, error) {
	files, err := List(path)
	if err != nil {
		return nil, err
	}
	m := make(map[string]*SingleDatabase)
	for _, file := range files {
		var db Database
		err = Load(file, &db)
		if err != nil {
			return nil, err
		}
		single, name := IsSingleTable(db)
		if single {
			ex, ok := m[name]
			if !ok {
				m[name] = &SingleDatabase{Filename: file, Table: name, Db: db}
			} else if ex.Filename < file {
				m[name] = &SingleDatabase{Filename: file, Table: name, Db: db}
			}
		}
	}
	arr := make([]SingleDatabase, 0)
	for _, element := range m {
		arr = append(arr, *element)
	}
	sort.Sort(DbByTableName(arr))
	for _, db := range arr {
		Sort(db.Db)
	}
	return  arr, nil
}
type DbByTableName []SingleDatabase

func (a DbByTableName) Len() int      { return len(a) }
func (a DbByTableName) Less(i, j int) bool {
	return a[i].Db.DatabaseChangeLog[0].Changes[0].CreateTable.TableName < a[j].Db.DatabaseChangeLog[0].Changes[0].CreateTable.TableName
}
func (a DbByTableName) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

func SearchDb(search string, a []SingleDatabase) (result int, searchCount int) {
	mid := len(a) / 2
	switch {
	case len(a) == 0:
		result = -1 // not found
	case a[mid].Db.DatabaseChangeLog[0].Changes[0].CreateTable.TableName > search:
		result, searchCount = SearchDb(search, a[:mid])
	case a[mid].Db.DatabaseChangeLog[0].Changes[0].CreateTable.TableName < search:
		result, searchCount = SearchDb(search, a[mid+1:])
		if result >= 0 { // if anything but the -1 "not found" result
			result += mid + 1
		}
	default: // a[mid] == search
		result = mid // found
	}
	searchCount++
	return
}
func GetDb(name string, dbs []SingleDatabase) *SingleDatabase {
	i, _ := SearchDb(name, dbs)
	if i >= 0 {
		return &dbs[i]
	}
	return nil
}
func DiffDataset(t1 []SingleDatabase, t2 []SingleDatabase) DatabaseDiff {
	add := make([]CreateTable, 0)
	drop := make([]CreateTable, 0)
	modify := make([]TableDiff, 0)
	for _, table := range t2 {
		table1 := GetDb(table.Db.DatabaseChangeLog[0].Changes[0].CreateTable.TableName, t1)
		if table1 == nil {
			add = append(add, table.Db.DatabaseChangeLog[0].Changes[0].CreateTable)
		} else {
			if !CompareTables(table1.Db.DatabaseChangeLog[0].Changes[0].CreateTable, table.Db.DatabaseChangeLog[0].Changes[0].CreateTable) {
				diff := DiffTable(table1.Db.DatabaseChangeLog[0].Changes[0].CreateTable, table.Db.DatabaseChangeLog[0].Changes[0].CreateTable)
				modify = append(modify, diff)
			}
		}
	}
	for _, table := range t1 {
		table2 := GetDb(table.Db.DatabaseChangeLog[0].Changes[0].CreateTable.TableName, t2)
		if table2 == nil {
			drop = append(drop, table.Db.DatabaseChangeLog[0].Changes[0].CreateTable)
		}
	}
	return DatabaseDiff{Add: add, Drop: drop, Modify: modify}
}
func DiffDirectories(path1, path2 string) (*DatabaseDiff, error) {
	db1, err := LoadDatabase(path1)
	if err !=  nil {
		return nil, err
	}
	db2, err := LoadDatabase(path2)
	if err !=  nil {
		return nil, err
	}
	d := DiffDataset(db1, db2)
	return &d, nil
}
