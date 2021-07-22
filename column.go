package main

import "strings"

// ColumnType contains a FieldName and real column name
type ColumnType struct {
	FieldName string
	// the column name in database
	Column string
	// the golang type of field
	GoType string
	// the name of arg
	VarName string
	// the unique index name, include primary
	UniqueIndexName string
	// the index name, without unique index
	IndexName string
}

type UniqueColumn struct {
	Columns   []ColumnType // support combine unique index
	IndexName string       //index name
}

func (col UniqueColumn) IsCompositeIndex() bool {
	return len(col.Columns) > 1
}

// generate function Sufix
func (col UniqueColumn) FunctionSufix() string {
	res := strings.Builder{}
	for i, column := range col.Columns {
		res.WriteString(column.FieldName)
		if i != len(col.Columns)-1 {
			res.WriteString("And")
		}
	}
	return res.String()
}

// generate function args
func (col UniqueColumn) FunctionArgs() string {
	res := strings.Builder{}
	for _, column := range col.Columns {
		res.WriteString(" , ")
		res.WriteString(column.VarName)
		res.WriteString(" ")
		res.WriteString(column.GoType)
	}
	return res.String()
}

// generate where conditon
func (col UniqueColumn) WhereCondition() string {
	res := strings.Builder{}
	for i, column := range col.Columns {
		res.WriteString(column.Column)
		res.WriteString(" = ? ")
		if i != len(col.Columns)-1 {
			res.WriteString("And ")
		}
	}
	return res.String()
}

// generate where args
func (col UniqueColumn) WhereArgs() string {
	res := strings.Builder{}
	for i, column := range col.Columns {
		res.WriteString(column.VarName)
		if i != len(col.Columns)-1 {
			res.WriteString(" , ")
		}
	}
	return res.String()
}

// extract the column name of field. ex: from tag gorm:"column:user_name;type:varchar(64);primaryKey" we get "user_name"
func extractColumn(fieldName, tagName string) (column ColumnType) {
	column.FieldName = fieldName

	opts := strings.Split(tagName, ";")
	for _, opt := range opts {

		// support composite index
		if strings.HasPrefix(opt, "primaryKey") || strings.HasPrefix(opt, "unique") || strings.HasPrefix(opt, "uniqueIndex") {
			kv := strings.Split(opt, ":")
			if len(kv) == 1 {
				column.UniqueIndexName = strings.TrimSpace(fieldName)
			} else if len(kv) > 1 {
				subStrs := strings.Split(kv[1], ",")
				if len(subStrs) == 0 || subStrs[0] == "" {
					//  gorm:"uniqueIndex:,sort:desc"
					column.UniqueIndexName = strings.TrimSpace(fieldName)
				} else {
					//  gorm:"uniqueIndex:idx_name,sort:desc"
					column.UniqueIndexName = strings.TrimSpace(kv[1])
				}

			}

		} else if strings.HasPrefix(opt, "index") {
			kv := strings.Split(opt, ":")
			if len(kv) == 1 {
				column.IndexName = strings.TrimSpace(fieldName)
			} else if len(kv) > 1 {
				subStrs := strings.Split(kv[1], ",")
				if len(subStrs) == 0 || subStrs[0] == "" {
					//  gorm:"index:,sort:desc"
					column.IndexName = strings.TrimSpace(fieldName)
				} else {
					//  gorm:"index:idx_name,sort:desc"
					column.IndexName = strings.TrimSpace(kv[1])
				}

			}
		} else {
			kv := strings.Split(opt, ":")
			if len(kv) != 2 {
				continue
			}

			key := strings.TrimSpace(kv[0])
			if key == "column" {
				column.Column = strings.TrimSpace(kv[1])
			}
		}

	}
	return

}
