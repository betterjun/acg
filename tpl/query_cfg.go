package tpl

import (
	"fmt"
	"strings"
	"unicode"
)

// 查询
type QueryCfg struct {
	Name    string     // 查询的名字
	Comment string     `yaml:"comment,omitempty"` // 查询的说明
	SQLstr  string     `yaml:"sql,omitempty"`     // 查询用到的sql语句
	In      *StructCfg `yaml:"inputs,omitempty"`  // 查询的输入参数
	Out     *StructCfg `yaml:"outputs,omitempty"` // 查询的输出参数
}

// query按查询名进行组织
type QueryMap map[string]*QueryCfg

func (q *QueryCfg) FormatSQL2() (sql string) {
	return q.SQLstr
}

func (q *QueryCfg) FormatSQL() (sql string, names []string) {
	sql = q.SQLstr
	names = make([]string, 0)
	//fieldNames := make([]string, 0, len(q.In.Columns))
	//for _, v := range q.In.Columns { // 只查看一级字段
	//	//if strings.ContainsAny()
	//	fieldNames = append(fieldNames, v.Name)
	//}

	//// 依次查找$符号
	//n := strings.Index(sql, "$")
	//if n<0 {
	//	return sql, names
	//}

	s := 0
	e := 0
	for i, c := range q.SQLstr {
		if c == rune('$') {
			s = i
		} else {
			if unicode.IsSpace(c) {
				e = i

				name := q.SQLstr[s : e+1]
				sql = strings.Replace(sql, name, "?", -1)
				names = append(names, name)
				// reset
				s = 0
				e = 0
			}
		}
	}

	return
}

type SqlRes struct {
	SQL   string
	Names []string
}

func (q *QueryCfg) FormatSQL3() SqlRes {
	sql := q.SQLstr
	names := make([]string, 0)
	s := -1
	e := -1
	for i, c := range q.SQLstr {
		if c == rune('$') {
			fmt.Println('$', i, c)
			s = i
		} else {
			if s == -1 {
				// do nothing
			} else {
				if unicode.IsSpace(c) || i == len(q.SQLstr) {
					e = i

					name := q.SQLstr[s+1 : e]
					sql = strings.Replace(sql, "$"+name, "?", -1)
					fmt.Println(sql)
					fmt.Println(q.SQLstr)
					names = append(names, name)
					// reset
					s = -1
					e = -1
				}
			}
		}
	}

	if s != -1 {
		name := q.SQLstr[s+1:]
		sql = strings.Replace(sql, "$"+name, "?", -1)
		fmt.Println(sql)
		fmt.Println(q.SQLstr)
		names = append(names, name)
		// reset
		s = -1
	}

	return SqlRes{
		SQL: sql, Names: names,
	}
}
