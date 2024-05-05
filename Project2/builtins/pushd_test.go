package builtins_test

import (
	"container/list"
	"errors"
	"testing"

	"github.com/ColtonDno/CSCE4600/Project2/builtins"
)

func TestSetAlias(t *testing.T) {
	dirs := list.New()
	tmp := t.TempDir()

	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name: "dir arg should change to that dir",
			args: args{
				args: []string{tmp},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// testing
			if err := builtins.PushDirectory(dirs, tt.args.args...); tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("ChangeDirectory() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			} else if err != nil {
				t.Fatalf("ChangeDirectory() unexpected error: %v", err)
			}

		})
	}
}
