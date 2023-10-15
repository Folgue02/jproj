package java

import "strings"

// Represents the configuration required for
// executing the command 'java'
type JavaCommand struct {
    JavaPath      string
    ClassPaths    []string
    MainClass     string
    ExecArguments []string
}

func NewJavaCommand(
    javaPath string,
    classPaths []string,
    mainClass string,
    args []string) JavaCommand {
    return JavaCommand{
        JavaPath:      javaPath,
        ClassPaths:    classPaths,
        MainClass:     mainClass,
        ExecArguments: args,
    }
}

// Returns the arguments for the 'java' command
func (j JavaCommand) Arguments() []string {
    args := make([]string, 0)

    if len(j.ClassPaths) > 0 {
        args = append(args, "-cp", strings.Join(j.ClassPaths, ":"))
    }

    args = append(args, j.MainClass)

    if len(j.ExecArguments) > 0 {
        args = append(args, j.ExecArguments...)
    }
    return args
}
