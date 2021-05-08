package javaLang

import (
	"docdog/src/notations"
	"reflect"
	"testing"
)

func TestJavaVariableHandler(t *testing.T) {
	type args struct {
		line      string
		index     int
		wholeFile []string
	}
	tests := []struct {
		name    string
		args    args
		want    notations.Variable
		wantErr bool
	}{{
		"got err",
		args{
			line:      "",
			index:     0,
			wholeFile: []string{},
		},
		notations.Variable{
			Name:        "",
			Description: "",
			Typ:         "",
			Notnull:     false,
			IsArray:     false,
		},
		true,
	},
		{
			"ignore variable",
			args{
				line:  "private String test;",
				index: 4,
				wholeFile: []string{"/*", "	@DD:IGNORE", "*/", "private String test;"},
			},
			notations.Variable{
				Name:        "test",
				Description: "",
				Typ:         "String",
				Notnull:     false,
				IsArray:     false,
			},
			false,
		},
		{
			"got not null and description",
			args{
				line:      "private String test;",
				index:     5,
				wholeFile: []string{"/*", "@DD:DESCRIPTION 'some var we use'", "@DD:NOTNULL", "*/", "private String test;"},
			},
			notations.Variable{
				Name:        "test",
				Description: "some var we use",
				Typ:         "String",
				Notnull:     true,
				IsArray:     false,
			},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JavaVariableHandler(tt.args.line, tt.args.index, tt.args.wholeFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("JavaVariableHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JavaVariableHandler() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestChckImpl(t *testing.T) {
	type args struct {
		fls []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "test impl Beta class",
			args: args{fls: []string{"import foo:", "public class Alpha implements Beta {"}},
			want: "Beta",
		},
		{
			name: "test no impl",
			args: args{fls: []string{"import foo:", "public class Alpha {"}},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ChckImpl(tt.args.fls); got != tt.want {
				t.Errorf("ChckImpl() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckLangTag(t *testing.T) {
	type args struct {
		wholeFile    []string
		isSpringBoot bool
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test no impl",
			args: args{
				wholeFile:    []string{"import foo;", "public class Alpha {"},
				isSpringBoot: false,
			},
			want: true,
		},
		{
			name: "test spring boot is active impl",
			args: args{
				wholeFile:    []string{"import org.springframework.web.bind.annotation.GetMapping;", "public class Alpha {"},
				isSpringBoot: true,
			},
			want: true,
		},
		{
			name: "test spring boot is active impl but not set in flags",
			args: args{
				wholeFile:    []string{"import org.springframework.web.bind.annotation.GetMapping;", "public class Alpha {"},
				isSpringBoot: false,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckLangTag(tt.args.wholeFile, tt.args.isSpringBoot); got != tt.want {
				t.Errorf("CheckLangTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateApiEp(t *testing.T) {
	type args struct {
		index        int
		wholeFile    []string
		verbose      *bool
		isSpringBoot bool
		wl           *int
	}
	tests := []struct {
		name string
		args args
		want notations.TempEndpoint
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateApiEp(tt.args.index, tt.args.wholeFile, tt.args.verbose, tt.args.isSpringBoot, tt.args.wl); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateApiEp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateArrayType(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateArrayType(tt.args.line); got != tt.want {
				t.Errorf("CreateArrayType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreatePrm(t *testing.T) {
	type args struct {
		l       string
		verbose *bool
	}
	tests := []struct {
		name string
		args args
		want notations.Params
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreatePrm(tt.args.l, tt.args.verbose); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreatePrm() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateRes(t *testing.T) {
	type args struct {
		l string
	}
	tests := []struct {
		name string
		args args
		want notations.Response
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateRes(tt.args.l); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateRes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestImp(t *testing.T) {
	type args struct {
		fls []string
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
			if got := Imp(tt.args.fls); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Imp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsAbstrc(t *testing.T) {
	type args struct {
		fls []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAbstrc(tt.args.fls); got != tt.want {
				t.Errorf("IsAbstrc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsArray(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsArray(tt.args.line); got != tt.want {
				t.Errorf("IsArray() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsArrayType(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsArrayType(tt.args.line); got != tt.want {
				t.Errorf("IsArrayType() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsItrf(t *testing.T) {
	type args struct {
		fls []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsItrf(tt.args.fls); got != tt.want {
				t.Errorf("IsItrf() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJavaVariableHandler1(t *testing.T) {
	type args struct {
		line      string
		index     int
		wholeFile []string
	}
	tests := []struct {
		name    string
		args    args
		want    notations.Variable
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := JavaVariableHandler(tt.args.line, tt.args.index, tt.args.wholeFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("JavaVariableHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("JavaVariableHandler() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPackgName(t *testing.T) {
	type args struct {
		fls []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := PackgName(tt.args.fls); got != tt.want {
				t.Errorf("PackgName() = %v, want %v", got, tt.want)
			}
		})
	}
}
