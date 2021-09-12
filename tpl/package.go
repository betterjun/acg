package tpl

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/viper"
)

// package 定义
type PackageCfg struct {
	Name       string       // 包名
	Type       []*TypeCfg   `yaml:"type"`  // 类型字典
	TypeDict   TypeDict     `yaml:"-"`     // 类型字典
	Column     []*ColumnCfg `yaml:"dict"`  // 字段字典
	ColumnDict ColumnDict   `yaml:"-"`     // 字段字典
	Model      []*ModelCfg  `yaml:"model"` // 包中的model配置
	ModelMap   ModelMap     `yaml:"-"`     // 包中的model配置

	CurrentModel *ModelCfg `yaml:"-"` // 包中的查询配置
	CurrentQuery *QueryCfg `yaml:"-"` // 包中的查询配置
}

func NewPackageCfg() *PackageCfg {
	return &PackageCfg{
		Type:       make([]*TypeCfg, 0),
		TypeDict:   make(map[string]*TypeCfg),
		Column:     make([]*ColumnCfg, 0),
		ColumnDict: make(map[string]*ColumnCfg),
		ModelMap:   make(map[string]*ModelCfg),
	}
}

// 解析package
func (pkgCfg *PackageCfg) parsePackage(config *viper.Viper) (err error) {
	/*
		解析package、type、dict、model
	*/

	pkgCfg.Name = config.GetString("package.name")

	dataCfg := config.Get("type")
	err = pkgCfg.parseType(dataCfg)
	if err != nil {
		return err
	}

	for b, i := true, 0; b; i++ {
		key := fmt.Sprintf("dict.%v", i)
		b = config.IsSet(key)
		fmt.Printf("dict.%v found: %v\n", i, b)
		if !b {
			break
		}

		dataCfg = config.Get(key)
		err = pkgCfg.parseDict(dataCfg)
		if err != nil {
			return err
		}
	}

	err = pkgCfg.parseModels(config)
	if err != nil {
		return err
	}

	return nil
}

// 解析type
func (pkgCfg *PackageCfg) parseType(dataCfg interface{}) (err error) {
	arr, ok := dataCfg.([]interface{})
	if !ok {
		return fmt.Errorf("parseType error: not array")
	}

	for _, v1 := range arr {
		switch v1.(type) {
		case map[interface{}]interface{}:
			newType := new(TypeCfg)
			for k2, v2 := range v1.(map[interface{}]interface{}) {
				switch v2.(type) {
				case string:
					k2Str := k2.(string)
					switch k2Str {
					case "type":
						newType.Type = v2.(string)
					case "goType":
						newType.GoType = v2.(string)
						// todo check the GoType is valid format
					case "goImport":
						newType.GoImport = v2.(string)
						// todo check the GoImport is valid format

					default:
						fmt.Sprintf("parseType key %q is ignored\n", k2)
					}
				default:
					return fmt.Errorf("parseType key %q is not a string\n", k2)
				}
			}

			pkgCfg.Type = append(pkgCfg.Type, newType)
			pkgCfg.TypeDict[newType.Type] = newType

		default:
			return fmt.Errorf("parseType format error: not a struct array\n")
		}
	}

	return err
}

// 解析dict
func (pkgCfg *PackageCfg) parseDict(dataCfg interface{}) (err error) {
	switch dataCfg.(type) {
	case map[interface{}]interface{}:
		column, err := parseColumnCfg(pkgCfg.TypeDict, pkgCfg.ColumnDict, dataCfg.(map[interface{}]interface{}))
		if err != nil {
			return err
		}
		pkgCfg.Column = append(pkgCfg.Column, column)
		pkgCfg.ColumnDict[column.Name] = column
	default:
		return fmt.Errorf("parseDict format error: not a struct\n")
	}

	return err
}

// 解析model数组
func (pkgCfg *PackageCfg) parseModels(config *viper.Viper) (err error) {
	for b, i := true, 0; b; i++ {
		key := fmt.Sprintf("model.%v", i)
		b = config.IsSet(key)
		fmt.Printf("key %v found: %v\n", key, b)
		if !b {
			break
		}

		modelCfg := config.Sub(key)
		model, err := pkgCfg.parseOneModel(modelCfg)
		if err != nil {
			return err
		}

		pkgCfg.Model = append(pkgCfg.Model, model)
		pkgCfg.ModelMap[model.Name] = model
	}

	return err
}

// 解析model
func (pkgCfg *PackageCfg) parseOneModel(modelCfg *viper.Viper) (model *ModelCfg, err error) {
	modelName := modelCfg.GetString("name")
	if len(modelName) == 0 {
		return nil, fmt.Errorf("parseOneModel no model name found")
	}

	model = NewModelCfg()
	model.Name = modelName
	model.Comment = modelCfg.GetString("comment")

	err = pkgCfg.parseColumns(model, modelCfg)
	if err != nil {
		return
	}
	err = pkgCfg.parseIndexes(model, modelCfg)
	if err != nil {
		return
	}

	err = pkgCfg.parseQueries(model, modelCfg)
	if err != nil {
		return
	}

	return
}

// 解析column
func (pkgCfg *PackageCfg) parseColumns(model *ModelCfg, columnCfg *viper.Viper) (err error) {
	for b, i := true, 0; b; i++ {
		key := fmt.Sprintf("column.%v", i)
		b = columnCfg.IsSet(key)
		fmt.Printf("key %v found: %v\n", key, b)
		if !b {
			break
		}

		dataCfg := columnCfg.Get(key)
		err = pkgCfg.parseOneColumn(model, dataCfg)
		if err != nil {
			return err
		}
	}

	return
}

func (pkgCfg *PackageCfg) parseOneColumn(model *ModelCfg, dataCfg interface{}) (err error) {
	switch dataCfg.(type) {
	case map[interface{}]interface{}:
		column, err := parseColumnCfg(pkgCfg.TypeDict, pkgCfg.ColumnDict, dataCfg.(map[interface{}]interface{}))
		if err != nil {
			return err
		}
		model.Columns = append(model.Columns, column)
		model.ColumnDict[column.Name] = column
	default:
		return fmt.Errorf("parseDict format error: not a struct\n")
	}

	return err
}

// 解析index
func (pkgCfg *PackageCfg) parseIndexes(model *ModelCfg, indexCfg *viper.Viper) (err error) {
	for b, i := true, 0; b; i++ {
		key := fmt.Sprintf("index.%v", i)
		b = indexCfg.IsSet(key)
		fmt.Printf("key %v found: %v\n", key, b)
		if !b {
			break
		}

		dataCfg := indexCfg.Get(key)
		err = pkgCfg.parseOneIndex(model, dataCfg)
		if err != nil {
			return err
		}
	}

	return err
}

// 解析index
func (pkgCfg *PackageCfg) parseOneIndex(model *ModelCfg, dataCfg interface{}) (err error) {
	switch dataCfg.(type) {
	case map[interface{}]interface{}:
		index, err := parseIndexCfg(model.ColumnDict, dataCfg.(map[interface{}]interface{}))
		if err != nil {
			return err
		}
		model.Indexes = append(model.Indexes, index)
		model.IndexDict[index.Name] = index
	default:
		return fmt.Errorf("parseOneIndex format error: not a struct\n")
	}

	return err
}

// 解析query
func (pkgCfg *PackageCfg) parseQueries(model *ModelCfg, queryCfg *viper.Viper) (err error) {
	// 合并字段字典
	mergedDict := make(map[string]*ColumnCfg)
	for k, v := range pkgCfg.ColumnDict {
		mergedDict[k] = v
	}

	for k, v := range model.ColumnDict {
		mergedDict[fmt.Sprintf("%v.%v", model.Name, k)] = v
	}

	for b, i := true, 0; b; i++ {
		key := fmt.Sprintf("query.%v", i)
		b = queryCfg.IsSet(key)
		fmt.Printf("key %v found: %v\n", key, b)
		if !b {
			break
		}

		err = pkgCfg.parseOneQuery(model, mergedDict, queryCfg.Sub(key))
		if err != nil {
			return err
		}
	}

	return err
}

// 解析query
func (pkgCfg *PackageCfg) parseOneQuery(model *ModelCfg, columnDict ColumnDict, queryCfg *viper.Viper) (err error) {
	name := queryCfg.GetString("name")
	if len(name) == 0 {
		return fmt.Errorf("parseOneQuery no query name found")
	}

	query := new(QueryCfg)
	query.Name = name
	query.Comment = queryCfg.GetString("comment")
	query.SQLstr = queryCfg.GetString("sql")
	query.Pager = queryCfg.GetBool("pager")

	dataCfg := queryCfg.Get("inputs")
	if dataCfg == nil {
		return fmt.Errorf("query(%v) has no inputs", query.Name)
	}
	query.In = NewStructCfg()
	query.In.Name = query.Name + "IN"
	//dataCfg = map[interface{}]interface{}{query.In.Name: dataCfg}
	err = parseQueryStruct(query.In, pkgCfg.TypeDict, columnDict, dataCfg)
	if err != nil {
		return err
	}

	dataCfg = queryCfg.Get("outputs")
	if dataCfg == nil {
		return fmt.Errorf("query(%v) has no outputs", query.Name)
	}
	query.Out = NewStructCfg()
	query.Out.Name = query.Name + "OUT"
	//dataCfg = map[interface{}]interface{}{query.Out.Name: dataCfg}

	err = parseQueryStruct(query.Out, pkgCfg.TypeDict, columnDict, dataCfg)
	if err != nil {
		return err
	}

	enablePagerFields(query)

	model.Query = append(model.Query, query)
	model.QueryMap[query.Name] = query

	return err
}

func enablePagerFields(query *QueryCfg) (err error) {
	/*
		如果没启用pager，则直接返回
		如果不是查询语句，则直接返回或报错
		查看输入是否有limit、offset字段，没有则插入

		如果输出是array，
			构造新的输出，将原来的输出数组作为result字段，并增加count、limit、offset字段
		如果输出是struct，
			查看是否有count、limit、offset、result字段
				如果都有，则不做处理
				如果有部分，则报错
				如果全都没有，则构造新的输出，把输出当作结构体，作为result字段，并增加count、limit、offset字段
	*/

	if !query.Pager {
		return
	}
	sql := convertSpaces(strings.ToLower(query.SQLstr))
	if strings.Index(sql, "select ") < 0 { // 不是查询语句
		//return fmt.Errorf("%v output is not a select statement, ignore pager flag", query.Name)
		query.Pager = false
		fmt.Printf("%v output is not a select statement, ignore pager flag\n", query.Name)
		return
	}

	// 处理输入
	if _, ok := query.In.ColumnDict["limit"]; !ok {
		column := &ColumnCfg{
			Name:    "limit",
			Type:    "int",
			Comment: "分页大小，自动添加",
		}

		query.In.Columns = append(query.In.Columns, column)
		query.In.ColumnDict["limit"] = column
	}
	if _, ok := query.In.ColumnDict["offset"]; !ok {
		column := &ColumnCfg{
			Name:    "offset",
			Type:    "int",
			Comment: "偏移量，自动添加",
		}

		query.In.Columns = append(query.In.Columns, column)
		query.In.ColumnDict["offset"] = column
	}

	// 处理输出
	if query.Out.Array {
		return replaceOutStruct(query)
	} else {
		_, ok1 := query.Out.ColumnDict["limit"]
		_, ok2 := query.Out.ColumnDict["offset"]
		_, ok3 := query.Out.ColumnDict["count"]
		_, ok4 := query.Out.ColumnDict["result"]

		if ok1 && ok2 && ok3 && ok4 { // 全都有，则不做处理
			return
		}

		//if !(ok1 && ok2 && ok3 && ok4) && (ok1 || ok2 || ok3 || ok4) { // 仅有部分
		if ok1 || ok2 || ok3 || ok4 { // 仅有部分
			return fmt.Errorf("%v output struct has some pager fields", query.Name)
		}

		// 全都没有
		return replaceOutStruct(query)
	}
	return
}

func replaceOutStruct(query *QueryCfg) (err error) {
	// 结构体配置
	newOut := NewStructCfg()
	newOut.Name = query.Out.Name
	newOut.Array = false

	column := &ColumnCfg{
		Name:    "limit",
		Type:    "int",
		Comment: "分页大小，自动添加",
	}
	newOut.Columns = append(newOut.Columns, column)
	newOut.ColumnDict["limit"] = column

	column = &ColumnCfg{
		Name:    "offset",
		Type:    "int",
		Comment: "偏移量，自动添加",
	}
	newOut.Columns = append(newOut.Columns, column)
	newOut.ColumnDict["offset"] = column

	column = &ColumnCfg{
		Name:    "count",
		Type:    "int64",
		Comment: "总数，自动添加",
	}
	newOut.Columns = append(newOut.Columns, column)
	newOut.ColumnDict["count"] = column

	columnArr := &ColumnCfg{
		Name:      "result",
		Comment:   "输出数据，自动添加",
		FieldType: 1,
		Children:  make(map[string]*ColumnCfg),
	}
	for _, v := range query.Out.Columns {
		columnArr.Children[v.Name] = v
	}
	newOut.Columns = append(newOut.Columns, columnArr)
	newOut.ColumnDict["result"] = columnArr

	query.Out = newOut
	return
}

func parseQueryStruct(structCfg *StructCfg, typeDict TypeDict, columnDict ColumnDict, dataCfg interface{}) (err error) {
	/*
	   inputs: {name: tb.name, page: [int, 页码], size: [int, 分页大小]}
	   outputs: [{age: ta.age, count: [int, 总数], objName: {name: tb.name}, data: [{id: ta.id, sex: [int, 性别]}]}]
	*/
	/*
		struct =  obj
		obj = "{" pair {pair} "}"
		pair = STRING ":" expression
		expression = STRING | array | obj
		array = "[" obj | STRING,STRING|STRING,STRING,STRING "]"
	*/

	// TODO 实现接口
	switch dataCfg.(type) {
	case map[interface{}]interface{}: // {name: tb.name, page: [int, 页码], size: [int, 分页大小]}
		rootColumn := new(ColumnCfg)
		err = parseStructFieldCfg(typeDict, columnDict, rootColumn, dataCfg.(map[interface{}]interface{}))
		if err != nil {
			return err
		}
		if rootColumn.Children != nil {
			structCfg.ColumnDict = rootColumn.Children
			for _, v := range rootColumn.Children {
				structCfg.Columns = append(structCfg.Columns, v)
			}
		}
	case []interface{}: // [{age: ta.age, count: [int, 总数], objName: {name: tb.name}, data: [{id: ta.id, sex: [int, 性别]}]}]
		rootColumn := new(ColumnCfg)
		err = parseStructFieldCfgArr(typeDict, columnDict, rootColumn, dataCfg.([]interface{}))
		if err != nil {
			return err
		}
		if rootColumn.Children != nil {
			structCfg.ColumnDict = rootColumn.Children
			for _, v := range rootColumn.Children {
				structCfg.Columns = append(structCfg.Columns, v)
			}
		}
		structCfg.Array = true

	default:
		return fmt.Errorf("parseQueryStruct format error: not a struct\n")
	}

	return
}

// 解析字段
func parseStructFieldCfg(typeDict TypeDict, columnDict ColumnDict, rootColumn *ColumnCfg, dataCfg map[interface{}]interface{}) (err error) {
	/*
		字典定义，支持下面三种写法
				name: dict.name
		        age: {type: int, comment: 年龄} # 对象定义方式
		        sex: [int, 性别, 1, false, false, false, false] # 数组定义方式，属性有顺序，最少1个属性，最多支持7个属性，多的被忽略
				data: [{id: ta.id, sex: [int, 性别]}]# 对象数组定义方式，数组中只有1个结构体

	*/

	// k为name，v为其他属性的map
	for k, vAttr := range dataCfg {
		subColumn := new(ColumnCfg)
		subColumn.Name = fmt.Sprintf("%v", k)
		switch vAttr.(type) {
		case string: // 字符串，参考字典定义的其他字段，用于重命名
			err = parseStructFieldCfgStr(columnDict, subColumn, vAttr.(string))
		case []interface{}: // 字符串数组，分别对应目前的属性
			err = parseStructFieldCfgArr(typeDict, columnDict, subColumn, vAttr.([]interface{}))
		case map[interface{}]interface{}: // 对象定义，key/value形式定义
			err = parseStructFieldCfg(typeDict, columnDict, subColumn, vAttr.(map[interface{}]interface{}))
		default:
			err = fmt.Errorf("parseFieldCfg wrong format for field")
		}

		if err != nil {
			return err
		}

		if rootColumn.Children == nil {
			rootColumn.Children = make(map[string]*ColumnCfg)
			rootColumn.FieldType = ColumnCfgTypeStruct
		}
		rootColumn.Children[subColumn.Name] = subColumn
	}

	return
}

func parseStructFieldCfgStr(columnDict ColumnDict, column *ColumnCfg, dataCfg string) (err error) {
	columnName, err := parseString(dataCfg)
	if err != nil {
		return err
	}

	// 在ColumnDict中查找，属性用以前的来定义
	c, ok := columnDict[columnName]
	if !ok {
		return fmt.Errorf("column name %v not found in dict", columnName)
	}

	newColumnName := column.Name
	*column = *c
	column.Name = newColumnName
	return
}

func parseStructFieldCfgArr(typeDict TypeDict, columnDict ColumnDict, column *ColumnCfg, columnCfgs []interface{}) (err error) {
	length := len(columnCfgs)
	switch length {
	case 1: // data: [{id: ta.id, sex: [int, 性别]}]
		subColumnMapCfg, ok := columnCfgs[0].(map[interface{}]interface{})
		if !ok {
			return fmt.Errorf("not a struct in array")
		}

		err = parseStructFieldCfg(typeDict, columnDict, column, subColumnMapCfg)
		if err != nil {
			return err
		}
		column.FieldType = ColumnCfgTypeArray
	case 2: // count: [int, 总数]
		for k, v := range columnCfgs {
			if k == 0 {
				column.Type, _ = parseString(v)
			} else if k == 1 {
				column.Comment, _ = parseString(v)
			}
		}

	case 3: // arrType: [array, int, 总数]
		for k, v := range columnCfgs {
			if k == 0 {
				column.Type, _ = parseString(v)
				// todo 检查类型是否正确
				if strings.Compare("array", strings.ToLower(strings.TrimSpace(column.Type))) == 0 {
					column.FieldType = ColumnCfgTypeArray
				}
			} else if k == 1 {
				column.Type, _ = parseString(v)
			} else if k == 2 {
				column.Comment, _ = parseString(v)
			}
		}
	default:
		return fmt.Errorf("parseStructFieldCfgArr format error: not a valid array\n")
	}

	return
}

//
//// 解析结构体
//func parseQueryStructCfg(structCfg *StructCfg, typeDict TypeDict, columnDict ColumnDict, dataCfg interface{}) (err error) {
//
//	return
//}
//
//// 解析字段
//func (pkgCfg *PackageCfg) parseField(dataCfg interface{}) (err error) {
//
//	return err
//}

// 解析根字段
func parseFieldCfgColumn(typeDict TypeDict, columnDict ColumnDict, dataCfg map[interface{}]interface{}) (column *ColumnCfg, err error) {
	/*
			struct = "{" STRING ":" expression "}"
			expression = STRING | array | obj
			//array = "[" STRING "]" | "[" obj "]"
			array = "[" STRING  {STRING} "]" | "[" obj "]"
			obj = "{" pair {pair} "}"
			pair = STRING ":" expression

		column = "{" STRING ":" expression "}"
		expression = STRING | array | obj
		array = "[" STRING "]" | "[" obj "]"
		obj = "{" pair {pair} "}"
		pair = STRING ":" expression
	*/

	if len(dataCfg) != 1 {
		return nil, fmt.Errorf("column config is not valid")
	}

	// k为name，v为其他属性的map
	for k, vAttr := range dataCfg {
		// column为根节点，用来生成结构体名
		column = new(ColumnCfg)
		column.Name = fmt.Sprintf("%v", k)
		//column.FieldType = 2

		switch vAttr.(type) {
		case string: // 字符串，参考字典定义的其他字段，用于重命名
			err = parseFieldCfgStr(columnDict, column, vAttr.(string))
		case []interface{}: // 字符串数组，分别对应目前的属性
			err = parseFieldCfgArr(typeDict, columnDict, column, vAttr.([]interface{}))
		case map[interface{}]interface{}: // 对象定义，key/value形式定义
			err = parseFieldCfgMap(typeDict, columnDict, column, vAttr.(map[interface{}]interface{}))
		default:
			err = fmt.Errorf("parseFieldCfg wrong format for field")
		}

		if err != nil {
			column.Children = make(map[string]*ColumnCfg)
			//column.Children[subColumn.Name] = subColumn
		}

		break // 只处理一个元素
	}

	return
}

// 解析字段
func parseFieldCfg(typeDict TypeDict, columnDict ColumnDict, dataCfg map[interface{}]interface{}) (column *ColumnCfg, err error) {
	/*
		字典定义，支持下面三种写法
				name: dict.name
		        age: {type: int, comment: 年龄} # 对象定义方式
		        sex: [int, 性别, 1, false, false, false, false] # 数组定义方式，属性有顺序，最少1个属性，最多支持7个属性，多的被忽略
				data: [{id: ta.id, sex: [int, 性别]}]# 对象数组定义方式，数组中只有1个结构体

	*/

	if len(dataCfg) != 1 {
		return nil, fmt.Errorf("dataCfg is not valid")
	}

	// k为name，v为其他属性的map
	for k, vAttr := range dataCfg {
		column = new(ColumnCfg)
		column.Name = fmt.Sprintf("%v", k)
		switch vAttr.(type) {
		case string: // 字符串，参考字典定义的其他字段，用于重命名
			err = parseFieldCfgStr(columnDict, column, vAttr.(string))
		case []interface{}: // 字符串数组，分别对应目前的属性
			err = parseFieldCfgArr(typeDict, columnDict, column, vAttr.([]interface{}))
		case map[interface{}]interface{}: // 对象定义，key/value形式定义
			err = parseFieldCfgMap(typeDict, columnDict, column, vAttr.(map[interface{}]interface{}))
		default:
			err = fmt.Errorf("parseFieldCfg wrong format for field")
		}

		break // 只处理一个元素
	}

	return
}

// 解析字段
func parseColumnCfg(typeDict TypeDict, columnDict ColumnDict, dataCfg map[interface{}]interface{}) (column *ColumnCfg, err error) {
	/*
		字典定义，支持下面三种写法
				name: dict.name
		        age: {type: int, comment: 年龄} # 对象定义方式
		        sex: [int, 性别, 1, false, false, false, false] # 数组定义方式，属性有顺序，最少1个属性，最多支持7个属性，多的被忽略
				data: [{id: ta.id, sex: [int, 性别]}]# 对象数组定义方式，数组中只有1个结构体

		column = "{" STRING ":" expression "}"
		expression = STRING | array | obj
		array = "[" STRING {STRING} "]"
		obj = "{" pair {pair} "}"
		pair = STRING ":" STRING
	*/

	if len(dataCfg) != 1 {
		return nil, fmt.Errorf("column config is not valid")
	}

	// k为name，v为其他属性的map
	for k, vAttr := range dataCfg {
		column = new(ColumnCfg)
		column.Name = fmt.Sprintf("%v", k)
		switch vAttr.(type) {
		case string: // 字符串，参考字典定义的其他字段，用于重命名
			err = parseFieldCfgStr(columnDict, column, vAttr.(string))
		case []interface{}: // 字符串数组，分别对应目前的属性
			err = parseColumnFieldCfgArr(typeDict, column, vAttr.([]interface{}))
		case map[interface{}]interface{}: // 对象定义，key/value形式定义
			err = parseColumnFieldCfgMap(typeDict, column, vAttr.(map[interface{}]interface{}))
		default:
			err = fmt.Errorf("parseFieldCfg wrong format for field")
		}

		break // 只处理一个元素
	}

	return
}

func parseFieldCfgStr(columnDict ColumnDict, column *ColumnCfg, dataCfg string) (err error) {
	columnName, err := parseString(dataCfg)
	if err != nil {
		return err
	}

	// 在ColumnDict中查找，属性用以前的来定义
	c, ok := columnDict[columnName]
	if !ok {
		return fmt.Errorf("column name %v not found in dict", columnName)
	}

	newColumnName := column.Name
	*column = *c
	column.Name = newColumnName
	return
}

func parseFieldCfgArr(typeDict TypeDict, columnDict ColumnDict, column *ColumnCfg, dataCfg []interface{}) (err error) {
	if len(dataCfg) == 1 {
		for k, v := range dataCfg { // 数组，且只有一个元素，第一个元素为结构体
			switch v.(type) {
			case map[interface{}]interface{}: // 结构体数组
				subDataCfg := v.(map[interface{}]interface{})
				//subColumn, err := parseFieldCfg(typeDict, columnDict, subDataCfg)
				subColumn := new(ColumnCfg)
				subColumn.Name = fmt.Sprintf("%v", k)
				err = parseFieldCfgMap(typeDict, columnDict, subColumn, subDataCfg)
				if err != nil {
					return err
				}

				column.FieldType = ColumnCfgTypeArray
				column.Children = subColumn.Children

				return nil

			case string: // todo 单类型数组，目前和属性定义冲突了
				//return nil
				// 不做处理，漏掉下面去处理

			default:
				return fmt.Errorf("parseFieldCfgArr wrong field format %v=%v\n", k, v)
			}

			break
		}
	}

	// columnName: type comment default primary notnull unique auto_increment
	for k, v := range dataCfg {
		switch k {
		case 0:
			column.Type, err = parseString(v)
			// todo 检查类型是否正确

			// 在type字典中查找
			if t, ok := typeDict[column.Type]; ok {
				column.GoType = t.GoType
				column.GoImport = t.GoImport
			}
		case 1:
			column.Comment, err = parseString(v)
		case 2:
			column.Default, err = parseString(v)
		case 3:
			column.Primary, err = parseBool(v)
		case 4:
			column.NotNull, err = parseBool(v)
		case 5:
			column.Unique, err = parseBool(v)
		case 6:
			column.AutoIncrement, err = parseBool(v)
		default:
			fmt.Printf("parseFieldCfgArr ignored field %v=%v\n", k, v)
		}

		if err != nil {
			return err
		}
	}
	column.FieldType = ColumnCfgTypeSimple

	return
}

func parseFieldCfgMap(typeDict TypeDict, columnDict ColumnDict, column *ColumnCfg, dataCfg map[interface{}]interface{}) (err error) {
	// 解析其他属性
	res, err := parseFieldAttr(dataCfg)
	if err != nil {
		return err
	}

	if len(res) == 1 { // 子结构体
		for k, v := range res {
			switch v.(type) {
			case map[interface{}]interface{}:
				subDataCfg := v.(map[interface{}]interface{})
				subColumn, err := parseFieldCfg(typeDict, columnDict, subDataCfg)
				if err != nil {
					return err
				}

				column.FieldType = ColumnCfgTypeStruct
				column.Children = subColumn.Children

				return nil
			default:
				return fmt.Errorf("parseFieldCfgMap wrong field format %v=%v\n", k, v)
			}
			break
		}
	}

	for k, v := range res {
		switch k {
		case "type":
			column.Type, err = parseString(v)
			// todo 检查类型是否正确

			// 在type字典中查找
			if t, ok := typeDict[column.Type]; ok {
				column.GoType = t.GoType
				column.GoImport = t.GoImport
			}
		case "comment":
			column.Comment, err = parseString(v)
		case "primary":
			column.Primary, err = parseBool(v)
		case "notnull":
			column.NotNull, err = parseBool(v)
		case "unique":
			column.Unique, err = parseBool(v)
		case "autoIncrement":
			column.AutoIncrement, err = parseBool(v)
		case "default":
			column.Default, err = parseString(v)
		default:
			fmt.Printf("parseFieldCfgMap ignored field %v=%v\n", k, v)
		}

		if err != nil {
			return err
		}
	}

	column.FieldType = ColumnCfgTypeStruct

	return
}

// 解析字段
func parseFieldAttr(dataCfg interface{}) (res map[string]interface{}, err error) {
	res = make(map[string]interface{})
	for k2, v2 := range dataCfg.(map[interface{}]interface{}) {
		res[k2.(string)] = v2
	}
	return res, err
}

func parseString(data interface{}) (string, error) {
	switch data.(type) {
	case string:
		return data.(string), nil
	case nil:
		return "", nil
	default:
		return fmt.Sprintf("%v", data), nil
	}

	return "", fmt.Errorf("not a string")
}

func parseBool(data interface{}) (bool, error) {
	switch data.(type) {
	case string:
		return strconv.ParseBool(data.(string))
	case bool:
		return data.(bool), nil
	}

	return false, fmt.Errorf("not a bool")
}

func parseInt64(data interface{}) (int64, error) {
	switch data.(type) {
	case string:
		return strconv.ParseInt(data.(string), 10, 64)
	case int64:
		return data.(int64), nil
	}

	return 0, fmt.Errorf("not a int64")
}

func parseFloat64(data interface{}) (float64, error) {
	switch data.(type) {
	case string:
		return strconv.ParseFloat(data.(string), 64)
	case float64:
		return data.(float64), nil
	}

	return 0, fmt.Errorf("not a float")
}

// 解析索引
func parseIndexCfg(columnDict ColumnDict, dataCfg map[interface{}]interface{}) (index *IndexCfg, err error) {
	/*
		索引定义，支持下面两种写法
			- idx_name: [[name], 姓名索引, index] # 数组定义方式，属性有顺序，最少1个属性，最多支持3个属性，多的被忽略
			- idx_age: {keys: [age], comment: 年龄索引, type: index} # 对象定义方式
	*/

	if len(dataCfg) != 1 {
		return nil, fmt.Errorf("dataCfg is not valid")
	}

	// k为name，v为其他属性的map
	for k, vAttr := range dataCfg {
		index = new(IndexCfg)
		index.Name = fmt.Sprintf("%v", k)
		switch vAttr.(type) {
		case []interface{}: // 字符串数组，分别对应目前的属性
			err = parseIndexCfgArr(columnDict, index, vAttr.([]interface{}))
		case map[interface{}]interface{}: // 对象定义，key/value形式定义
			err = parseIndexCfgMap(columnDict, index, vAttr.(map[interface{}]interface{}))
		default:
			err = fmt.Errorf("parseIndexCfg wrong format for field")
		}

		break // 只处理一个元素
	}

	return
}

func parseIndexCfgArr(columnDict ColumnDict, index *IndexCfg, dataCfg []interface{}) (err error) {
	// indexName: [field name list], comment, type
	// type can be [primary/index/unique]

	for k, v := range dataCfg {
		switch k {
		case 0:
			err = parseIndexCfgKeys(columnDict, index, v)
		case 1:
			index.Comment, err = parseString(v)
		case 2:
			index.Type, err = parseIndexCfgType(v)
		default:
			fmt.Printf("parseIndexCfgArr ignored field %v=%v\n", k, v)
		}

		if err != nil {
			return err
		}
	}

	return
}

func parseIndexCfgMap(columnDict ColumnDict, index *IndexCfg, dataCfg map[interface{}]interface{}) (err error) {
	// 解析其他属性
	res, err := parseFieldAttr(dataCfg)
	if err != nil {
		return err
	}

	// idx_age: {keys: [age], comment: 年龄索引, type: index}
	for k, v := range res {
		switch k {
		case "keys":
			err = parseIndexCfgKeys(columnDict, index, v)
		case "comment":
			index.Comment, err = parseString(v)
		case "type":
			index.Type, err = parseIndexCfgType(v)
		default:
			fmt.Printf("parseIndexCfgMap ignored field %v=%v\n", k, v)
		}

		if err != nil {
			return err
		}
	}

	return
}

// 解析index里面的keys
func parseIndexCfgKeys(columnDict ColumnDict, index *IndexCfg, data interface{}) error {
	// 解析数组
	switch data.(type) {
	case []interface{}:
		index.Keys = make([]string, 0)
		for _, v2 := range data.([]interface{}) {
			key := fmt.Sprint(v2)
			if _, ok := columnDict[key]; !ok {
				panic(fmt.Errorf("index(%v)'s key=%v not found in model defines", index.Name, key))
			}
			index.Keys = append(index.Keys, key)
		}
	default:
		return fmt.Errorf("parseIndexCfgKeys not supported index(%v) key list format", index.Name)
	}

	return nil
}

// 解析index里面的type
func parseIndexCfgType(v interface{}) (t string, err error) {
	t, err = parseString(v)
	if err != nil {
		return "", err
	}

	idxType := strings.ToLower(t)
	switch idxType {
	case "primary":
		fallthrough
	case "index":
		fallthrough
	case "unique":
	default:
		err = fmt.Errorf("unsupported index type %v", t)
	}

	return t, err
}

func parseColumnFieldCfgArr(typeDict TypeDict, column *ColumnCfg, dataCfg []interface{}) (err error) {
	// columnName: type comment default primary notnull unique auto_increment
	for k, v := range dataCfg {
		switch k {
		case 0:
			column.Type, err = parseString(v)
			// todo 检查类型是否正确

			// 在type字典中查找
			if t, ok := typeDict[column.Type]; ok {
				column.GoType = t.GoType
				column.GoImport = t.GoImport
			}
		case 1:
			column.Comment, err = parseString(v)
		case 2:
			column.Default, err = parseString(v)
		case 3:
			column.Primary, err = parseBool(v)
		case 4:
			column.NotNull, err = parseBool(v)
		case 5:
			column.Unique, err = parseBool(v)
		case 6:
			column.AutoIncrement, err = parseBool(v)
		default:
			fmt.Printf("parseColumnFieldCfgArr ignored field %v=%v\n", k, v)
		}

		if err != nil {
			return err
		}
	}
	column.FieldType = ColumnCfgTypeSimple

	return
}

func parseColumnFieldCfgMap(typeDict TypeDict, column *ColumnCfg, dataCfg map[interface{}]interface{}) (err error) {
	// 解析其他属性
	res, err := parseFieldAttr(dataCfg)
	if err != nil {
		return err
	}

	for k, v := range res {
		switch k {
		case "type":
			column.Type, err = parseString(v)
			// todo 检查类型是否正确

			// 在type字典中查找
			if t, ok := typeDict[column.Type]; ok {
				column.GoType = t.GoType
				column.GoImport = t.GoImport
			}
		case "comment":
			column.Comment, err = parseString(v)
		case "primary":
			column.Primary, err = parseBool(v)
		case "notnull":
			column.NotNull, err = parseBool(v)
		case "unique":
			column.Unique, err = parseBool(v)
		case "autoIncrement":
			column.AutoIncrement, err = parseBool(v)
		case "default":
			column.Default, err = parseString(v)
		default:
			fmt.Printf("parseColumnFieldCfgMap ignored field %v=%v\n", k, v)
		}

		if err != nil {
			return err
		}
	}
	column.FieldType = ColumnCfgTypeSimple

	return
}

//
//type Package struct {
//	Template     string   `json:"template"` // 模板路径，可选，不配置则使用默认模板
//	Module       string   `json:"module"`   // 模块名
//	Name         string   `json:"name"`     // 包名
//	Model        string   `json:"model"`    // model的包名，默认为model
//	OutDir       string   `json:"out_dir"`  // 代码输出目录，默认当前目录
//	Tables       []Table  `json:"tables"`   // 表的数据
//	TableNames   []string // 表名数组
//	CurrentTable *Table   // 当前表格
//}
//
//// 填充默认值
//func (o *Package) FillDefault() {
//	if len(o.Model) == 0 {
//		o.Model = "model"
//	}
//
//	o.TableNames = make([]string, 0, len(o.Tables))
//	for i, _ := range o.Tables {
//		v := &o.Tables[i]
//		v.FillDefault()
//		o.TableNames = append(o.TableNames, v.Name)
//	}
//
//	return
//}
//
//func (o *Package) GetModulePath() string {
//	if o.Module == "" {
//		return "."
//	}
//	return o.Module
//}
//
//func (o *Package) GetModuleFileName() string {
//	return fmt.Sprintf("%v/go.mod", o.GetModulePath())
//}
//
//func (o *Package) GetPackagePath() string {
//	return fmt.Sprintf("%v/%v", o.GetModulePath(), o.Name)
//}
//
//func (o *Package) GetSubPackagePath(subPkgName string) string {
//	return fmt.Sprintf("%v/%v", o.GetPackagePath(), subPkgName)
//}
//
//func (o *Package) GetSubPackageFilePath(subPkgName, fileName string) string {
//	return fmt.Sprintf("%v/%v/%v.go", o.GetPackagePath(), subPkgName, fileName)
//}
