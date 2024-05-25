package path_resolver

import (
	"path/filepath"
	"runtime"
)

var root string

func GetRoot() string {
	_, b, _, _  := runtime.Caller(0)
	root = filepath.Join(filepath.Dir(b), "../")
	return root
}