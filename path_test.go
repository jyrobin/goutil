// Copyright (c) 2021 Jing-Ying Chen. Subject to the MIT License.

package goutil

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestPath(t *testing.T) {
	file, err := ioutil.TempFile("", "prefix.*.txt")
	if err != nil {
		t.Fatal(err)
	}

	fpath := file.Name()
	dir, fname := filepath.Dir(fpath), filepath.Base(fpath)
	defer os.Remove(fpath)

	fmt.Printf("%s %s\n", dir, fname)

	if !FileExists(fpath, "") {
		t.Fatal("should exists")
	}
	if !FileExists(dir, fname) {
		t.Fatal("should exists")
	}
	if !DirExists(dir, "") {
		t.Fatal("should exists")
	}
	if FileExists(dir, "") {
		t.Fatal("should not exists")
	}

	if str, _ := ResolvePath("xxx", fpath); str != fpath {
		t.Fatalf("%s != %s", str, fpath)
	}
	if str, _ := ResolvePath(dir, fname); str != fpath {
		t.Fatalf("%s != %s", str, fpath)
	}

	if finfo, str, _ := GetExistingFileInfo(dir, fname); finfo == nil || str != fpath {
		t.Fatalf("nil file info or bad path")
	}

	if str, _ := ResolveExistingFile(dir, fname); str != fpath {
		t.Fatalf("%s != %s", str, fpath)
	}
	if str, _ := ResolveExistingDir(dir, ""); str != dir {
		t.Fatalf("%s != %s", str, dir)
	}
}
