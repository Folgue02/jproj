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
    "github.com/folgue02/jproj/actions/validate"
)

type Action struct {
    ActionFunction func([]string) error
    HelpMsg string
}

var Actions *map[string]Action

// Initializes the global variable 'Commands'
func InitializeCommands() {
    Actions = &map[string]Action{
        "help": { ListActionsAction, "Lists all possible actions" },
        "createproject": { createproject.CreateProjectActionHandler, "Creates a new project." },
        "listproject": { listproject.ListProjectActionHandler, "Displays the information about the project." },
        "build": { build.BuildActionHandler, "Builds/Compiles the current project." },
        "clean": { clean.CleanActionHandler, "Cleans/Removes the compiled sources of the project." },
        "new": { newElement.NewElementActionHandler, "Adds/Creates a new element to the project." },
        "run": { run.RunProjectActionHandler, "Runs the current project." },
        "deps": { deps.ManageDependenciesActionHandler, "Manage dependencies." },
        "jar": { jar.CreateJarActionHandler, "Creates a jar based on the project specified." },
        "version": { versionAction, "Displays jproj's version." },
        "validate": { validate.ValidateActionHandler, "Validates the environment." },
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
                log.Printf("Error: Error while executing action '%s':\n", actionName)
                log.Printf("   %v\n", err)
            }
            return true
        }
    }
    return false
}
