package utils

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type GrepMode uint

const (
    GrepFiles GrepMode = 0
    GrepDir
    GrepAll
)

const (
    DefaultFilePermission = 0650
    DefaultDirPermission = 0750
)

func GrepFilesByExtension(targetPath, extension string, mode GrepMode) ([]string, error) {
    paths := make([]string, 0)

    err := filepath.Walk(targetPath, func(path string, fileInfo fs.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if strings.HasSuffix(path, "." + extension) {
            if mode == GrepFiles && !fileInfo.IsDir() ||
                (mode == GrepAll || mode == GrepDir) {
                paths = append(paths, path)
            } 
        }
        return nil
    })

    if err != nil {
        return nil, err
    }

    return paths, nil
}

// Removes all contents from the specified directory without erasing 
// the directory itself. In case something fails (dirPath doesn't exist, 
// cannot remove certain item...), the value will hold an error.
// 
// fileFilter is a function that gets executed with each entry's absolute 
// path as an argument. If fileFilter returns 'false', the entry won't be 
// removed.
func CleanDirectory(dirPath string, fileFilter func(string) bool) error {
    entries, err := os.ReadDir(dirPath)

    if err != nil {
        return err
    }

    for _, entry := range entries {
        entryPath := path.Join(dirPath, entry.Name())
        if (fileFilter(entryPath)) {
            err = os.RemoveAll(entryPath)

            if err != nil {
                return err
            }
        }
    }

    return nil
}

// Same as 'CleanDirectory' but uses an empty filter as fileFilter.
func CleanAllDirectory(dirPath string) error {
    return CleanDirectory(dirPath, func(s string) bool { return true })
}

// Checks if the directory specified contains a valid structure.
// 
// Checks performed:
//
// - ${dirPath} exists
//
// - ./jproj.json exists
//
// - ./src/ exists
func CheckValidityOfDirectory(dirPath string) bool {
    // TODO: Look for a better way to do this.
    dirPathStat, err := os.Stat(dirPath)
    configPathStat, cerr := os.Stat(path.Join(dirPath, "jproj.json"))
    srcPathStat, serr := os.Stat(path.Join(dirPath, "src"))

    if err != nil || !dirPathStat.IsDir() ||
        cerr != nil || configPathStat.IsDir() ||
        serr != nil || !srcPathStat.IsDir() {
        return false
    } else {
        return true
    }
}
