package actions

import (
	"fmt"
	"log"
	"os"
	"text/tabwriter"
)

func listactions(args []string) error {
    tw := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
    // fmt.Fprintln(tw, "Action\tDescription")
    // fmt.Fprintln(tw, "------\t-----------")

    if len(args) > 1 {
        count := 0
        for _, actionName := range args[1:] {
            for cn, cmd := range *Commands {
                if cn == actionName {
                    fmt.Fprintf(tw, "%s:\t%s\n", cn, cmd.HelpMsg)
                    count += 1
                }
            }
        }
        if count > 0 {
            tw.Flush()
        } else {
            log.Println("No actions found with the names you've specified.")
        }
    } else {
        for cn, cmd := range *Commands {
            fmt.Fprintf(tw, "%s:\t%s\n", cn, cmd.HelpMsg)
        }
        tw.Flush()

    }
    return nil
}
