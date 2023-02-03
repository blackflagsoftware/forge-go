package main

import "testing"

func Test_getProjectPath(t *testing.T) {
	type args struct {
		goPath    string
		directory string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"full path",
			args{
				"/test/path/here",
				"/test/path/here/src/my/directory/here",
			},
			"my/directory/here",
		},
		{
			"part path",
			args{
				"/test/path/here",
				"my/directory/here",
			},
			"my/directory/here",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getProjectPath(tt.args.goPath, tt.args.directory); got != tt.want {
				t.Errorf("getProjectPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
