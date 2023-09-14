package configuration

import (
	"fmt"
	"strconv"
	"strings"
)

type JavaManifest map[string]any

// Turns the java manifest object into a real
// valid manifest string.
func (j JavaManifest) ToManifest() (*string, error) {
    s := strings.Builder {}

    for headerName, headerValue := range j {
        var headerValueString string
        switch headerValue.(type) {
        case int:
            value, _ := headerValue.(int)
            headerValueString = strconv.Itoa(value)

        case string:
            value, _ := headerValue.(string)
            headerValueString = value

        case []string:
            value, _ := headerValue.([]string)
            sb := strings.Builder {}

            for _, p := range value {
                sb.WriteString(fmt.Sprintf("\"%s\" ", p))
            }
            headerValueString = sb.String()

        case float64:
            value, _ := headerValue.(float64)
            headerValueString = fmt.Sprintf("%g", value)

        default:
            return nil, fmt.Errorf("Wrong type found in manifest object: %v", headerValue)
        }
        // TODO: Continue here
        s.WriteString(fmt.Sprintf("%s: %s", headerName, headerValueString))
        s.WriteRune('\n')
    }

    resultString := s.String()
    return &resultString, nil;
}
