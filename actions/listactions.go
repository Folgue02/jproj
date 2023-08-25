package actions

import (
    "fmt"
)

func listactions(args []string) {
    for cn, cmd := range *Commands {
        fmt.Printf("[%s]: %s\n", cn, cmd.HelpMsg)
    }
}
