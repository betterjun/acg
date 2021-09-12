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

	"github.com/betterjun/acg/tpl"

	"github.com/spf13/viper"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "init a yaml configure file for your code generator",
	Long: `init a yaml configure file for your code generator. For example:
init example.yaml
.`,
	Run: func(cmd *cobra.Command, args []string) {
		defer fmt.Println("init called")
		runInitCmd(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	//initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	initCmd.Flags().BoolP("ok", "o", false, "Help message for ok")
}

func runInitCmd(cmd *cobra.Command, args []string) {
	fmt.Println("init called")

	fmt.Printf("init args=%v\n", args)

	//viper.Set("package.name", "example")
	//
	//viper.Set("type", []map[string]string{{"type": "Decimal", "goType": "decimal.Decimal", "goImport": "github.com/zj/decimal"}})
	//
	//dictArr := make([]map[string]interface{}, 0)
	//dictArr = append(dictArr, map[string]interface{}{
	//	"id":         map[string]interface{}{"type": "int", "comment": "唯一标识", "primary": true, "notnull": true, "unique": true, "auto_increment": true, "default": 0},
	//	"created_at": map[string]interface{}{"type": "int64", "comment": "创建时间", "default": "NULL"},
	//	"updated_at": "created_at",
	//	"sex":        []interface{}{"string", "性别", "默认为空", false, true, false, true},
	//})
	//viper.Set("dict", dictArr)

	pkg := tpl.NewPackageCfg()
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

	viper.Set("package", pkg)

	viper.WriteConfigAs("example.yaml") // will error since it has already been written
}
