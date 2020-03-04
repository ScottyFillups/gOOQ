package modelgen

const modelTemplate = `
// THIS FILE WAS AUTOGENERATED - ANY EDITS TO THIS WILL BE LOST WHEN IT IS REGENERATED

{{ $schema := .Schema }}
package {{ .Package }}

import "github.com/lumina-tech/gooq/pkg/gooq"

{{ range $_, $table := .Tables }}
type {{ $table.ModelType }} struct {
  {{ range $_, $f := $table.Fields -}}
  {{ snakeToCamelID $f.Name }} {{ $f.Type }} ` + "`db:\"{{ $f.Name }}\" json:\"{{ $f.Name }}\"`" + `
  {{ end }}
}
{{ end }}
`

const tableTemplate = `
// THIS FILE WAS AUTOGENERATED - ANY EDITS TO THIS WILL BE LOST WHEN IT IS REGENERATED

{{ $schema := .Schema }}
package {{ .Package }}

import "github.com/lumina-tech/gooq/pkg/gooq"

{{ range $_, $table := .Tables -}}

type {{ $table.TableType }}Constraints struct {
  {{ range $_, $f := $table.Constraints -}}
  {{ snakeToCamel $f.Name }} gooq.DatabaseConstraint
  {{ end }}
}

type {{ $table.TableType }} struct {
	gooq.TableImpl
	Asterisk gooq.StringField
  {{ range $_, $f := $table.Fields -}}
  {{ snakeToCamel $f.Name }} gooq.{{ $f.GooqType }}Field
  {{ end }}
  Constraints *{{ $table.TableType }}Constraints
}

func new{{ capitalize $table.TableType }}Constraints() *{{ $table.TableType }}Constraints {
  constraints := &{{ $table.TableType }}Constraints{}
  {{ range $_, $f := $table.Constraints -}}
  constraints.{{ snakeToCamel $f.Name }} = gooq.DatabaseConstraint{
    Name: "{{$f.Name}}",
    Predicate: null.NewString("{{$f.Predicate.String}}", {{$f.Predicate.Valid}}),
  }
  {{ end -}}
  return constraints
}

func new{{ capitalize $table.TableType }}() *{{ $table.TableType }} {
  instance := &{{ $table.TableType }}{}
	instance.Initialize("{{ $schema }}", "{{ $table.TableName }}")
	instance.Asterisk = gooq.NewStringField(instance, "*")
  {{ range $_, $f := $table.Fields -}}
  instance.{{ snakeToCamelID $f.Name }} = gooq.New{{ $f.GooqType }}Field(instance, "{{ $f.Name }}")
  {{ end -}}
  instance.Constraints = new{{ $table.ModelType }}Constraints()
  return instance
}

func (t *{{ $table.TableType }}) As(alias string) *{{ $table.TableType }} {
  instance := new{{ $table.ModelType }}()
  instance.TableImpl = *instance.TableImpl.As(alias)
  return instance
}

func (t *{{ $table.TableType }}) GetColumns() []gooq.Expression {
	return []gooq.Expression{
  {{ range $_, $f := $table.Fields -}}
  t.{{ snakeToCamelID $f.Name }},
  {{ end -}}
  }
}

func (t *{{ $table.TableType }}) ScanRow(
	db gooq.DBInterface, stmt gooq.Fetchable,
) (*{{ $table.QualifiedModelType }}, error) {
	result := {{ $table.QualifiedModelType }}{}
	if err := gooq.ScanRow(db, stmt, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (t *{{ $table.TableType }}) ScanRows(
	db gooq.DBInterface, stmt gooq.Fetchable,
) ([]{{ $table.QualifiedModelType }}, error) {
	results := []{{ $table.QualifiedModelType }}{}
	if err := gooq.ScanRows(db, stmt, &results); err != nil {
		return nil, err
	}
	return results, nil
}

var {{ $table.TableSingletonName }} = new{{ capitalize $table.TableType }}()
{{ end }}
`
