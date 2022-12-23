package diff

import (
	"os"
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
			if diff := cmp.Diff(tt.args.want, tt.args.got); diff != "" {
				t.Errorf("%v", diff)
			}
		})
	}
}

func TestDiffRegister(t *testing.T) {
	result := DiffAll(AzurermRegisters())
	result.FixDocuments()
	t.Logf("%s", result.ToString())
}

func TestDiffAll(t *testing.T) {
	result := DiffAll(AzurermRegistersAll())
	result.FixDocuments()
	t.Logf("%s", result.ToString())

	md := result.CrossCheckIssues()
	t.Logf(md)

	if err := os.WriteFile("./diffcross.md", []byte(md), 0666); err != nil {
		t.Fatalf("write fail: %v", err)
	}
}
