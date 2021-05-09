package notations

import "testing"

func TestCommentEndTag(t *testing.T) {
	type args struct {
		line string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test_is_endtag",
			args: args{line: "*/"},
			want: true,
		},
		{
			name: "test_is_something",
			args: args{line: "//"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CommentEndTag(tt.args.line); got != tt.want {
				t.Errorf("CommentEndTag() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHasIgnoreNotation(t *testing.T) {
	type args struct {
		index     int
		wholeFile []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test",
			args: args{
				index:     4,
				wholeFile: []string{"package com.something;", "import oracle.sucks.foo;", "import oracle.sucks.bar;", "public abstract class Alpha {"},
			},
			want: false,
		},
		{
			name: "test_happy_path",
			args: args{
				index:     4,
				wholeFile: []string{"package com.something;", "package com.something;", "import oracle.sucks.foo;", "@DD:IGNORE", "import oracle.sucks.bar;"},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := HasIgnoreNotation(tt.args.index, tt.args.wholeFile); got != tt.want {
				t.Errorf("HasIgnoreNotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsConnectionMethodNotation(t *testing.T) {
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
			if got := IsConnectionMethodNotation(tt.args.line); got != tt.want {
				t.Errorf("IsConnectionMethodNotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsController(t *testing.T) {
	type args struct {
		line []byte
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
			if got := IsController(tt.args.line); got != tt.want {
				t.Errorf("IsController() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsDescriptionNotation(t *testing.T) {
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
			if got := IsDescriptionNotation(tt.args.line); got != tt.want {
				t.Errorf("IsDescriptionNotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsEp(t *testing.T) {
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
			if got := IsEp(tt.args.line); got != tt.want {
				t.Errorf("IsEp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsIgnoreNotation(t *testing.T) {
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
			if got := IsIgnoreNotation(tt.args.line); got != tt.want {
				t.Errorf("IsIgnoreNotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsNotNullNotation(t *testing.T) {
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
			if got := IsNotNullNotation(tt.args.line); got != tt.want {
				t.Errorf("IsNotNullNotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsParamNotation(t *testing.T) {
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
			if got := IsParamNotation(tt.args.line); got != tt.want {
				t.Errorf("IsParamNotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsPayloadNotation(t *testing.T) {
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
			if got := IsPayloadNotation(tt.args.line); got != tt.want {
				t.Errorf("IsPayloadNotation() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsResponseNotation(t *testing.T) {
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
			if got := IsResponseNotation(tt.args.line); got != tt.want {
				t.Errorf("IsResponseNotation() = %v, want %v", got, tt.want)
			}
		})
	}
}
