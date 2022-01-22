// Copyright (c) 2021 Jing-Ying Chen. Subject to the MIT License.

package goutil

import (
	"fmt"
	"os"
	"path/filepath"
)

func ResolvePath(basePath, relPath string) (string, error) {
	if filepath.IsAbs(relPath) {
		return relPath, nil
	}

	abs, err := filepath.Abs(basePath)
	if err != nil {
		return filepath.Join(basePath, relPath), err
	}

	return filepath.Join(abs, relPath), nil
}

func GetExistingFileInfo(basePath, relPath string) (os.FileInfo, string, error) {
	var fi os.FileInfo
	fpath, err := ResolvePath(basePath, relPath)
	if err == nil {
		if fi, err = os.Stat(fpath); err == nil {
			return fi, fpath, nil
		}
	}

	return nil, fpath, err
}

func FileExists(basePath, relPath string) bool {
	fi, _, _ := GetExistingFileInfo(basePath, relPath)
	return fi != nil && fi.Mode().IsRegular()
}

func DirExists(basePath, relPath string) bool {
	fi, _, _ := GetExistingFileInfo(basePath, relPath)
	return fi != nil && fi.Mode().IsDir()
}

func ResolveExistingFile(basePath, relPath string) (string, error) {
	fi, fpath, err := GetExistingFileInfo(basePath, relPath)
	if fi == nil {
		return fpath, err
	}
	if !fi.Mode().IsRegular() {
		err = fmt.Errorf("%s %s is not a file", basePath, relPath)
	}
	return fpath, err
}

func ResolveExistingDir(basePath, relPath string) (string, error) {
	fi, fpath, err := GetExistingFileInfo(basePath, relPath)
	if fi == nil {
		return fpath, err
	}
	if !fi.Mode().IsDir() {
		err = fmt.Errorf("%s %s is not a file", basePath, relPath)
	}
	return fpath, err
}
