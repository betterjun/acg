package tpl

type Table struct {
	Name        string   `json:"name"`   // 表名
	Note        string   `json:"note"`   // 表的说明
	Fields      []Column `json:"fields"` // 所有字段
	PrimaryKeys []string // 主键名数组，从Fields中提取出来的
	Params      []string `json:"params"`  // 用于查询
	Indexes     []string `json:"indexes"` // 索引数组

	AutoIndexes map[string]*AutoIndex `json:"AutoIndexes"` // 根据字段，自动建立的索引数组
}

type AutoIndex struct {
	Name   string    // 索引名
	Fields []*Column `json:"fields"` // 所有字段
}

// 填充默认值
func (o *Table) FillDefault() {
	for _, v := range o.Fields {
		v.FillDefault()
	}

	// 构建自动索引
	o.AutoIndexes = make(map[string]*AutoIndex)
	for i, _ := range o.Fields {
		v := &o.Fields[i]
		o.addAutoIndex("index", v)
		o.addAutoIndex("unique_index", v)
		o.addPrimaryIndex("primary", v)
	}

	return
}

func (o *Table) addPrimaryIndex(keyName string, column *Column) {
	primary := column.Attrs.GetBool(keyName)
	if !primary {
		return
	}

	name := "idx_" + keyName
	if v, ok := o.AutoIndexes[name]; !ok {
		o.AutoIndexes[name] = &AutoIndex{Name: name, Fields: []*Column{column}}
	} else {
		v.Fields = append(v.Fields, column)
	}
}

func (o *Table) addAutoIndex(keyName string, column *Column) {
	name := column.Attrs.GetString(keyName)
	if len(name) == 0 {
		return
	}

	if v, ok := o.AutoIndexes[name]; !ok {
		o.AutoIndexes[name] = &AutoIndex{Name: name, Fields: []*Column{column}}
	} else {
		v.Fields = append(v.Fields, column)
	}
}

func (o *Table) GetModelPath() string {
	//return fmt.Sprintf("%v/%v", o.getModulePath(), o.Package)
	return ""
}

func (o *Table) GetModelFileName(tablename string) string {
	//return fmt.Sprintf("%v/%v.go", o.getModelPath(), tablename)
	return ""
}

func (e *Table) HasDateTimeField() bool {
	for _, v := range e.Fields {
		switch v.Type {
		case "date", "datetime", "time", "timestamp":
			return true
		}
	}
	return false
}

func (e *Table) GetPrimaryIntKeys() (keys []string) {
	keys = make([]string, 0)
	for _, v := range e.Fields {
		if !v.Attrs.Exists(GormFieldTag_Primary) {
			continue
		}

		switch v.Type {
		case "tinyint", "int", "smallint", "mediumint":
			keys = append(keys, v.Name)
		case "bigint":
			keys = append(keys, v.Name)
		case "decimal", "double":
			keys = append(keys, v.Name)
		case "float":
			keys = append(keys, v.Name)
		}
	}
	return keys
}

func (e *Table) GetPrimaryStringKeys() (keys []string) {
	keys = make([]string, 0)
	for _, v := range e.Fields {
		if !v.Attrs.Exists(GormFieldTag_Primary) {
			continue
		}

		switch v.Type {
		case "char", "enum", "varchar", "longtext", "mediumtext", "text", "tinytext":
			keys = append(keys, v.Name)
		}
	}
	return keys
}
