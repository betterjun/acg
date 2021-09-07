package model

import "{{ .Name }}/pkg/db"
import "fmt"

func Setup() {
	{{range $, $v := .Model}}
        {{- with  $v -}}
    	db.GetDB().AutoMigrate(&{{ formatName $v.Name }}{})
    	{{- end }}
    {{end}}

// 设置链接数量
	sqlDB, _ := db.GetDB().DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	dbh := db.GetDB()
	user := &User{}
	err := user.Create(dbh)
	if err != nil {
		panic(err)
	}

	user.Name = "zj"
	err = user.Update(dbh)
	if err != nil {
		panic(err)
	}

	nu, err := user.FindBySex(dbh, 0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", nu)

	nu, err = user.FindByID(dbh, 1)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", nu)

	nu, err = user.FindByAge(dbh, 0)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", nu)
}
