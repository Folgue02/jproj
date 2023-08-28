package utils

import (
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const (
    GrepFiles uint = 0
    GrepDir
    GrepAll
)

func GrepFilesByExtension(targetPath, extension string, mode uint) ([]string, error) {
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
