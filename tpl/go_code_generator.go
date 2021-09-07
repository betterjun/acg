package tpl

//
//import (
//	"bytes"
//	"fmt"
//	"io/ioutil"
//	"os"
//	"text/template"
//)
//
//// 代码生成器接口
//type CodeGenerator interface {
//}
//
//// go代码生成器
//type GoCodeGenerator struct {
//}
//
//func (g *GoCodeGenerator) G2(p *Package, options *Options) (err error) {
//
//	os.MkdirAll(p.GetSubPackagePath(p.Model), 777)
//
//	fieldName := &Field{
//		FieldType: 1,
//		Name:      "name",
//		Type:      "string",
//		Comment:   "姓名",
//		IsArray:   false,
//	}
//
//	fieldAge := &Field{
//		FieldType: 1,
//		Name:      "age",
//		Type:      "int",
//		Comment:   "年龄",
//		IsArray:   false,
//	}
//
//	fieldIntArray := &Field{
//		FieldType: 2,
//		Name:      "intArray",
//		Type:      "[]int",
//		Comment:   "整数数组",
//		IsArray:   true,
//	}
//
//	fieldStruct := &Field{
//		FieldType: 3,
//		Name:      "children",
//		Type:      "struct",
//		Comment:   "子结构体，二层嵌套",
//		IsArray:   false,
//	}
//
//	fieldStruct.Children = make(map[string]*Field)
//	fieldStruct.Children["class"] = &Field{
//		FieldType: 1,
//		Name:      "classId",
//		Type:      "int64",
//		Comment:   "班级id",
//		IsArray:   false,
//	}
//	fieldStruct.Children["classMateId"] = &Field{
//		FieldType: 2,
//		Name:      "classMateIds",
//		Type:      "[]int64",
//		Comment:   "班级成员id列表",
//		IsArray:   true,
//	}
//
//	thirdStruct := &Field{
//		FieldType: 3,
//		Name:      "embededStruct",
//		Type:      "struct",
//		Comment:   "三层嵌套测试",
//		IsArray:   false,
//		Children:  make(map[string]*Field),
//	}
//	thirdStruct.Children["strArr"] = &Field{
//		FieldType: 2,
//		Name:      "stringArray",
//		Type:      "[]string",
//		Comment:   "第四层嵌套的数组",
//		IsArray:   true,
//	}
//	fieldStruct.Children["embededStruct"] = thirdStruct
//
//	out := &StructParam{}
//	out.Fields = make(map[string]*Field)
//	out.Fields["name"] = fieldName
//	out.Fields["age"] = fieldAge
//	out.Fields["intArr"] = fieldIntArray
//	out.Fields["subMap"] = fieldStruct
//
//	res, err := generateModel2(p, out)
//	if err != nil {
//		fmt.Printf("生成失败，错误:%v\n", err)
//		return err
//	}
//	err = ioutil.WriteFile(p.GetSubPackageFilePath(p.Model, "query"), res, 666)
//	if err != nil {
//		fmt.Printf("生成失败，错误:%v\n", err)
//		return err
//	}
//
//	return err
//
//	return nil
//}
//
//func (g *GoCodeGenerator) Generate(p *Package, options *Options) (err error) {
//	os.MkdirAll(p.GetSubPackagePath(p.Model), 777)
//
//	/*
//		所有表格循环
//		所有模板循环
//
//	*/
//	for _, t := range p.Tables {
//		p.CurrentTable = &t
//		res, err := generateModel(p)
//		if err != nil {
//			fmt.Printf("表格%v生成失败，错误:%v\n", t.Name, err)
//			return err
//		}
//		err = ioutil.WriteFile(p.GetSubPackageFilePath(p.Model, t.Name), res, 666)
//		if err != nil {
//			fmt.Printf("表格%v生成失败，错误:%v\n", t.Name, err)
//			return err
//		}
//	}
//
//	p.CurrentTable = nil
//
//	return err
//}
//
//func generateModel(p *Package) ([]byte, error) {
//	templ, err := ioutil.ReadFile("./tmpl/crud/model.tmpl")
//	if err != nil {
//		return nil, err
//	}
//	t := template.Must(template.New("escape").Funcs(template.FuncMap{
//		"formatName": formatName,
//	}).Parse(string(templ)))
//	nb := bytes.NewBuffer(nil)
//	if err := t.Execute(nb, p); err != nil {
//		return nil, err
//	}
//
//	t = t
//
//	return nb.Bytes(), nil
//}
//
//func generateModel2(p *Package, o *StructParam) ([]byte, error) {
//	templ, err := ioutil.ReadFile("./tmpl/crud/query.tmpl")
//	if err != nil {
//		return nil, err
//	}
//	t := template.Must(template.New("escape").Funcs(template.FuncMap{
//		"formatName": formatName,
//	}).Parse(string(templ)))
//	nb := bytes.NewBuffer(nil)
//	if err := t.Execute(nb, o); err != nil {
//		return nil, err
//	}
//
//	t = t
//
//	return nb.Bytes(), nil
//}
//
////
////// 自定义函数使用
////func (t *Template) Funcs(funcMap FuncMap) *Template {
////
////}
//
//func formatName(name string) string {
//	return FormatFieldName(stringifyFirstChar(name))
//}
