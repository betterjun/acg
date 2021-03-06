/**
 * @Author zj
 * @create 2021/09/05 4:44 PM
 */

package model

import (
	"gorm.io/gorm"
	{{ .CurrentQuery.GoImportsString }}
)


{{$modelName := formatName .CurrentModel.Name }}
{{$queryName := formatName .CurrentQuery.Name }}
{{$contextName := printf "%vContext" $modelName }}
{{$queryInterface := printf "%v%v" $contextName $queryName}}


{{define "field.tmpl"}}
    {{- if eq .FieldType 0 -}}
        {{- formatName .Name }} {{.GoTypeString}} `json:"{{.Name}}" form:"{{.Name}}" gorm:"{{.Name}}"` // {{.Comment -}}
    {{- else if eq .FieldType 1 -}}
        {{$lt := len .Type }}
        {{- if eq $lt 0 -}}
            {{- formatName .Name }} []struct{
                {{ range $k, $v := .Children }}
                    {{- template "field.tmpl" $v -}}
                {{- end -}}
            } `json:"{{.Name}}" form:"{{.Name}}"` // {{.Comment -}}
        {{- else -}}
             {{- formatName .Name }} []{{.GoTypeString}} `json:"{{.Name}}" form:"{{.Name}}" gorm:"{{.Name}}"` // {{.Comment -}}
        {{- end -}}
    {{- else if eq .FieldType 2 -}}
        {{- formatName .Name }} []struct{
            {{ range $k, $v := .Children }}
                {{- template "field.tmpl" $v -}}
            {{- end -}}
        } `json:"{{.Name}}" form:"{{.Name}}"` // {{.Comment -}}
    {{- end }}
{{ end}}

{{$inStructName := formatName .CurrentQuery.In.Name }}
{{$outStructName := formatName .CurrentQuery.Out.Name }}

type {{ $inStructName }} struct {
{{range $k, $v := .CurrentQuery.In.Columns -}}
    {{- template "field.tmpl" $v -}}
{{- end -}}
}

{{ if eq .CurrentQuery.Out.Array false }}

    {{ if eq .CurrentQuery.Pager false}}
        type {{ $outStructName }} struct {
            {{range $k, $v := .CurrentQuery.Out.Columns -}}
                {{- template "field.tmpl" $v -}}
            {{- end -}}
        }

        // {{ .CurrentQuery.Comment }}
        func (q *{{ $contextName }}) {{ $queryName }}(db *gorm.DB, in *{{ $inStructName }}) (out *{{ $outStructName }}, err error) {
            var mm interface{} = q
            if f, ok := mm.({{ $queryInterface }}Interface); ok && f != nil {
                return f.{{ $queryInterface }}Impl(db, in)
            }

            out = new({{ $outStructName }})
            {{ $res := .CurrentQuery.FormatSQL }}
            err = db.Raw("{{$res.SQL}}", {{- range $k, $v := $res.Names }}
                                       in.{{ formatName $v }},
                                       {{- end}}).Scan(out).Error

            return
        }

    {{ else }}
        type {{ $outStructName }}Data struct {
            {{range $k, $v := .CurrentQuery.Out.Columns -}}
                {{- if eq $v.Name "result" -}}
                    {{range $k1, $v1 := $v.Children -}}
                        {{- template "field.tmpl" $v1 -}}
                    {{ end -}}
                {{- end -}}
            {{- end -}}
        }

        type {{ $outStructName }} struct {
            {{range $k, $v := .CurrentQuery.Out.Columns -}}
                {{ if eq $v.Name "result" }}
                    {{- formatName $v.Name }} []{{ $outStructName }}Data `json:"{{$v.Name}}" form:"{{$v.Name}}"` // {{$v.Comment -}}
                {{ else }}
                    {{- template "field.tmpl" $v -}}
                {{ end -}}
            {{- end }}
        }


        // {{ .CurrentQuery.Comment }}
        func (q *{{ $contextName }}) {{ $queryName }}(db *gorm.DB, in *{{ $inStructName }}) (out *{{ $outStructName }}, err error) {
            var mm interface{} = q
            if f, ok := mm.({{ $queryInterface }}Interface); ok && f != nil {
                return f.{{ $queryInterface }}Impl(db, in)
            }

            out = new({{ $outStructName }})
            out.Result = make([]{{ $outStructName }}Data,0)

            {{ $res := .CurrentQuery.FormatSQL }}
            err = db.Raw("{{$res.SQL}}", {{- range $k, $v := $res.Names }}
                                       in.{{ formatName $v }},
                                       {{- end}}).Scan(&out.Result).Error
            if err !=nil {
                return nil, err
            }

            {{ $res := .CurrentQuery.FormatCountSQL }}
            err = db.Raw("{{$res.SQL}}", {{- range $k, $v := $res.Names }}
                                       in.{{ formatName $v }},
                                       {{- end}}).Count(&out.Count).Error

            return
        }

    {{ end }}


    type {{ $queryInterface }}Interface interface {
        {{ $queryInterface }}Impl(db *gorm.DB, in *{{ $inStructName }}) (out *{{ $outStructName }}, err error)
    }

{{ else }}
    type {{ $outStructName }} struct {
    {{range $k, $v := .CurrentQuery.Out.Columns -}}
        {{- template "field.tmpl" $v -}}
    {{- end -}}
    }

// {{ .CurrentQuery.Comment }}
func (q *{{ $contextName }}) {{ $queryName }}(db *gorm.DB, in *{{ $inStructName }}) (out []{{ $outStructName }}, err error) {
    var mm interface{} = q
	if f, ok := mm.({{ $queryInterface }}Interface); ok && f != nil {
		return f.{{ $queryInterface }}Impl(db, in)
	}

    out = make([]{{ $outStructName }},0)

    {{ $res := .CurrentQuery.FormatSQL }}
    err = db.Raw("{{$res.SQL}}", {{- range $k, $v := $res.Names }}
                               in.{{ formatName $v }},
                               {{- end}}).Scan(&out).Error
    if err !=nil {
        return nil, err
    }

    return out, err
}



type {{ $queryInterface }}Interface interface {
    {{ $queryInterface }}Impl(db *gorm.DB, in *{{ $inStructName }}) (out []{{ $outStructName }}, err error)
}

{{ end }}

