/**
 * @Author zj
 * @create 2021/09/05 4:44 PM
 */

package model

import (
	"gorm.io/gorm"
)


{{$tableName := formatName .CurrentQuery.Name }}

// {{ .CurrentQuery.Comment }}
type {{ $tableName }}Context struct {

}




{{define "field.tmpl"}}
    {{- if eq .FieldType 0 -}}
        {{- formatName .Name }} {{.Type}} `json:"{{.Name}}" form:"{{.Name}}"` // {{.Comment -}}
    {{- else if eq .FieldType 1 -}}
        {{$lt := len .Type }}
        {{- if eq $lt 0 -}}
            {{- formatName .Name }} []struct{
                {{ range $k, $v := .Children }}
                    {{- template "field.tmpl" $v -}}
                {{- end -}}
            } `json:"{{.Name}}" form:"{{.Name}}"` // {{.Comment -}}
        {{- else -}}
             {{- formatName .Name }} []{{.Type}} `json:"{{.Name}}" form:"{{.Name}}"` // {{.Comment -}}
        {{- end -}}
    {{- else if eq .FieldType 2 -}}
        {{- formatName .Name }} []struct{
            {{ range $k, $v := .Children }}
                {{- template "field.tmpl" $v -}}
            {{- end -}}
        } `json:"{{.Name}}" form:"{{.Name}}"` // {{.Comment -}}
    {{- end }}
{{ end}}



type {{ .CurrentQuery.In.Name }} struct {
{{range $k, $v := .CurrentQuery.In.Columns -}}
    {{- template "field.tmpl" $v -}}
{{- end -}}
}


type {{ .CurrentQuery.Out.Name }} struct {
	Count  int `json:"count" form:"count"`   // 总数
	Offset int `json:"offset" form:"offset"` // 偏移量
	Size   int `json:"size" form:"size"`     // 分页大小
	Result []struct {     // 结果数据
{{range $k, $v := .CurrentQuery.Out.Columns -}}
    {{- template "field.tmpl" $v -}}
{{- end -}}
    } `json:"result" form:"result"`
}




// {{ .CurrentQuery.Comment }}
func (q *{{ $tableName }}Context) {{ $tableName }}(db *gorm.DB, in *{{ .CurrentQuery.In.Name }}) (out *{{ .CurrentQuery.Out.Name }}, err error) {
    var mm interface{} = q
    out =new({{ .CurrentQuery.Out.Name }})

{{ $res := .CurrentQuery.FormatSQL3 }}
    err = db.Raw("{{$res.SQL}}", {{- range $k, $v := $res.Names }}
                               in.{{ formatName $v }},
                               {{- end}}).Scan(out).Error

    return
}

