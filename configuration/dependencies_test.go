package configuration

import (
	"testing"
)

func TestMavenURL(t *testing.T) {
	expected := "https://repo1.maven.org/maven2/org/junit/jupiter/junit-jupiter-api/5.10.0/junit-jupiter-api-5.10.0.jar"
	result := NewDependency("org.junit.jupiter", "junit-jupiter-api", "5.10.0").BuildMavenURLDownload()

	if result != expected {
		t.Errorf("(expected: %s, result: %s)", expected, result)
	}
}

func TestDepStringParsing(t *testing.T) {
	depString := "org.junit.jupiter:junit-jupiter-api:5.10.0"
	expected := NewDependency("org.junit.jupiter", "junit-jupiter-api", "5.10.0")
	result, err := NewDependencyFromString(depString)

	if err != nil {
		t.Errorf("NewDependencyFromString() -> _, %v", err)
	}

	if expected != *result {
		t.Errorf("(expected: %v, result: %v)", expected, result)
	}
}

func TestMissingColonDepStringParsing(t *testing.T) {
	depString := "org.junit.jupiterjunit-jupiter-api:5.10.0"
	result, err := NewDependencyFromString(depString)

	if err == nil {
		t.Errorf("No error, and this was returned: %v", result)
	}
}

func TestTooManyColonDepStringParsing(t *testing.T) {
	depString := "org.junit.jupiter:junit-jupiter-api:5.10.0:"

	result, err := NewDependencyFromString(depString)

	if err == nil {
		t.Errorf("No error, and this was returned: %v", result)
	}	
}