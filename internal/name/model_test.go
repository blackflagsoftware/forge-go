package name

import "testing"

func TestName_DetermineAbbr(t *testing.T) {
	type fields struct {
		Lower string
	}
	type args struct {
		knownAliases []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			"try 1, no known alias",
			fields{Lower: "test_name_1"},
			args{knownAliases: []string{}},
			"tes",
		},
		{
			"try 2, known alias",
			fields{Lower: "test_name_2"},
			args{knownAliases: []string{"tes"}},
			"te",
		},
		{
			"try 3, known alias",
			fields{Lower: "test_name"},
			args{knownAliases: []string{"tes", "te"}},
			"tn",
		},
		{
			"try 4, known alias",
			fields{Lower: "test_name1"},
			args{knownAliases: []string{"tes", "te", "tn"}},
			"tena",
		},
		{
			"try 5, known alias",
			fields{Lower: "test_name"},
			args{knownAliases: []string{"tes", "te", "tn", "tena"}},
			"test_name",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			n := &Name{
				Lower: tt.fields.Lower,
			}
			if got := n.DetermineAbbr(tt.args.knownAliases); got != tt.want {
				t.Errorf("Name.DetermineAbbr() = %v, want %v", got, tt.want)
			}
		})
	}
}
