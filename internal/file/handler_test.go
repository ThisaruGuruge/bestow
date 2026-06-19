package file_test

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/ThisaruGuruge/bestow/internal/file"
)

func TestGetExistingFileType(t *testing.T) {
	tests := []struct {
		name      string
		setup     func(t *testing.T, src, dst string)
		wantType  file.ExistingType
		wantError bool
		err       string
	}{
		{
			name: "Managed Symlink",
			setup: func(t *testing.T, src, dst string) {
				if err := os.Symlink(src, dst); err != nil {
					t.Fatal(err)
				}
			},
			wantType: file.ExistingManagedSymlink,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()
			src := filepath.Join(dir, "src_file")
			dst := filepath.Join(dir, "dst_file")
			if err := os.WriteFile(src, []byte("src file"), 0o644); err != nil {
				t.Fatal(err)
			}
			if tc.setup != nil {
				tc.setup(t, src, dst)
			}
			h := file.NewHandler(slog.New(slog.NewTextHandler(io.Discard, nil)))
			actual, err := h.GetExistingFileType(src, dst)
			if err != nil && tc.wantError && tc.err != err.Error() {
				t.Fatalf("actual %s, expected %s", err.Error(), tc.err)
			}
			if (err != nil) != tc.wantError {
				t.Fatalf("err = %v, wantErr %v", err, tc.wantError)
			}
			if !tc.wantError && actual != tc.wantType {
				t.Fatalf("actual %q, expected %q", actual, tc.wantType)
			}
		})
	}
}
