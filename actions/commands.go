package actions

import (
	"log"
)

type Command struct {
    ActionFunction func([]string)
    HelpMsg string
}

var Commands *map[string]Command

func InitializeCommands() {
    Commands = &map[string]Command {
        "help": Command { help, "Displays a help message" },
        "listactions": Command { listactions, "Lists all possible actions" },
        "createproject": Command { createproject, "Creates a new project." },
        "listproject": Command { listproject, "Displays the information about the project." },
        "build": Command { buildProject, "Builds/Compiles the current project." },
        "clean": Command { clean, "Cleans/Removes the compiled sources of the project." },
        "new": Command { newElement, "Adds/Creates a new element to the project." },
    }
}

func ExecuteAction(actionName string, args []string) bool {
    // In case they are not defined.
    if Commands == nil {
        InitializeCommands()
    }

    for cn, cmd := range *Commands {
        if cn == actionName {
            log.Printf("==> Executing action with name '%s'...\n", actionName)
            cmd.ActionFunction(args)
            return true
        }
    }
    return false
}
