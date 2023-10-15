package java

import (
    "strings"

    "github.com/folgue02/jproj/utils"
)

// Represents the configuration required for
// executing the command 'javac'
type JavacCommand struct {
    // Path/Command to the compiler
    CompilerPath string

    // Sources to be compiled
    Sources []string

    // Class paths to be included
    ClassPaths []string

    // Output dir
    TargetDir string
}

func NewJavacCommand(cc string, sources []string, classPaths []string, targetDir string) JavacCommand {
    return JavacCommand{
        CompilerPath: cc,
        Sources:      sources,
        ClassPaths:   classPaths,
        TargetDir:    targetDir,
    }
}

// Returns the arguments for the 'javac' command
func (j JavacCommand) Arguments() []string {
    cmd := make([]string, 0)
    cmd = append(cmd, j.Sources...)

    cmd = append(cmd, "-d", j.TargetDir)

    if len(j.ClassPaths) > 0 {
        cmd = append(cmd,
            "-cp",
            strings.Join(j.ClassPaths, ":"))
    }

    return cmd
}

// Adds the specified sources to the .Sources attribute.
func (j *JavacCommand) AddSources(newSources ...string) {
    j.Sources = append(j.Sources, newSources...)
}

// Adds all files ended in '.java' in the specified
// directory (recursively) to the sources.
func (j *JavacCommand) AddDirAsSource(dirPath string) error {
    sources, err := utils.GrepFilesByExtension(dirPath, "java", utils.GrepFiles)

    if err != nil {
        return err
    }

    j.Sources = append(j.Sources, sources...)

    return nil
}
