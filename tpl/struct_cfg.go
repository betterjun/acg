package tpl

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
