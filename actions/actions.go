package actions

import (
	"log"
    "github.com/folgue02/jproj/actions/build"
    "github.com/folgue02/jproj/actions/clean"
    "github.com/folgue02/jproj/actions/createproject"
    "github.com/folgue02/jproj/actions/deps"
    jar "github.com/folgue02/jproj/actions/jar"
    "github.com/folgue02/jproj/actions/listproject"
    "github.com/folgue02/jproj/actions/run"
    newElement "github.com/folgue02/jproj/actions/newElement"
)

type Action struct {
    ActionFunction func([]string) error
    HelpMsg string
}

var Actions *map[string]Action

// Initializes the global variable 'Commands'
func InitializeCommands() {
    Actions = &map[string]Action{
        "help": { ListActions, "Lists all possible actions" },
        "createproject": { createproject.CreateProject, "Creates a new project." },
        "listproject": { listproject.ListProject, "Displays the information about the project." },
        "build": { build.BuildProject, "Builds/Compiles the current project." },
        "clean": { clean.CleanProject, "Cleans/Removes the compiled sources of the project." },
        "new": { newElement.NewElement, "Adds/Creates a new element to the project." },
        "run": { run.RunProject, "Runs the current project." },
        "deps": { deps.ManageDependencies, "Manage dependencies." },
        "jar": { jar.CreateJar, "Creates a jar based on the project specified." },
        "version": { version, "Displays jproj's version." },
    }
}

// Executes the 
func ExecuteAction(actionName string, args []string) bool {
    // In case they are not defined.
    if Actions == nil {
        InitializeCommands()
    }

    for cn, cmd := range *Actions {
        if cn == actionName {
            //log.Printf("==> Executing action with name '%s'...\n", actionName)
            if err := cmd.ActionFunction(args); err != nil {
                log.Printf("Error: Error while executing action '%s': %v\n", actionName, err)
            }
            return true
        }
    }
    return false
}
