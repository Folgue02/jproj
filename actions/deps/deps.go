package deps

import (
	"fmt"

	"github.com/folgue02/jproj/utils/actionmanager"
)

func DepMgrActionHandler(args []string) error {
    var actionName string
    var actionArgs []string
    if len(args) == 1 {
        return fmt.Errorf("No action specified.")
    } else if len(args) == 2 {
        actionName = args[1]
        actionArgs = args[1:] 
    } else {
        actionName = args[1]
        actionArgs = args[1:]
    }

    depActions := actionmanager.ActionCollection {
        "add": { 
            ActionHandler: AddActionHandler, 
            HelpMsg: "Adds a new dependency to the configuration",
        },
        "clean": { 
            ActionHandler: CleanActionHandler, 
            HelpMsg: "Removes the jar dependencies stored in the project's lib.",
        },
        "fetch": { 
            ActionHandler: FetchActionHandler, 
            HelpMsg: "Fetches the dependencies specified in the project's configuration." ,
        },
    }

    err, ok := depActions.ExecuteAction(actionName, actionArgs)

    if !ok {
        return fmt.Errorf("No dependency management action found with name '%s'.", actionName)
    }

    return err
}
