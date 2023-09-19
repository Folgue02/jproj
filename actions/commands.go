package actions

import (
	"log"
)

type Command struct {
    ActionFunction func([]string) error
    HelpMsg string
}

var Commands *map[string]Command

func InitializeCommands() {
    Commands = &map[string]Command {
        "help": { listactions, "Lists all possible actions" },
        "createproject": { createproject, "Creates a new project." },
        "listproject": { listproject, "Displays the information about the project." },
        "build": { buildProject, "Builds/Compiles the current project." },
        "clean": { clean, "Cleans/Removes the compiled sources of the project." },
        "new": { newElement, "Adds/Creates a new element to the project." },
        "run": { runProject, "Runs the current project." },
        "deps": { manageDependencies, "Manage dependencies." },
        "jar": { CreateJar, "Creates a jar based on the project specified." },
        "version": { version, "Displays jproj's version." },
    }
}

func ExecuteAction(actionName string, args []string) bool {
    // In case they are not defined.
    if Commands == nil {
        InitializeCommands()
    }

    for cn, cmd := range *Commands {
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
