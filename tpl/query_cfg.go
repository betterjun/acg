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
	Pager   bool       `yaml:"pager,omitempty"`   // 当sql为查询时，且输出为数组时，是否启用分页，默认不启用
	In      *StructCfg `yaml:"inputs,omitempty"`  // 查询的输入参数
	Out     *StructCfg `yaml:"outputs,omitempty"` // 查询的输出参数
}

// query按查询名进行组织
type QueryMap map[string]*QueryCfg

type SqlRes struct {
	SQL   string
	Names []string
}

func (q *QueryCfg) FormatSQL() SqlRes {
	sql, names := handleSql(q.SQLstr)
	return SqlRes{
		SQL: sql, Names: names,
	}
}

func (q *QueryCfg) FormatCountSQL() SqlRes {
	if !q.Out.Array {
		return SqlRes{}
	}

	/*
		先换掉select字段
		再去掉limit和offset语句
		再替换字段
	*/
	sqlCount := strings.ToLower(q.SQLstr)
	sqlCount = convertSpaces(sqlCount)
	p1 := strings.Index(sqlCount, "select ")
	p2 := strings.Index(sqlCount, " from ")
	if p1 < 0 || p2 < 7 {
		// 语法错误了，暂时返回空
		return SqlRes{}
	}

	str := sqlCount[p1+7 : p2] //起始位置要去掉"select "
	sqlCount = strings.Replace(sqlCount, str, "count(1)", 1)

	p3 := strings.Index(sqlCount, " limit ")
	if p3 > 0 {
		sqlCount = sqlCount[:p3]
	}
	sql, names := handleSql(sqlCount)
	return SqlRes{
		SQL: sql, Names: names,
	}
}

func convertSpaces(s string) (ret string) {
	ret = s
	ret = strings.Replace(ret, "\t", " ", -1)
	ret = strings.Replace(ret, "\n", " ", -1)
	ret = strings.Replace(ret, "\v", " ", -1)
	ret = strings.Replace(ret, "\f", " ", -1)
	ret = strings.Replace(ret, "\r", " ", -1)

	return ret
}

func handleSql(sqlInput string) (sql string, names []string) {
	sql = sqlInput
	names = make([]string, 0)
	s := -1
	e := -1
	for i, c := range sqlInput {
		if c == rune('$') {
			fmt.Println('$', i, c)
			s = i
		} else {
			if s == -1 {
				// do nothing
			} else {
				if isSqlSplitter(c) || i == len(sqlInput) {
					e = i

					name := sqlInput[s+1 : e]
					sql = strings.Replace(sql, "$"+name, "?", -1)
					fmt.Println(sql)
					fmt.Println(sqlInput)
					names = append(names, name)
					// reset
					s = -1
					e = -1
				}
			}
		}
	}

	if s != -1 {
		name := sqlInput[s+1:]
		sql = strings.Replace(sql, "$"+name, "?", -1)
		fmt.Println(sql)
		fmt.Println(sqlInput)
		names = append(names, name)
		// reset
		s = -1
	}

	return sql, names
}

func isSqlSplitter(c rune) bool {
	if unicode.IsSpace(c) || c == ',' || c == ')' || c == ']' || c == '=' || c == '!' || c == '<' || c == '>' {
		return true
	} else {
		return false
	}
}
