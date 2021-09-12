package tpl

import (
	"fmt"
	"strings"
)

// model定义
type ModelCfg struct {
	// Struct  StructCfg  // model中的字段定义
	Name       string                // 表明
	Comment    string                // 表的说明
	Columns    []*ColumnCfg          `yaml:"column"` // model中的字段定义
	ColumnDict map[string]*ColumnCfg `yaml:"-"`      // model中的字段定义
	Indexes    []*IndexCfg           `yaml:"index"`  // model中的索引定义
	IndexDict  map[string]*IndexCfg  `yaml:"-"`      // model中的索引定义
	Query      []*QueryCfg           `yaml:"query"`  // 包中的查询配置
	QueryMap   QueryMap              `yaml:"-"`      // 包中的查询配置
}

// 获取goimports导入数组
func (m *ModelCfg) GoImportsArray() []string {
	importsMap := make(map[string]bool, 0)

	for _, v := range m.Columns {
		importsMap = v.goImportsMap(importsMap)
	}

	for _, v := range m.Query {
		importsMap = v.goImportsMap(importsMap)
	}

	imports := make([]string, 0)
	for k, _ := range importsMap {
		imports = append(imports, k)
	}

	return imports
}

// 获取goimports导入字符串，用换行符分割
func (m *ModelCfg) GoImportsString() string {
	imports := m.GoImportsArray()
	return strings.Join(imports, "\n")
}

// 获取当前字段的gorm tag
func (m *ModelCfg) GetGormIndexTag(columnName string) string {
	idx := ""
	for _, v := range m.Indexes {
		if v.hasKey(columnName) {
			if v.Type == "unique" {
				idx += fmt.Sprintf("uniqueIndex:%v;", v.Name)
			} else {
				idx += fmt.Sprintf("index:%v;", v.Name)
			}
		}
	}

	return idx
}

// 获取字段的gorm tag，包含本身以及索引的
func (m *ModelCfg) GetGormTag(columnName string) string {
	column := m.ColumnDict[columnName]
	if column == nil {
		panic(fmt.Errorf("column %v not existed in model %v", columnName, m.Name))
	}
	gormTag := column.GetGormTag()
	gormTag += m.GetGormIndexTag(columnName)
	return gormTag
}

// model按表名进行组织
type ModelMap map[string]*ModelCfg

func NewModelCfg() *ModelCfg {
	return &ModelCfg{
		Columns:    make([]*ColumnCfg, 0),
		ColumnDict: make(map[string]*ColumnCfg),
		Indexes:    make([]*IndexCfg, 0),
		IndexDict:  make(map[string]*IndexCfg),
		Query:      make([]*QueryCfg, 0),
		QueryMap:   make(map[string]*QueryCfg),
	}
}
