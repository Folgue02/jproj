package utils

import (
	"path"
	"strings"
)

var javaKwds = []string {
    "import",
    "native",
    "package",
    "class",
    "return",
    "enum",
    "interface",
    "public",
    "static",
    "void",
    "int",
    "float",
    "short",
    "long",
    "char",
    "private",
    "protected",
    "final",
    "const",
    "goto",
    "abstract",
}

type JavaPackagePath struct {
    Original string
}

func NewJavaPackagePath(original string) JavaPackagePath {
    return JavaPackagePath { original }
}

func NewJavaPackagePathFromSlice(originalSlice []string) JavaPackagePath {
    original := strings.Join(originalSlice, ".")
    return NewJavaPackagePath(original)
}

func (j *JavaPackagePath) DividePackageName() []string {
    return strings.Split(j.Original, ".") // TODO: Care for empty spaces
}

func (j *JavaPackagePath) ToJavaFilePath() string {
    dividedPackageName := DividePackageName(j.Original)
    dividedPackageName[len(dividedPackageName)-1] += ".java"
    return path.Join(dividedPackageName...)
}

func (j *JavaPackagePath) ToJavaDirPath() string {
    return path.Join(j.DividePackageName()...)
}

func (j *JavaPackagePath) Base() string {
    divided := j.DividePackageName()

    if len(divided) == 0 {
        return ""
    } else {
        return divided[len(divided)-1]
    }
}

func (j *JavaPackagePath) Dir() string {
    divided := j.DividePackageName()

    if len(divided) == 0 {
        return ""
    } else {
        jpp := NewJavaPackagePathFromSlice(divided[:len(divided)-1])
        return jpp.ToJavaDirPath()
    }
}

func (j *JavaPackagePath) DirPackagePath() string {
    return strings.Replace(j.Dir(), "/", ".", -1)
}

func validatePackageName(pkgName string) bool {
    for _, javaKwd := range javaKwds {
        if strings.Contains(pkgName, javaKwd) {
            return false
        }
    }

    if strings.Contains(pkgName, "..") || strings.HasPrefix(pkgName, ".") ||
        strings.HasSuffix(pkgName, ".") {
        return false
    }

    /* TODO: Finish this
    for i, c := range pkgName {
    }*/
    return true
}

// Splits the packagePath into a slice containing all the
// segments.
func DividePackageName(packagePath string) []string {
    return strings.Split(packagePath, ".") // TODO: Care for empty spaces
}

// Turns packagePath (i.e. com.company.app.App) into 
// its file path version (i.e. com/company/app/App.java)
func JavaPathToFilePath(packagePath string) string {
    dividedPackageName := DividePackageName(packagePath)
    dividedPackageName[len(dividedPackageName)-1] += ".java"
    return path.Join(dividedPackageName...)
}

// Turns packagePath (i.e. com.company.app) into
// its dir path version (i.e. com/company/app)
func JavaPathToDirPath(packagePath string) string {
    return path.Join(DividePackageName(packagePath)...)
}

