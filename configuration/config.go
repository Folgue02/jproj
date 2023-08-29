package configuration

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
)

type Configuration struct {
    ProjectAuthor string `json:"author"`
    ProjectDescription string `json:"description"`
    ProjectName string `json:"project_name"`
    ProjectVersion string `json:"project_version"`
    ProjectTarget string `json:"project_target_path"`
    ProjectLib string `json:"project_lib_path"`
    MainClassPath string `json:"main_class_path"`
    Dependencies []Dependency `json:"dependencies"`
}

func NewConfiguration(projectName string) Configuration {
    return Configuration {
        ProjectName: projectName,
        ProjectAuthor: "anon",
        ProjectDescription: "A Java project.",
        ProjectVersion: "1.0",
        ProjectTarget: "./target",
        ProjectLib: "./lib",
        MainClassPath: "App",
        Dependencies: []Dependency {},
    }
}

// Based on the configuration, this method will create a project.
func (c Configuration) CreateProject(baseDirectory string) error {
    baseDirectory = path.Join(baseDirectory, c.ProjectName)

    configFilePath := path.Join(baseDirectory, "jproj.json")

    log.Printf("Creating directory %s...", baseDirectory)
    if err := os.MkdirAll(baseDirectory, 0750); err != nil {
        return err
    }


    for _, dir := range []string { c.ProjectTarget, c.ProjectLib, "./src" } {
        err := os.MkdirAll(path.Join(baseDirectory, dir), 0750)

        if err != nil {
            return err
        }
    }

    // TODO: Use this.SaveConfiguration()
    jsonConfiguration, _ := json.MarshalIndent(c, "", "    ")

    if err := os.WriteFile(configFilePath, jsonConfiguration, 0750); err != nil {
        return err
    }

    return nil
}

// Checks if there is a dependency with the same name as the one
// specified.
func (c *Configuration) DependencyExists(name string) bool {
    for _, dep := range c.Dependencies {
        if dep.Name == name {
            return true
        }
    }

    return false
}

func (c *Configuration) SaveConfiguration(destPath string) error {
    destStat, err := os.Stat(destPath)
    if err != nil {
        return fmt.Errorf("Cannot stat path: %v", err)
    }

    if destStat.IsDir() {
        destPath = path.Join(destPath, "jproj.json")

    }

    configContent, _ := json.MarshalIndent(c, "", "    ")

    err = os.WriteFile(destPath, configContent, 0750)

    if err != nil {
        return fmt.Errorf("Cannot write to config file due to the following error: %v", err)
    } else {
        return nil
    }
}

// Fetches the dependencies and stores them in the dependencies directory inside
// of the project. The 'projectDirectory' string must contain the string pointing
// to the project's location, (i.e., if the project is named 'newproject', the 
// projectDirectory variable should be '/home/user/source/newproject/').
//
// NOTE: This method overwrites the .jar files that are already stored in the project.
func (c Configuration) FetchDependencies(projectDirectory string) error {
    depsPath := path.Join(projectDirectory, c.ProjectLib)

    depsPathStat, err := os.Stat(depsPath)
    
    if err != nil {
        return fmt.Errorf("Cannot stat project's lib: %v", err)
    }

    if !depsPathStat.IsDir() {
        return fmt.Errorf("The project's lib ('%s') is not a directory.", depsPath)
    }

    for i, dep := range c.Dependencies {
        jarName, jarURL := dep.GetJarName(), dep.BuildMavenURLDownload()
        log.Printf("[%d]: Fetching dependency '%s' from '%s'...\n", i + 1, jarName, jarURL)
        jarFile, err := os.Create(path.Join(depsPath, jarName))

        if err != nil {
            return fmt.Errorf("Error while creating jar file ('%s'): %v", jarName, err)
        }

        if err := dep.fetchDependency(jarFile); err != nil {
            return fmt.Errorf("Error writing to jar file ('%s'): %v", jarName, err)
        }
    }

    log.Printf("Dependencies fetched.\n")
    return nil
}

func LoadConfigurationFromString(rawString string) (*Configuration, error) {
    var config Configuration
    err := json.Unmarshal([]byte(rawString), &config)

    if err!= nil {
        return nil, err 
    }

    return &config, err
}

func (c Configuration) String() string {
    return fmt.Sprintf(`
Name: %s
Target directory: %s`, c.ProjectName, c.ProjectTarget)
}

func LoadConfigurationFromFile(filePath string) (*Configuration, error) {
    filePathStat, err := os.Stat(filePath)

    if err != nil {
        return nil, err
    } else if filePathStat.IsDir() {
        filePath = path.Join(filePath, "jproj.json")
    }

    fileContent, err := os.ReadFile(filePath)

    if err != nil {
        return nil, err
    }

    return LoadConfigurationFromString(string(fileContent))
}
