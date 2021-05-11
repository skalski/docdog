package helper

import (
	"reflect"
	"testing"
)

func TestBytesToStringArrayByLinebreaks(t *testing.T) {
	type args struct {
		data []byte
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test_happy_path",
			args: args{data: []byte("something\nthat\nis sliceable;")},
			want: []string{
				"something", "that", "is sliceable;",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BytesToStringArrayByLinebreaks(tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("BytesToStringArrayByLinebreaks() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFind(t *testing.T) {
	type args struct {
		slice []string
		val   string
	}
	tests := []struct {
		name  string
		args  args
		want  int
		want1 bool
	}{
		{
			name: "test_happy_path",
			args: args{
				slice: []string{
					"something", "that", "is sliceable;",
				},
				val: "something",
			},
			want:  0,
			want1: true,
		},
		{
			name: "test_happy_path_find_second",
			args: args{
				slice: []string{
					"something", "that", "is sliceable;",
				},
				val: "that",
			},
			want:  1,
			want1: true,
		},
		{
			name: "test_happy_path_find_nothing",
			args: args{
				slice: []string{
					"something", "that", "is sliceable;",
				},
				val: "somethingelse",
			},
			want:  -1,
			want1: false,
		},
		{
			name: "test_happy_path_find_nothing_pt_2",
			args: args{
				slice: []string{
					"something", "that", "is sliceable;",
				},
				val: "some",
			},
			want:  -1,
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := Find(tt.args.slice, tt.args.val)
			if got != tt.want {
				t.Errorf("Find() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("Find() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestGetStringFromQouteLine(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		{
			name:       "test",
			args:       struct{ str string }{str: "'some str'"},
			wantResult: "some str",
		},
		{
			name:       "test from tag",
			args:       struct{ str string }{str: "@DD:DESCRIPTION: 'some str'"},
			wantResult: "some str",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := GetStringFromQouteLine(tt.args.str); gotResult != tt.wantResult {
				t.Errorf("GetStringFromQouteLine() = %v, want %v", gotResult, tt.wantResult)
			}
		})
	}
}

func TestInfoLog(t *testing.T) {
	type args struct {
		msg     string
		source  string
		verbose *bool
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func TestSeparateLineByTags(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SeparateLineByTags(tt.args.line); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SeparateLineByTags() = %v, want %v", got, tt.want)
			}
		})
	}
}
