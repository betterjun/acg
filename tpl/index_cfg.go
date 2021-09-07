package tpl

// 索引配置
type IndexCfg struct {
	Name    string   // 索引名
	Keys    []string // 简单、结构体、数组
	Comment string   // 说明
	Type    string   `yaml:"type,omitempty"` // 索引类型, index或unique
}
