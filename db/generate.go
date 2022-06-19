package db

import (
	"fmt"
	"strings"
)

func GenerateDiff(diff DatabaseDiff) string {
	addOutput := generateAdd(diff.Add)
	dropOutput := generateDrop(diff.Drop)
	modifyOutput := generateModify(diff.Modify)

	output := fmt.Sprintf(`databaseChangeLog:%v%v%v`, addOutput, dropOutput, modifyOutput)

	return output
}

func generateAdd(add []CreateTable) string {
	columnsVals := make([]string, 0)
	createTableChangeTypes := make([]string, 0)
	outputs := make([]string, 0)

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

			output := fmt.Sprintf(`
- changeSet:
    id: CREATE_TABLE_%v
    author: test.user (generated)
    labels: sit-init
    changes: %v`, add[index].TableName, createTableChangeTypes[index])
			outputs = append(outputs, output)
		}
	} else {
		return ""
	}

	return strings.Join(outputs, "")
}

func generateDrop(drop []CreateTable) string {
	dropTableChangeTypes := make([]string, 0)
	outputs := make([]string, 0)

	if drop != nil && len(drop) > 0 {
		for i, v := range drop {
			dropStmt := fmt.Sprintf(`
    - dropTable:
        cascadeConstraints: true
        tableName: %v`, v.TableName)
			dropTableChangeTypes = append(dropTableChangeTypes, dropStmt)

			output := fmt.Sprintf(`
- changeSet:
    id: DROP_TABLE_%v
    author: test.user (generated)
    labels: sit-drop
    changes: %v`, drop[i].TableName, dropTableChangeTypes[i])
			outputs = append(outputs, output)
		}
	} else {
		return ""
	}

	return strings.Join(outputs, "")
}

func generateModify(modify []TableDiff) string {
	outputs := make([]string, 0)

	if modify != nil && len(modify) > 0 {
		for i, vMod := range modify {
			count := 0
			modifyTableChangeTypes := make([]string, 0)
			for _, vAdd := range vMod.Add {
				addColumnStmt := fmt.Sprintf(`
    - addColumn:
        tableName: %v
        columns:
        - column:
            name: %v
            type: %v`, vMod.TableName, vAdd.Column.Name, vAdd.Column.Type)
				modifyTableChangeTypes = append(modifyTableChangeTypes, addColumnStmt)
				count++
			}
			for _, vDrop := range vMod.Drop {
				dropColumnStmt := fmt.Sprintf(`
    - dropColumn:
        columnName: %v
        tableName: %v`, vDrop.Column.Name, vMod.TableName)
				modifyTableChangeTypes = append(modifyTableChangeTypes, dropColumnStmt)
				count++
			}
			for _, vModify := range vMod.Modify {
				modifyColumnStmt := fmt.Sprintf(`
    - modifyDataType:
        columnName: %v
        newDataType: %v
        tableName: %v`, vModify.Column.Name, vModify.Column.Type, vMod.TableName)
				modifyTableChangeTypes = append(modifyTableChangeTypes, modifyColumnStmt)
				count++
			}
			for y := 0; y < count; y++ {
				output := fmt.Sprintf(`
- changeSet:
    id: MODIFY_TABLE_%v_%v
    author: test.user (generated)
    labels: sit-change
    changes: %v`, modify[i].TableName, y, modifyTableChangeTypes[y])

				outputs = append(outputs, output)
			}
		}
	} else {
		return ""
	}

	return strings.Join(outputs, "")
}
