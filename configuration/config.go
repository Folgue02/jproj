package configuration

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
    fileContent, err := ioutil.ReadFile(filePath)

    if err != nil {
        return nil, err
    }

    return LoadConfigurationFromString(string(fileContent))
}
