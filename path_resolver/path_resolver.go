package path_resolver

import (
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _  = runtime.Caller(0)
	ROOT = filepath.Join(filepath.Dir(b), "../")
)