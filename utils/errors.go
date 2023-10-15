package utils

import (
	"fmt"
	"strings"
)

func BeautifyError(err error, action func(string)) {
    sections := strings.Split(err.Error(), ":") 
    for i, section := range sections {
        if i < len(sections) - 1 && !strings.HasSuffix(section, ":") {
            action(fmt.Sprintf("%s%s:\n", strings.Repeat("  ", i), section))
        } else {
            action(fmt.Sprintf("%s%s\n", strings.Repeat("  ", i), section))
        }
    }
}
