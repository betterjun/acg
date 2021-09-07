package tpl

// model定义
type ModelCfg struct {
	// Struct  StructCfg  // model中的字段定义
	Name       string                // 表明
	Comment    string                // 表的说明
	Columns    []*ColumnCfg          `yaml:"column"` // model中的字段定义
	ColumnDict map[string]*ColumnCfg `yaml:"-"`      // model中的字段定义
	Indexes    []*IndexCfg           `yaml:"index"`  // model中的索引定义
	IndexDict  map[string]*IndexCfg  `yaml:"-"`      // model中的索引定义
}

// model按表名进行组织
type ModelMap map[string]*ModelCfg

func NewModelCfg() *ModelCfg {
	return &ModelCfg{
		Columns:    make([]*ColumnCfg, 0),
		ColumnDict: make(map[string]*ColumnCfg),
		Indexes:    make([]*IndexCfg, 0),
		IndexDict:  make(map[string]*IndexCfg),
	}
}
