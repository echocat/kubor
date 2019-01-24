package functions

import (
	"fmt"
	"github.com/levertonai/kubor/template"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
)

var FuncPathExists = Function{
	Parameters: Parameters{{
		Name: "path",
	}},
	Returns: Return{
		Description: "Is <true> if provided <path> does exist.",
	},
}.MustWithFunc(func(context template.ExecutionContext, path string) (bool, error) {
	return statCheckOfContext(context, path, nil)
})

var FuncIsFile = Function{
	Parameters: Parameters{{
		Name: "path",
	}},
	Returns: Return{
		Description: "<true> if provided <path> does exist and is a file.",
	},
}.MustWithFunc(func(context template.ExecutionContext, path string) (bool, error) {
	return statCheckOfContext(context, path, func(fi os.FileInfo) bool {
		return !fi.IsDir()
	})
})

var FuncIsDir = Function{
	Parameters: Parameters{{
		Name: "path",
	}},
	Returns: Return{
		Description: "<true> if provided <path> does exist and is a directory.",
	},
}.MustWithFunc(func(context template.ExecutionContext, path string) (bool, error) {
	return statCheckOfContext(context, path, func(fi os.FileInfo) bool {
		return fi.IsDir()
	})
})

var FuncFileSize = Function{
	Parameters: Parameters{{
		Name: "path",
	}},
	Returns: Return{
		Description: "The size of the provided <path>. It will fail with an error if the specified file does not exist or is not a file.",
	},
}.MustWithFunc(func(context template.ExecutionContext, path string) (int64, error) {
	if fi, err := statOfContext(context, path); err != nil {
		return 0, err
	} else {
		return fi.Size(), nil
	}
})

var FuncPathExt = Function{
	Parameters: Parameters{{
		Name: "path",
	}},
	Returns: Return{
		Description: "The file extension.",
	},
}.MustWithFunc(func(p string) string {
	return path.Ext(p)
})

var FuncPathBase = Function{
	Parameters: Parameters{{
		Name: "path",
	}},
	Returns: Return{
		Description: "The last element of path.",
	},
}.MustWithFunc(func(p string) string {
	return path.Base(p)
})

var FuncPathDir = Function{
	Parameters: Parameters{{
		Name: "path",
	}},
	Returns: Return{
		Description: "All but the last element of path, typically the path's directory.",
	},
}.MustWithFunc(func(p string) string {
	return path.Dir(p)
})

var FuncPathClean = Function{
	Parameters: Parameters{{
		Name: "path",
	}},
	Returns: Return{
		Description: "Shortest path name equivalent to path by purely lexical processing.",
	},
}.MustWithFunc(func(p string) string {
	return path.Clean(p)
})

var FuncReadFile = Function{
	Parameters: Parameters{{
		Name: "file",
	}},
	Returns: Return{
		Description: "The content of the provided <file>.",
	},
}.MustWithFunc(func(context template.ExecutionContext, file string) (string, error) {
	if resolved, err := resolvePathOfContext(context, file); err != nil {
		return "", fmt.Errorf("cannot resolve path of '%s': %v", file, err)
	} else if b, err := ioutil.ReadFile(resolved); err != nil {
		return "", fmt.Errorf("cannot read path '%s' (source:%s ): %v", resolved, file, err)
	} else {
		return string(b), nil
	}
})

var FuncsPath = Functions{
	"readFile":   FuncReadFile,
	"fileSize":   FuncFileSize,
	"pathExists": FuncPathExists,
	"isFile":     FuncIsFile,
	"isDir":      FuncIsDir,
	"pathExt":    FuncPathExt,
	"pathBase":   FuncPathBase,
	"pathDir":    FuncPathDir,
	"pathClean":  FuncPathClean,
}
var CategoryPath = Category{
	Functions: FuncsPath,
}

func resolvePathOfContext(context template.ExecutionContext, path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	var dir string
	sourceFile := context.GetTemplate().GetSourceFile()
	if sourceFile != nil {
		dir = filepath.Dir(*sourceFile)
	} else {
		if cwd, err := os.Getwd(); err != nil {
			return "", err
		} else {
			dir = cwd
		}
	}
	cleaned := filepath.Clean(dir + path)
	return filepath.Abs(cleaned)
}

func statOfContext(context template.ExecutionContext, path string) (os.FileInfo, error) {
	if resolved, err := resolvePathOfContext(context, path); err != nil {
		return nil, err
	} else {
		return os.Stat(resolved)
	}
}

func statCheckOfContext(context template.ExecutionContext, path string, predicate func(fi os.FileInfo) bool) (bool, error) {
	if fi, err := statOfContext(context, path); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	} else if predicate != nil {
		return predicate(fi), nil
	} else {
		return true, nil
	}
}
