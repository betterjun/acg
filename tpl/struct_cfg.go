package tpl

import "strings"

// 结构体配置
type StructCfg struct {
	Name       string                // 结构体名
	Array      bool                  // 输出是否为数组
	Columns    []*ColumnCfg          `yaml:"column"` // 包含的字段
	ColumnDict map[string]*ColumnCfg `yaml:"-"`      // 包含的字段
}

func NewStructCfg() *StructCfg {
	return &StructCfg{
		Columns:    make([]*ColumnCfg, 0),
		ColumnDict: make(map[string]*ColumnCfg),
	}
}

// 获取goimports导入数组
func (c *StructCfg) GoImportsArray(imports []string) []string {
	if imports == nil {
		imports = make([]string, 0)
	}

	for k, _ := range c.goImportsMap(nil) {
		imports = append(imports, k)
	}

	return imports
}

// 获取goimports导入字符串，用换行符分割
func (c *StructCfg) GoImportsString() string {
	imports := make([]string, 0)
	for k, _ := range c.goImportsMap(nil) {
		imports = append(imports, k)
	}

	return strings.Join(imports, "\n")
}

func (c *StructCfg) goImportsMap(imports map[string]bool) map[string]bool {
	for _, v := range c.Columns {
		imports = v.goImportsMap(imports)
	}

	return imports
}
