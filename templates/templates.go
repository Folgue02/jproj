package templates

import (
	"bytes"
	"text/template"
)

// TODO: Change this to use text/templates

const (
	ClassTemplate = `package {{ .Package }};

public class {{ .Name }} {

}`
	InterfaceTemplate = `package {{ .Package }};

public interface {{ .Name }} {
    
}`
	EnumTemplate = `package {{ .Package }};

public enum {{ .Name }} {

}`

	NoPackageClassTemplate = `
public class {{ .Name }} {

}`
	NoPackageInterfaceTemplate = `
public interface {{ .Name }} {

}`
	NoPackageEnumTemplate = `
public enum {{ .Name }} {

}`
)

type NewElementInfo struct {
    Package string
    Name string
}

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
    result := bytes.NewBufferString("")

    templ := template.New("Data")

    if elementPackage == "" {
        switch elementType {
        case ElementTypeClass:
            templ.Parse(NoPackageClassTemplate)
        case ElementTypeInterface:
            templ.Parse(NoPackageInterfaceTemplate)
        case ElementTypeEnum:
            templ.Parse(NoPackageEnumTemplate)
        default:
            return nil, false
        }
    } else {
        switch elementType {
        case ElementTypeClass:
            templ.Parse(ClassTemplate)
        case ElementTypeInterface:
            templ.Parse(InterfaceTemplate)
        case ElementTypeEnum:
            templ.Parse(EnumTemplate)
        default:
            return nil, false
        }
    }
    templ.Execute(result, NewElementInfo { Name: elementName, Package: elementPackage })
    resultString := result.String()
	return &resultString, true
}
