package tpl

// 索引配置
type IndexCfg struct {
	Name    string   // 索引名
	Keys    []string // 简单、结构体、数组
	Comment string   // 说明
	Type    string   `yaml:"type,omitempty"` // 索引类型, index或unique
}

// 查看index中是否包含特定字段
func (i *IndexCfg) hasKey(columnName string) bool {
	for _, v := range i.Keys {
		if v == columnName {
			return true
		}
	}

	return false
}
