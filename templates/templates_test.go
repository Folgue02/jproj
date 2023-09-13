package templates

import "testing"

func TestTemplates(t *testing.T) {
    expected := `package me.user.App;

public class Something {

}`
    output, happened := GenerateTemplate(ElementTypeClass, "Something", "me.user.App")

    if !happened {
        t.Errorf("Didn't do template.")
    }

    if *output != expected {
        t.Errorf("Templates differ:\n Expected: %s\n Got: %s", expected, *output)
    }

}
