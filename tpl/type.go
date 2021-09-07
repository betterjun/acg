package tpl

// 数据类型
type TypeCfg struct {
	Type     string // 在yaml文件中的字段名
	GoType   string `yaml:"goType,omitempty"`   // 生成go代码时的type名
	GoImport string `yaml:"goImport,omitempty"` // 生成go代码时的import路径
}

// 类型字典
type TypeDict map[string]*TypeCfg
