package configuration

import (
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
)

// https://repo1.maven.org/maven2/{groupID separated by slashes}/{name}/{version}/{name}-{version}.jar

type Dependency struct {
	GroupID string `json:"groupID"`
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Parses the string passed with the following format:
// {groupID}:{name}:{version} (i.e. org.junit.jupiter:junit-jupiter-api:5.10.0)
func NewDependencyFromString(depString string) (*Dependency, error) {
	splitDep := strings.Split(depString, ":")

	if len(splitDep) != 3 {
		return nil, fmt.Errorf("Wrong format dependency.")
	} else {
		dep := NewDependency(splitDep[0], splitDep[1], splitDep[2])
		return &dep, nil
	}
}

func NewDependency(groupID, name, version string) Dependency {
	return Dependency{
		groupID,
		name,
		version,
	}
}

// Returns the download URL for a dependency in the Maven repository.
func (d Dependency) BuildMavenURLDownload() string {
	slashedGroupID := path.Join(strings.Split(d.GroupID, ".")...)
	return fmt.Sprintf("https://repo1.maven.org/maven2/%s/%s/%s/%s-%s.jar", slashedGroupID, d.Name, d.Version, d.Name, d.Version)
}

func (d Dependency) GetJarName() string {
	return fmt.Sprintf("%s-%s.jar", d.Name, d.Version)
}

// Fetches the dependency, and writes its contents to 'w'.
func (d Dependency) fetchDependency(w io.Writer) error {
	response, err := http.Get(d.BuildMavenURLDownload())

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("The status code ('%d') differs from 200.", response.StatusCode)
	}

	if _, err = io.Copy(w, response.Body); err != nil {
		return err
	}

	return nil
}

func (d Dependency) String() string {
	return fmt.Sprintf("%s:%s:%s", d.GroupID, d.Name, d.Version)
}
