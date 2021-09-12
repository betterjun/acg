/**
 * @Author zj
 * @create 2020/12/14 4:44 PM
 */

package model

import (
	"gorm.io/gorm"
	{{ .CurrentModel.GoImportsString }}
)


{{$tableName := formatName .CurrentModel.Name }}

// {{ .CurrentModel.Comment }}
type {{ $tableName }} struct {
{{- range $k, $v := .CurrentModel.Columns }}
    {{- with  $v -}}

    {{$le:= len .Default}}
    {{ formatName .Name}} {{.GoTypeString}}  `json:"{{.Name}}" form:"{{.Name}}" gorm:"{{.Name}} {{- if .Primary -}} ;primary {{- end -}}   {{- if .NotNull -}} ;notnull {{- end -}}   {{- if .Unique -}} ;unique {{- end -}} {{- if gt $le 0 -}} ;default= {{- .Default -}} {{- end -}}"` // {{ .Comment -}}
    {{- end -}}
{{end}}
}

// {{ .CurrentModel.Name }}查询上下文
type {{ $tableName }}Context struct {

}


/*
func Migrate{{ $tableName }}(db *gorm.DB) error {
	err := db.AutoMigrate(&{{ $tableName }}{}).Error
	if err!=nil {
    	return err
	}

{{- range $k, $v := .CurrentModel.Indexes }}
    {{- with  $v -}}
    {{- if eq .Type "unique" }}
        db.Model(&{{ $tableName }}{}).AddUniqueIndex("&{{ .Name }}",
        	{{- range $k1, $v1 := .Keys }}
        	"{{ $v1 }}",
            {{end}}
        	)
    {{else}}
    	db.Model(&{{ $tableName }}{}).AddIndex("&{{ .Name }}",
    	{{- range $k1, $v1 := .Keys }}
    	"{{ $v1 }}",
        {{end}}
    	)
    {{- end -}}
        {{- end -}}
{{end}}

}
*/

func (m *{{ $tableName }}) TableName() string {
	return "{{ .CurrentModel.Name }}"
}

func (m *{{ $tableName }}) Create(db *gorm.DB) error {
	return db.Model(&{{$tableName}}{}).Create(m).Error
}

func (m *{{ $tableName }}) Delete(db *gorm.DB) error {
	return db.Model(&{{$tableName}}{}).Where("id=?", m.ID).Delete(m).Error
}

func (m *{{ $tableName }}) Update(db *gorm.DB) error {
    mm := new({{ $tableName }})
    err := db.Model(&{{$tableName}}{}).Where("id=?", m.ID).First(mm).Error
    if err!=nil {
        return err
    }

    // compare the fields
{{- range $k, $v := .CurrentModel.Columns }}
    {{- with  $v }}
    if mm.{{formatName .Name}} != m.{{formatName .Name}} {
        mm.{{formatName .Name}} = m.{{formatName .Name}}
    }
    {{- end}}
{{- end}}

	return db.Model(&{{$tableName}}{}).Where("id=?", m.ID).Updates(mm).Error
}


{{range $k, $v := .CurrentModel.Columns}}
    {{- with  $v }}
func (m *{{ $tableName }}) FindBy{{formatName .Name}}(db *gorm.DB, {{formatName .Name}} {{.GoTypeString}}) (ret *{{ $tableName }}, err error) {
    var mm interface{} = m
    /*
	if f, ok := mm.({{$tableName}}FindBy{{formatName .Name}}); ok && !isNil(f) {
		return f.{{ $tableName }}FindBy{{formatName .Name}}(db, {{formatName .Name}})
	}
	*/
	if f, ok := mm.({{$tableName}}FindBy{{formatName .Name}}); ok && f != nil {
		return f.{{ $tableName }}FindBy{{formatName .Name}}(db, {{formatName .Name}})
	}

    ret = new({{ $tableName }})
    err = db.Model(&{{ $tableName }}{}).Where("{{.Name}}=?", {{formatName .Name}}).First(ret).Error
	return ret, err
}
    {{- end}}
{{end}}
