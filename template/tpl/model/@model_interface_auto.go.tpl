/**
 * @Author zj
 * @create 2020/12/14 4:44 PM
 */

package model

import (
	"gorm.io/gorm"
)


{{$tableName := formatName .CurrentModel.Name }}


// auto generated interface, implement them in struct {{$tableName}} in this file.
{{range $k, $v := .CurrentModel.Columns}}
    {{- with  $v }}
type {{$tableName}}FindBy{{formatName .Name}} interface {
    {{ $tableName }}FindBy{{formatName .Name}}(db *gorm.DB, {{formatName .Name}} {{.GoTypeString}}) (ret *{{ $tableName }}, err error)
}
    {{- end}}
{{end}}


// append your interface here
