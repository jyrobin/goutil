// Copyright (c) 2021 Jing-Ying Chen. Subject to the MIT License.

package goutil

import (
	"fmt"
	"os"
	"path/filepath"
)

func joinPair(basePath, relPath string) string {
	if relPath == "" {
		return basePath
	}
	if basePath == "" || filepath.IsAbs(relPath) {
		return relPath
	}

	return filepath.Join(basePath, relPath)
}

// return fpath even if bad basePath
func ResolvePath(basePath string, relPaths ...string) (string, error) {
	for _, relPath := range relPaths {
		basePath = joinPair(basePath, relPath)
	}

	abs, err := filepath.Abs(basePath)
	if err != nil {
		return basePath, err
	}

	return abs, nil
}

// return fpath even not exists
func GetExistingFileInfo(basePath string, relPaths ...string) (os.FileInfo, string, error) {
	var fi os.FileInfo = nil
	fpath, err := ResolvePath(basePath, relPaths...)
	if err == nil {
		fi, err = os.Stat(fpath)
	}
	return fi, fpath, err
}

func FileExists(basePath string, relPaths ...string) bool {
	fi, _, _ := GetExistingFileInfo(basePath, relPaths...)
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
		err = fmt.Errorf("%s %s is not a folder", basePath, relPath)
	}
	return fpath, err
}
