// Copyright © 2021 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/betterjun/acg/tpl"

	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("generate called")

		fm := NewFileMap()
		scanDirs(templateDir, fm)
		fmt.Println(*fm)

		if len(args) == 0 {
			cmd.Usage()
			return
		}

		//outDir, err := cmd.Flags().GetString("output")
		//if err != nil {
		//	cmd.Usage()
		//	return
		//}

		for _, v := range args {
			generateCode(v, fm, outputDir)
		}

	},
}

var templateDir string
var outputDir string

func init() {
	rootCmd.AddCommand(generateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	generateCmd.Flags().StringVarP(&outputDir, "output", "o", "out", "output directory for code generation, default is out")
	generateCmd.Flags().StringVarP(&templateDir, "template", "t", "template", "template directory for code generation, default is template")

}

/*
package.tmpl
model.tmpl
query.tmpl

other.tmpl

t.txt
go.mod

文件中带有

@model.go.tmpl
@model_logic.go.tmpl
@model_handler.go.tmpl

扫描文件
所有package.tmpl文件，用全局数据过一遍
所有model.tmpl文件，用model数组过一遍
所有query.tmpl文件，用query数组过一遍

所有普通tmpl文件，用全局数据过一遍
非模板文件，直接拷贝过去

用模板生成的文件保持模板文件的目录结构不变，且都去掉.tmpl

*/
type FileMap struct {
	PackageTemplate map[string]string
	ModelTemplate   map[string]string
	QueryTemplate   map[string]string
	OtherTemplate   map[string]string
	NormalFile      map[string]string
}

func NewFileMap() *FileMap {
	return &FileMap{
		PackageTemplate: make(map[string]string),
		ModelTemplate:   make(map[string]string),
		QueryTemplate:   make(map[string]string),
		OtherTemplate:   make(map[string]string),
		NormalFile:      make(map[string]string),
	}
}

// 递归扫描目录
func scanDirs(dirName string, fm *FileMap) {
	files, err := ioutil.ReadDir(dirName)
	if err != nil {
		log.Println(err)
	}
	parentPath := dirName
	for _, file := range files {
		fullName := path.Join(parentPath, file.Name())
		if file.IsDir() {
			scanDirs(fullName, fm)
			continue
		}

		fileName := file.Name()
		if !strings.HasSuffix(fileName, ".tpl") {
			fm.NormalFile[fullName] = fileName
			continue
		}

		if strings.Index(fileName, "package") >= 0 {
			fm.PackageTemplate[fullName] = fileName
		} else if strings.Index(fileName, "model") >= 0 {
			fm.ModelTemplate[fullName] = fileName
		} else if strings.Index(fileName, "query") >= 0 {
			fm.QueryTemplate[fullName] = fileName
		} else {
			fm.OtherTemplate[fullName] = fileName
		}
	}
}

func removeFirstDir(file string) string {
	n := strings.Index(file, "/")
	if n < 0 {
		return file
	}
	return file[n+1:]
}

func createDstDir(file string) {
	dstPath, _ := path.Split(file)
	os.MkdirAll(dstPath, os.ModePerm)
}

func generateCode(cfgFile string, fm *FileMap, outPath string) {
	pkg := getTestPkg(cfgFile)
	os.MkdirAll(outPath, os.ModePerm)

	// 普通文件，直接复制
	for k, _ := range fm.NormalFile {
		outFile := path.Join(outPath, removeFirstDir(k))
		createDstDir(outFile)
		err := os.Link(k, outFile)
		if err != nil {
			fmt.Printf("普通%v生成%v失败，错误:%v\n", k, outFile, err)
		}
	}

	// 一般模板文件
	for k, _ := range fm.OtherTemplate {
		outFile := path.Join(outPath, strings.Replace(removeFirstDir(k), ".tpl", "", -1))
		createDstDir(outFile)
		err := tpl.GenerateTemplate(k, outFile, pkg)
		if err != nil {
			fmt.Printf("模板%v生成%v失败，错误:%v\n", k, outFile, err)
		}
	}

	// 包模板文件
	for k, _ := range fm.PackageTemplate {
		outFile := strings.Replace(strings.Replace(removeFirstDir(k), ".tpl", "", -1), "@package", pkg.Name, -1)
		outFile = path.Join(outPath, outFile)
		createDstDir(outFile)
		err := tpl.GenerateTemplate(k, outFile, pkg)
		if err != nil {
			fmt.Printf("模板%v生成%v失败，错误:%v\n", k, outFile, err)
		}
	}

	// model模板文件
	for k, _ := range fm.ModelTemplate {
		for _, m := range pkg.Model {
			outFile := strings.Replace(strings.Replace(removeFirstDir(k), ".tpl", "", -1), "@model", m.Name, -1)
			outFile = path.Join(outPath, outFile)
			createDstDir(outFile)
			pkg.CurrentModel = m
			err := tpl.GenerateTemplate(k, outFile, pkg)
			if err != nil {
				fmt.Printf("模板%v生成%v失败，错误:%v\n", k, outFile, err)
			}

			// query模板文件
			for k, _ := range fm.QueryTemplate {
				for _, q := range m.Query {
					queryFile := strings.Replace(strings.Replace(removeFirstDir(k), ".tpl", "", -1), "@query", fmt.Sprintf("%v_%v", m.Name, q.Name), -1)
					queryFile = path.Join(outPath, queryFile)
					createDstDir(queryFile)
					pkg.CurrentQuery = q
					err := tpl.GenerateTemplate(k, queryFile, pkg)
					if err != nil {
						fmt.Printf("模板%v生成%v失败，错误:%v\n", k, queryFile, err)
					}
				}
			}
		}
	}
}

//
//func generateTemplate(tmplFile, outFile string, data interface{}) error {
//	templ, err := ioutil.ReadFile(tmplFile)
//	if err != nil {
//		return err
//	}
//	//t := template.Must(template.New("escape").Funcs(template.FuncMap{
//	//	"formatName": formatName,
//	//}).Parse(string(templ)))
//
//	t := template.Must(template.New("escape").Parse(string(templ)))
//	nb := bytes.NewBuffer(nil)
//	if err := t.Execute(nb, data); err != nil {
//		return err
//	}
//
//	err = ioutil.WriteFile(outFile, nb.Bytes(), os.ModePerm)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func getTestPkg(cfgFile string) *tpl.PackageCfg {
	//pkg := tpl.NewPackageCfg()
	////viper.SetConfigFile(cfgFile)
	//
	//config := viper.New()
	//config.SetConfigFile(cfgFile)
	////config.AddConfigPath(path)     //设置读取的文件路径
	////config.SetConfigName("cr_all") //设置读取的文件名
	////config.SetConfigType("yaml")   //设置文件的类型
	////尝试进行配置读取
	//if err := config.ReadInConfig(); err != nil {
	//	panic(err)
	//}

	pkg, err := tpl.ParseConfigure(cfgFile)
	if err != nil {
		panic(err)
	}

	return pkg

	pkg.Name = "example"

	pkg.Type = append(pkg.Type, &tpl.TypeCfg{
		Type:     "Decimal",
		GoType:   "decimal.Decimal",
		GoImport: "github.com/zj/decimal",
	})

	pkg.Column = append(pkg.Column, &tpl.ColumnCfg{
		Name:     "id",
		Type:     "int",
		Comment:  "唯一标识",
		GoType:   "decimal.Decimal",
		GoImport: "github.com/zj/decimal",
	})

	pkg.Column = append(pkg.Column, &tpl.ColumnCfg{
		Name:    "name",
		Type:    "string",
		Comment: "姓名",
	})

	pkg.Model = append(pkg.Model, &tpl.ModelCfg{
		Name:    "user",
		Comment: "用户表",
		Columns: []*tpl.ColumnCfg{
			{
				Name:          "id",
				Type:          "int64",
				Comment:       "唯一标识",
				Primary:       true,
				NotNull:       true,
				Unique:        true,
				AutoIncrement: true,
				Default:       "NULL",
				GoType:        "int",
				GoImport:      "",
				FieldType:     0,
			},
			{
				Name:    "name",
				Type:    "string",
				Comment: "姓名",
			},
		},
		Indexes: []*tpl.IndexCfg{{
			Name:    "idx_id",
			Keys:    []string{"id"},
			Comment: "id的索引",
			Type:    "index",
		}},
	})

	pkg.Model[0].Query = append(pkg.Model[0].Query, &tpl.QueryCfg{
		Name:    "queryUserInfo",
		Comment: "query user info",
		SQLstr:  "select * from user where id=$id",
		In: &tpl.StructCfg{
			Name: "inPara",
			Columns: []*tpl.ColumnCfg{
				&tpl.ColumnCfg{
					Name:    "page",
					Type:    "int",
					Comment: "分页数",
				},
				&tpl.ColumnCfg{
					Name:    "size",
					Type:    "int",
					Comment: "分页大小",
				},
			},
		},
		Out: &tpl.StructCfg{
			Name: "outPara",
			Columns: []*tpl.ColumnCfg{
				{
					Name:    "count",
					Type:    "int",
					Comment: "总数",
				},
				{
					Name:      "users",
					Type:      "objArr",
					FieldType: 1,
					Children: map[string]*tpl.ColumnCfg{
						"id": {
							Name:    "id",
							Type:    "int",
							Comment: "唯一标识",
						},
						"name": {
							Name:    "name",
							Type:    "string",
							Comment: "姓名",
						},
					},
				},
			},
		},
	})

	return pkg
}
