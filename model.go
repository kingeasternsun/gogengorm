package main

import (
	"os"
	"path/filepath"
	"sort"
	"text/template"
)

type ModelType struct {
	PackageName   string
	ModelName     string
	Columns       []ColumnType
	UniqueColumns []UniqueColumn //TODO if detect unique index, we should generate different gorm code.
}

// extract the unique index columns from the normal columns
func extractUniqueIndex(models []ModelType) {

	for i := range models {

		var normalColumns []ColumnType
		uniqueIndexColumns := make(map[string][]ColumnType)

		for j := range models[i].Columns {
			if models[i].Columns[j].UniqueIndexName == "" {
				normalColumns = append(normalColumns, models[i].Columns[j])
			} else {
				uniqueIndexColumns[models[i].Columns[j].UniqueIndexName] = append(uniqueIndexColumns[models[i].Columns[j].UniqueIndexName], models[i].Columns[j])
			}
		}

		models[i].Columns = normalColumns
		for k, v := range uniqueIndexColumns {
			models[i].UniqueColumns = append(models[i].UniqueColumns, UniqueColumn{IndexName: k, Columns: v})
		}

		sort.Slice(models[i].UniqueColumns, func(j, k int) bool {
			return models[i].UniqueColumns[j].IndexName > models[i].UniqueColumns[k].IndexName
		})

	}

}

func generateGormCode(models []ModelType, fileDir string) error {

	tpl, err := template.ParseFiles("./gorm.template")
	if err != nil {
		return err
	}

	for _, model := range models {
		name, _ := convertFieldName("snakecase", model.ModelName)
		goPath := filepath.Join(fileDir, name+"_db.go")
		rd, err := os.Create(goPath)
		if err != nil {
			return err
		}

		err = tpl.ExecuteTemplate(rd, "gorm.template", model)
		if err != nil {
			rd.Close()
			return err
		}
		rd.Close()
	}

	return nil

}
