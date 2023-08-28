package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
)

type Configuration struct {
    ProjectName string `json:"project_name"`
    ProjectTarget string `json:"project_target"`
}

func LoadConfigurationFromString(rawString string) (*Configuration, error){
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

    fileContent, err := ioutil.ReadFile(filePath)

    if err != nil {
        return nil, err
    }

    return LoadConfigurationFromString(string(fileContent))
}
