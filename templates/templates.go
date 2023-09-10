package templates

import (
	"fmt"
	"strings"
)

// TODO: Change this to use text/templates

const (
	ClassTemplate = `
package <Package>;

public class <Name> {
}`
	InterfaceTemplate = `
package <Package>;

public interface <Name> {
}`
	EnumTemplate = `
package <Package>;

public enum <Name> {
}`
)

const (
	ElementTypeClass     = "class"
	ElementTypeInterface = "interface"
	ElementTypeEnum      = "enum"
	ElementTypePackage   = "package"
)

// Returns a template if the element type matches "class", "enum" or "interface".
// If it doesn't, (nil, false) will be returned. NOTE: `elementName` doesn't get
// checked, so this element name can be invalid.
func GenerateTemplate(elementType, elementName, elementPackage string) (*string, bool) {
	var result string
	switch elementType {
	case ElementTypeClass:
		result = strings.Replace(
			strings.Replace(ClassTemplate, "<Name>", elementName, -1),
			"<Package>",
			elementPackage, -1)
	case ElementTypeInterface:
		result = strings.Replace(
			strings.Replace(InterfaceTemplate, "<Name>", elementName, -1),
			"<Package>",
			elementPackage, -1)
	case ElementTypeEnum:
		result = strings.Replace(
			strings.Replace(EnumTemplate, "<Name>", elementName, -1),
			"<Package>",
			elementPackage, -1)
	default:
		return nil, false
	}
	fmt.Printf("Generated template: %s\n", result)
	return &result, true
}
