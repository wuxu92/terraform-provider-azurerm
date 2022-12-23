package check

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSliceDiff(t *testing.T) {
	type args struct {
		want []string
		got  []string
	}
	tests := []struct {
		name     string
		args     args
		wantDiff int
	}{
		{
			name: "aaa",
			args: args{
				want: []string{"abc", "def"},
				got:  []string{"def", "abc"},
			},
			wantDiff: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if missed, odd := SliceDiff(tt.args.want, tt.args.got, true); len(missed)+len(odd) != tt.wantDiff {
				t.Errorf("SliceDiff() missed: %v, odd: %v", missed, odd)
			}
			if diff := cmp.Diff(tt.args.want, tt.args.got); diff == "" {
				t.Errorf("%s should have diff in cmp.Diff", tt.name)
			}
		})
	}
}

func TestDiffRegister(t *testing.T) {
	result := DiffAll(AzurermRegisters())
	if err := result.FixDocuments(); err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", result.ToString())
}

func TestDiffAll(t *testing.T) {
	result := DiffAll(AzurermRegistersAll())
	if err := result.FixDocuments(); err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", result.ToString())
}
