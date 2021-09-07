package tpl

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
