package tpl

import (
	"fmt"
	"strings"
)

// 字段的数据字典
type ColumnDict map[string]*ColumnCfg

type ColumnCfgType int

const (
	// 0单字段、1数组、2子结构体
	ColumnCfgTypeSimple = 0
	ColumnCfgTypeArray  = 1
	ColumnCfgTypeStruct = 2
)

// 字段配置
type ColumnCfg struct {
	// "name", "type", "comment", "primary", "notnull", "unique", "auto_increment", "default"
	Name          string // 字段名
	Type          string // 在yaml文件中的字段名
	Comment       string `yaml:"comment,omitempty"`       // 字段说明，从属性表中抽出来的
	Primary       bool   `yaml:"primary,omitempty"`       // 是否是主键
	NotNull       bool   `yaml:"notnull,omitempty"`       // 是否不允许为空
	Unique        bool   `yaml:"unique,omitempty"`        // 是否必须唯一
	AutoIncrement bool   `yaml:"autoIncrement,omitempty"` // 是否自增
	Default       string `yaml:"default,omitempty"`       // 默认值
	GoType        string `yaml:"goType,omitempty"`        // 生成go代码时的type名
	GoImport      string `yaml:"goImport,omitempty"`      // 生成go代码时的import路径
	FieldType     int    `yaml:"fieldType,omitempty"`     // 0单字段、1数组、2子结构体

	//Type  ColumnType `json:"type"`
	// column: name type comment default primary notnull unique auto_increment
	Attrs Attribute `yaml:"-" json:"attrs"`

	Children map[string]*ColumnCfg `yaml:"children,omitempty"` // 如果是子结构体，则包含子结构体的字段
}

func NewColumnCfg() *ColumnCfg {
	return &ColumnCfg{
		Children: make(map[string]*ColumnCfg),
	}
}

// 获取类型，优先使用GoType，次选Type
func (c *ColumnCfg) GoTypeString() string {
	if len(c.GoType) > 0 {
		return c.GoType
	}

	return c.Type
}

// 获取goimports导入数组
func (c *ColumnCfg) GoImportsArray(imports []string) []string {
	if imports == nil {
		imports = make([]string, 0)
	}

	for k, _ := range c.goImportsMap(nil) {
		imports = append(imports, k)
	}

	return imports
}

// 获取goimports导入字符串，用换行符分割
func (c *ColumnCfg) GoImportsString() string {
	imports := make([]string, 0)
	for k, _ := range c.goImportsMap(nil) {
		imports = append(imports, k)
	}

	return strings.Join(imports, "\n")
}

// 获取goimports导入数组
func (c *ColumnCfg) goImportsMap(imports map[string]bool) map[string]bool {
	if imports == nil {
		imports = make(map[string]bool, 0)
	}
	switch c.FieldType {
	case 0:
		if len(c.GoImport) > 0 {
			if _, ok := imports[c.GoImport]; !ok {
				imports[c.GoImport] = true
			}
		}
	case 1:
		fallthrough
	case 2:
		for _, v := range c.Children {
			imports = v.goImportsMap(imports)
		}
	}

	return imports
}

// 获取当前字段的gorm tag
func (c *ColumnCfg) GetGormTag() string {
	//    gorm:"{{.Name}}
	//   {{- if .Primary -}} ;primary {{- end -}}   {{- if .NotNull -}} ;notnull {{- end -}}
	//  {{- if .Unique -}} ;unique {{- end -}} {{- if gt $le 0 -}} ;default= {{- .Default -}} {{- end -}}"

	primaryStr := ""
	if c.Primary {
		primaryStr = "primary;"
	}
	notNullStr := ""
	if c.NotNull {
		notNullStr = "notnull;"
	}
	uniqueStr := ""
	if c.Unique {
		uniqueStr = "unique;"
	}
	defaultStr := ""
	if len(c.Default) > 0 {
		defaultStr = fmt.Sprintf("default=%v;", c.Default)
	}

	return fmt.Sprintf(`%v;%v%v%v%v`,
		c.Name, primaryStr, notNullStr, uniqueStr, defaultStr)
}
