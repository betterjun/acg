package model

import "{{ .Name }}/pkg/db"
import "fmt"

// 重新生成测试代码
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

	{
        user := &User{}
        user.Name = "zj"
        err := user.Create(dbh)
        if err != nil {
            panic(err)
        }

	    // 增加
        cuc := new(UserContext)
        cucout, err := cuc.CreateUser(dbh, &CreateUserIN{
        	Name        :"zj",
        	Age         :30,
        	Sex         :1,
        	Introduction:"intro",
        	Money       :2,
        })
        if err!=nil {
            panic(err)
        }

        // 修改
        uucout, err := cuc.UpdateUser(dbh, &UpdateUserIN{
        	Name        :"zj2",
        	Age         :10,
        	Sex         :0,
        	Introduction:"简介",
        	Money       :3,
        	ID          :user.ID,
        })
        if err!=nil {
            panic(err)
        }

        // 删除
        ducout, err := cuc.DeleteUser(dbh, &DeleteUserIN{
        	ID          :user.ID,
        })
        if err!=nil {
            panic(err)
        }

        fmt.Println(cucout, uucout, ducout)
	}




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
