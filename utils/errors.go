package utils

import (
	"fmt"
	"strings"
)

type ErrorBundle struct {
    Errors []error
}

func (e *ErrorBundle) Add(newError error) {
    e.Errors = append(e.Errors, newError)
}

func (e ErrorBundle) Error() string {
    sb := strings.Builder {}

    for _, err := range e.Errors {
        sb.WriteString(fmt.Sprintf("%s\n", err))
    }

    return sb.String()
}

func (e ErrorBundle) Len() int {
    return len(e.Errors)
}

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
