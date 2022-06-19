package db

import (
	"fmt"
	"strings"
)

func GenerateFile(diff DatabaseDiff) string {
	addOutput := generateAdd(diff.Add)
	dropOutput := generateDrop(diff.Drop)
	modifyOutput := generateModify(diff.Modify)

	output := fmt.Sprintf(`databaseChangeLog:%v%v%v`, addOutput, dropOutput, modifyOutput)

	return output
}

func generateAdd(add []CreateTable) string {
	columnsVals := make([]string, 0)
	createTableChangeTypes := make([]string, 0)
	var output string

	if add != nil && len(add) > 0 {
		for index, _ := range add {
			columns := add[index].Columns
			for i, v := range columns {
				if columns[i].Constraints != nil {
					if columns[i].Constraints.PrimaryKey {
					col := fmt.Sprintf(`
        - column:
            constraints:
              nullable: %v
              primaryKey: %v
              primaryKeyName: %v
            name: %v
            type: %v`, v.Constraints.Nullable, v.Constraints.PrimaryKey, v.Constraints.PrimaryKeyName, v.Name, v.Type)
					columnsVals = append(columnsVals, col)
					} else {
					col := fmt.Sprintf(`
        - column:
            constraints:
              nullable: %v
            name: %v
            type: %v`, v.Constraints.Nullable, v.Name, v.Type)
					columnsVals = append(columnsVals, col)
					}
				} else {
					col := fmt.Sprintf(`
        - column:
            name: %v
            type: %v`, v.Name, v.Type)
					columnsVals = append(columnsVals, col)
				}
			}
			createTableChangeType := fmt.Sprintf(`
    - createTable:
        columns: %v
        tableName: %v`, strings.Join(columnsVals, ""), add[0].TableName)
			createTableChangeTypes = append(createTableChangeTypes, createTableChangeType)

			output = fmt.Sprintf(`
- changeSet:
    id: CREATE_TABLE_%v
    author: test.user (generated)
    labels: sit-init
    changes: %v`, add[index].TableName, strings.Join(createTableChangeTypes, ""))
		}
	} else {
		return ""
	}

	return output
}

func generateDrop(drop []CreateTable) string {
	dropTableChangeTypes := make([]string, 0)
	var output string

	if drop != nil && len(drop) > 0 {
		for i, v := range drop {
			dropStmt := fmt.Sprintf(`
    - dropTable:
        cascadeConstraints: true
        tableName: %v`, v.TableName)
			dropTableChangeTypes = append(dropTableChangeTypes, dropStmt)

			output = fmt.Sprintf(`
- changeSet:
    id: DROP_TABLE_%v
    author: test.user (generated)
    labels: sit-drop
    changes: %v`, drop[i].TableName, strings.Join(dropTableChangeTypes, ""))
		}
	} else {
		return ""
	}

	return output
}

func generateModify(modify []TableDiff) string {
	modifyTableChangeTypes := make([]string, 0)
	var output string

	if modify != nil && len(modify) > 0 {
		for i, vMod := range modify {
			for _, vAdd := range vMod.Add {
				addColumnStmt := fmt.Sprintf(`
    - addColumn:
        tableName: %v
        columns:
        - column:
            name: %v
            type: %v`, vMod.TableName, vAdd.Column.Name, vAdd.Column.Type)
				modifyTableChangeTypes = append(modifyTableChangeTypes, addColumnStmt)
			}
			for _, vDrop := range vMod.Drop {
				dropColumnStmt := fmt.Sprintf(`
    - dropColumn:
        columnName: %v
        tableName: %v`, vDrop.Column.Name, vMod.TableName)
				modifyTableChangeTypes = append(modifyTableChangeTypes, dropColumnStmt)
			}
			for _, vModify := range vMod.Modify {
				dropColumnStmt := fmt.Sprintf(`
    - modifyDataType:
        columnName: %v
        newDataType: %v
        tableName: %v`, vModify.Column.Name, vModify.Column.Type, vMod.TableName)
				modifyTableChangeTypes = append(modifyTableChangeTypes, dropColumnStmt)
			}
			output = fmt.Sprintf(`
- changeSet:
    id: MODIFY_TABLE_%v
    author: test.user (generated)
    labels: sit-change
    changes: %v`, modify[i].TableName, strings.Join(modifyTableChangeTypes, ""))
		}
	} else {
		return ""
	}
	return output
}