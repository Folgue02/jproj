package actionmanager

import (
	"bytes"
	"fmt"
	"text/tabwriter"
)

// Handler of the action, also known as the entry point.
type ActionHandler func([]string) error

// Represents an action, an executable ActionFunction
// that takes a collection of arguments and returns an error 
// and a descriptive message.
type Action struct {
    ActionHandler
    HelpMsg string
}

// Represents a collection of actions mapped by 
// their name as action.
type ActionCollection map[string]Action

// If there is an action with the name specified in the arguments, it will be executed,
// returning whatever gets returned from the action's handler with a 'true' value. 
// If there is no function with such name, 'nil, false' will be returned.
//
// NOTE: If the actionName is 'help' or 'listActions', then a table of the actions 
// and their helpMsgs will be printed, and 'nil, true' gets returned.
func (a ActionCollection) ExecuteAction(actionName string, args []string) (error, bool) {
    if actionName == "help" || actionName == "listactions" {
        fmt.Println(a.listActions())
        return nil, true
    }

    action, ok := a[actionName]

    if !ok {
        return nil, false
    }

    return action.ActionHandler(args), true
}

// Formats a string to contain a table of the action names pointing at their
// help message.
func (a ActionCollection) listActions() string {
    sb := bytes.NewBufferString("")
    tw := tabwriter.NewWriter(sb, 0, 0, 3, ' ', 0)

    for cn, cmd := range a {
        fmt.Fprintf(tw, "%s:\t%s\n", cn, cmd.HelpMsg)
    }
    fmt.Fprintf(tw, "help:\tList all available actions.")
    tw.Flush()

    return sb.String()
}
