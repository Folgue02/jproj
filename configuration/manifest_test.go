package configuration

import "testing"

func TestGenerateManifest(t *testing.T) {
    expected := `Created-By: User
Main-Class: me.user.app.App
Class-Path: "./target" "./lib" 
Version: 2.3
`
    input := Configuration {
        Manifest: JavaManifest { 
            "Created-By": "User", 
            "Main-Class": "me.user.app.App",
            "Class-Path": []string { 
                "./target",
                "./lib",
            },
            "Version": 2.3,
        },
    }

    result, err := input.Manifest.ToManifest()

    if err != nil {
        t.Errorf("Error: %v", err)
    }

    if *result != expected {
        t.Errorf("Different (\n---Expected:\n%s\n---Result:\n %s", expected, *result)
    }
}
