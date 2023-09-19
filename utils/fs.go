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
func CleanDirectory(dirPath string) error {
    entries, err := os.ReadDir(dirPath)

    if err != nil {
        return err
    }

    for _, entry := range entries {
        err = os.RemoveAll(path.Join(dirPath, entry.Name()))

        if err != nil {
            return err
        }
    }

    return nil
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
