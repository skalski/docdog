package spring

import (
	"docdog/src/notations"
	"reflect"
	"testing"
)

func TestCreateApiEndpoint(t *testing.T) {
	type args struct {
		index     int
		wholeFile []string
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
			if got := CreateApiEndpoint(tt.args.index, tt.args.wholeFile); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateApiEndpoint() = %v, want %v", got, tt.want)
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

func TestFetchMethodFromMapping(t *testing.T) {
	type args struct {
		s string
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
			if got := FetchMethodFromMapping(tt.args.s); got != tt.want {
				t.Errorf("FetchMethodFromMapping() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetProtocolFormMappingTag(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name       string
		args       args
		wantResult string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := GetProtocolFormMappingTag(tt.args.str); gotResult != tt.wantResult {
				t.Errorf("GetProtocolFormMappingTag() = %v, want %v", gotResult, tt.wantResult)
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotResult := GetStringFromQouteLine(tt.args.str); gotResult != tt.wantResult {
				t.Errorf("GetStringFromQouteLine() = %v, want %v", gotResult, tt.wantResult)
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

func TestIsFinal(t *testing.T) {
	type args struct {
		s string
		i int
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsFinal(tt.args.s, tt.args.i); got != tt.want {
				t.Errorf("IsFinal() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSBVariableHandler(t *testing.T) {
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
			got, err := SBVariableHandler(tt.args.line, tt.args.index, tt.args.wholeFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("SBVariableHandler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SBVariableHandler() got = %v, want %v", got, tt.want)
			}
		})
	}
}
