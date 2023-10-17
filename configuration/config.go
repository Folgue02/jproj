package configuration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"text/tabwriter"

	"github.com/folgue02/jproj/utils"
)

type Configuration struct {
	ProjectAuthor      string       `json:"author"`
	ProjectDescription string       `json:"description"`
	ProjectName        string       `json:"project_name"`
	ProjectVersion     string       `json:"project_version"`
	ProjectTarget      string       `json:"project_target_path"`
	ProjectLib         string       `json:"project_lib_path"`
	ProjectBin         string       `json:"project_bin_path"`
	Manifest           JavaManifest `json:"manifest,omitempty"`
	MainClassPath      string       `json:"main_class_path,omitempty"`
	Dependencies       []Dependency `json:"dependencies"`
    IncludedProjects   []string     `json:"included_projects"`
}

func NewConfiguration(projectName string) Configuration {
	return Configuration{
		ProjectName:        projectName,
		ProjectAuthor:      "anon",
		ProjectDescription: "A Java project.",
		ProjectVersion:     "1.0",
		ProjectTarget:      "./target/classes",
		ProjectLib:         "./lib",
		ProjectBin:         "./target/bin",
		MainClassPath:      "App",
		Dependencies:       []Dependency{},
        IncludedProjects:   []string{},
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

	for _, dir := range []string{c.ProjectTarget, c.ProjectLib, "./src", c.ProjectBin} {
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

// Checks the integrity of the project, checking the existence of directories.
// If anything is wrong with the project, this method will return an error with
// a string containing each problem. NOTE: This method doesn't check if the
// configuration file exists.
func (c Configuration) Validate(location string) error {
	paths := []string{
		c.ProjectBin,
		c.ProjectLib,
		"./src/",
	}

	errors := utils.ErrorBundle {} 

    // Check directories related to the current project
	for _, p := range paths {
		fullP := path.Join(location, p)

		pathStat, err := os.Stat(fullP)

		if err != nil {
			errors.Add(fmt.Errorf("Cannot stat project directory '%s': %s", fullP, err))
            continue
		}

		if !pathStat.IsDir() {
            errors.Add(fmt.Errorf("Path '%s' is not a directory", fullP))
            continue
		}
	}

	if errors.Len() > 0 {
        return errors
	}
	return nil
}

func (c Configuration) ValidateIncluded(location string) error {
    errors := utils.ErrorBundle {}

    for _, includedPath := range c.IncludedProjects {
        includedPath = filepath.Join(location, includedPath)

        includedConfig, err := LoadConfigurationFromFile(includedPath)

        if err != nil {
            errors.Add(fmt.Errorf("Cannot load project's ('%s') configuration: %v", includedPath, err))
            continue
        }

        err = includedConfig.Validate(includedPath)

        if err != nil {
            errors.Add(fmt.Errorf("Errors found while validating included project '%s': %v", includedPath, err))
        }
    }

    if errors.Len() == 0 {
        return nil
    }
    return errors
}

// Checks if there is a dependency with the same name as one
// specified inside of the .Dependencies attribute.
func (c *Configuration) DependencyExists(name string) bool {
	for _, dep := range c.Dependencies {
		if dep.Name == name {
			return true
		}
	}

	return false
}

// Returns the name of what would be the output jar from the project.
func (c Configuration) OutputJarName() string {
	return fmt.Sprintf("%s-%s.jar", c.ProjectName, c.ProjectVersion)
}

// Saves the configuration to the path specified in the case that the
// path specified points to a directory, "jproj.json" will be added at
// the end.
func (c Configuration) SaveConfiguration(destPath string) error {
	destStat, err := os.Stat(destPath)
	if err != nil {
		return fmt.Errorf("Cannot stat path: %v", err)
	}

	if destStat.IsDir() {
		destPath = path.Join(destPath, "jproj.json")

	}

	configContent, _ := json.MarshalIndent(c, "", "    ")

	err = os.WriteFile(destPath, configContent, utils.DefaultFilePermission)

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
		log.Printf("[%d]: Fetching dependency '%s' from '%s'...\n", i+1, jarName, jarURL)
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

// Returns a 'Configuration' object based on the string passed (this should be
// in json format representing the configuration)
func LoadConfigurationFromString(rawString string) (*Configuration, error) {
	var config Configuration
	err := json.Unmarshal([]byte(rawString), &config)

	if err != nil {
		return nil, err
	}

	return &config, err
}

// Checks if the configuration of the project contains a main class or
// not, if the .MainClassPath attribute is set to "", it would mean that
// the project is not meant to be executed.
func (c Configuration) IsExecutableProject() bool {
	return c.MainClassPath != ""
}

// Checks if the given name of the jar is associated with one
// of the dependencies in the project's configuration.
func (c Configuration) IsJarInDependencies(jarName string) bool {
	jarName = filepath.Base(jarName)
	for _, dep := range c.Dependencies {
		if dep.GetJarName() == jarName {
			return true
		}
	}
	return false
}

func (c Configuration) String() string {
	var result string
	resultWriter := bytes.NewBufferString(result)
	tw := tabwriter.NewWriter(resultWriter, 0, 0, 3, ' ', 0)
	fmt.Fprintf(tw, "Author name\t%s\n", c.ProjectAuthor)
	fmt.Fprintf(tw, "Description\t%s\n", c.ProjectDescription)
	fmt.Fprintf(tw, "Version\t%s\n", c.ProjectVersion)
	fmt.Fprintf(tw, "Target directory\t%s\n", c.ProjectTarget)
	fmt.Fprintf(tw, "Lib directory\t%s\n", c.ProjectLib)
	fmt.Fprintf(tw, "Bin directory\t%s\n", c.ProjectBin)
	if c.MainClassPath != "" {
		fmt.Fprintf(tw, "Main class\t%s\n", c.MainClassPath)
	}

	if len(c.Dependencies) > 0 {
		fmt.Fprintf(tw, "Dependencies\t%s\n", c.Dependencies[0])

		for _, dep := range c.Dependencies {
			fmt.Fprintf(tw, "\t%s\n", dep)
		}
	}

	tw.Flush()
	return resultWriter.String()
}

// Returns a 'Configuration' object created by reading the file specified.
//
// NOTE: If a directory is specified instead of a file, this function will
// attempt to read the {filePath}/jproj.json file.
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

// Returns a list of the .jar files inside of the .ProjectLib directory of
// the project.
func (c Configuration) ListJarInLib(filePath string) ([]string, error) {
	return utils.GrepFilesByExtension(path.Join(filePath, c.ProjectLib), "jar", utils.GrepFiles)
}

func (c Configuration) ListIncludedTargetPaths(configPath string) ([]string, error) {
    paths := make([]string, 0)

    for _, includedPath := range c.IncludedProjects {
        includedConfig, err:= LoadConfigurationFromFile(filepath.Join(configPath, includedPath))
        if err != nil {
            return nil, fmt.Errorf("Cannot load configuration from included project '%s': %v", filepath.Join(configPath, includedPath), err)
        }
        includedPath = filepath.Join(configPath, includedPath, includedConfig.ProjectTarget)
        paths = append(paths, includedPath)
    }

    return paths, nil
}
