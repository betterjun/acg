package tpl

import (
	"fmt"
)

type Column struct {
	Name  string     `json:"name"`  // 字段名字
	Note  string     `json:"note"`  // 字段说明
	Type  ColumnType `json:"type"`  // 字段类型
	Attrs Attribute  `json:"attrs"` // 字段属性
}

// 填充默认值
func (o *Column) FillDefault() {

}

func (o *Column) FormattedName() string {
	return FormatFieldName(stringifyFirstChar(o.Name))
}

func (o *Column) GoTypeName() string {
	nullable := false
	return mysqlTypeToGoType(string(o.Type), nullable)
}

func (o *Column) GormTag() string {
	gormTags := o.Attrs.GetAllGormTag()
	return fmt.Sprintf("gorm:\"column:%s%s;comment:'%s'\"", o.Name, gormTags, o.Note)
}

func (o *Column) JsonTag() string {
	return fmt.Sprintf("json:\"%s\"", o.Name)
}

//
//func (o *Column) GenerateField(col *Column, options *Options) string {
//	fieldName := FormatFieldName(stringifyFirstChar(col.Name))
//	nullable := false
//	valueType := mysqlTypeToGoType(string(col.Type), nullable)
//
//	gormTags := col.Attrs.GetAllGormTag()
//
//	var annotations []string
//	if options.GenJson {
//		annotations = append(annotations, fmt.Sprintf("json:\"%s\"", col.Name))
//	}
//	if options.GenGorm {
//		annotations = append(annotations, fmt.Sprintf("gorm:\"column:%s%s;comment:'%s'\"", col.Name, gormTags, col.Note))
//	}
//	if len(annotations) > 0 {
//		return fmt.Sprintf("\t%s %s `%s` // %s\n",
//			fieldName,
//			valueType,
//			strings.Join(annotations, " "),
//			col.Note)
//	} else {
//		return fmt.Sprintf("\t%s %s // %s\n",
//			fieldName,
//			valueType,
//			col.Note)
//	}
//}
