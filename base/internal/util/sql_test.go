package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSearchBuilder_AppendCompare(t *testing.T) {
	type fields struct {
		Params []string
		Values []interface{}
	}
	type args struct {
		param   string
		compare string
		value   interface{}
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantParam []string
	}{
		{
			"success - 1",
			fields{
				[]string{},
				[]interface{}{},
			},
			args{
				"id",
				"=",
				1,
			},
			[]string{"id = ?"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SearchBuilder{
				Params: tt.fields.Params,
				Values: tt.fields.Values,
			}
			s.AppendCompare(tt.args.param, tt.args.compare, tt.args.value)
			assert.Equal(t, tt.wantParam, s.Params, "params are not equal")
		})
	}
}

func TestSearchBuilder_AppendLike(t *testing.T) {
	type fields struct {
		Params []string
		Values []interface{}
	}
	type args struct {
		param string
		value string
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantParam []string
	}{
		{
			"successful - 1",
			fields{
				[]string{},
				[]interface{}{},
			},
			args{
				"addr",
				"street",
			},
			[]string{"addr LIKE %street%"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SearchBuilder{
				Params: tt.fields.Params,
				Values: tt.fields.Values,
			}
			s.AppendLike(tt.args.param, tt.args.value)
			assert.Equal(t, tt.wantParam, s.Params, "params are not equal")
		})
	}
}

func TestSearchBuilder_AppendNull(t *testing.T) {
	type fields struct {
		Params []string
		Values []interface{}
	}
	type args struct {
		param    string
		wantNull bool
	}
	tests := []struct {
		name      string
		fields    fields
		args      args
		wantParam []string
	}{
		{
			"successful - null",
			fields{
				[]string{},
				[]interface{}{},
			},
			args{
				"name",
				true,
			},
			[]string{"name IS NULL"},
		},
		{
			"successful - not null",
			fields{
				[]string{},
				[]interface{}{},
			},
			args{
				"name",
				false,
			},
			[]string{"name IS NOT NULL"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SearchBuilder{
				Params: tt.fields.Params,
				Values: tt.fields.Values,
			}
			s.AppendNull(tt.args.param, tt.args.wantNull)
			assert.Equal(t, tt.wantParam, s.Params, "params are not equal")
		})
	}
}

func TestSearchBuilder_String(t *testing.T) {
	type fields struct {
		Params []string
		Values []interface{}
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"successful",
			fields{
				[]string{"name IS NOT NULL", "addr LIKE %home%"},
				[]interface{}{},
			},
			"WHERE name IS NOT NULL\n\t\tAND addr LIKE %home%",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SearchBuilder{
				Params: tt.fields.Params,
				Values: tt.fields.Values,
			}
			got := s.String()
			assert.Equal(t, tt.want, got, "output is not equal")
		})
	}
}

func TestBuildSearchStringWithTenant(t *testing.T) {
	type args struct {
		search []ParamSearch
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 []interface{}
	}{
		{
			"successful - LIKE",
			args{
				[]ParamSearch{
					{
						Column:  "name",
						Compare: "LIKE",
						Value:   "test",
					},
				},
			},
			"WHERE name LIKE %test%",
			[]interface{}{"tenant_01"},
		},
		{
			"successful - COMPARE",
			args{
				[]ParamSearch{
					{
						Column:  "name",
						Compare: "=",
						Value:   "test",
					},
				},
			},
			"WHERE name = ?",
			[]interface{}{"test"},
		},
		{
			"successful - NULL",
			args{
				[]ParamSearch{
					{
						Column:  "name",
						Compare: "NULL",
						Value:   "",
					},
				},
			},
			"WHERE name IS NULL",
			[]interface{}{},
		},
		{
			"successful - NOT NULL",
			args{
				[]ParamSearch{
					{
						Column:  "name",
						Compare: "NOT NULL",
						Value:   "",
					},
				},
			},
			"WHERE name IS NOT NULL",
			[]interface{}{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := BuildSearchString(tt.args.search)
			assert.Equal(t, tt.want, got, "where stmt is not equal")
			assert.Equal(t, tt.want1, got1, "values are equal")
		})
	}
}
