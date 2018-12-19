package functions

import (
	"fmt"
	"io/ioutil"
	"kubor/template"
	"os"
	"path/filepath"
)

var _ = Register(Function{
	Name:     "pathExists",
	Category: "path",
	Parameters: Parameters{{
		Name: "path",
	}},
	Returns: Return{
		Description: "Is <true> if provided <path> does exist.",
	},
	Func: func(context template.ExecutionContext, path string) (bool, error) {
		return statCheckOfContext(context, path, nil)
	},
}, Function{
	Name:     "isFile",
	Category: "path",
	Parameters: Parameters{{
		Name: "path",
	}},
	Returns: Return{
		Description: "<true> if provided <path> does exist and is a file.",
	},
	Func: func(context template.ExecutionContext, path string) (bool, error) {
		return statCheckOfContext(context, path, func(fi os.FileInfo) bool {
			return !fi.IsDir()
		})
	},
}, Function{
	Name:     "isDir",
	Category: "path",
	Parameters: Parameters{{
		Name: "path",
	}},
	Returns: Return{
		Description: "<true> if provided <path> does exist and is a directory.",
	},
	Func: func(context template.ExecutionContext, path string) (bool, error) {
		return statCheckOfContext(context, path, func(fi os.FileInfo) bool {
			return fi.IsDir()
		})
	},
}, Function{
	Name:     "fileSize",
	Category: "path",
	Parameters: Parameters{{
		Name: "path",
	}},
	Returns: Return{
		Description: "The size of the provided <path>. It will fail with an error if the specified file does not exist or is not a file.",
	},
	Func: func(context template.ExecutionContext, path string) (int64, error) {
		if fi, err := statOfContext(context, path); err != nil {
			return 0, err
		} else {
			return fi.Size(), nil
		}
	},
}, Function{
	Name:     "readFile",
	Category: "path",
	Parameters: Parameters{{
		Name: "file",
	}},
	Returns: Return{
		Description: "The content of the provided <file>.",
	},
	Func: func(context template.ExecutionContext, file string) (string, error) {
		if resolved, err := resolvePathOfContext(context, file); err != nil {
			return "", fmt.Errorf("cannot resolve path of '%s': %v", file, err)
		} else if b, err := ioutil.ReadFile(resolved); err != nil {
			return "", fmt.Errorf("cannot read path '%s' (source:%s ): %v", resolved, file, err)
		} else {
			return string(b), nil
		}
	},
})

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
