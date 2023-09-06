package split

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/tenntenn/golden"
)

var (
	flagUpdate bool
)

func init() {
	flag.BoolVar(&flagUpdate, "update", false, "update golden files")
}

func testTarget(dir string) error {
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("こんにちは"), 0700); err != nil {
		return err
	}

	if err := os.WriteFile(filepath.Join(dir, "b.txt"), []byte("world"), 0700); err != nil {
		return err
	}

	return nil
}

func Test(t *testing.T) {
	dir := t.TempDir()
	if err := testTarget(dir); err != nil {
		t.Fatal("unexpected error:", err)
	}

	got := golden.Txtar(t, dir)
	if diff := golden.Check(t, flagUpdate, "testdata", "mytest", got); diff != "" {
		t.Error(diff)
	}
}
