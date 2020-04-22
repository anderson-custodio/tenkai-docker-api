package model

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/jinzhu/gorm"
)

//ListDockerTagsRequest structure
type ListDockerTagsRequest struct {
	ImageName string `json:"imageName"`
	From      string `json:"from"`
}

//ListDockerTagsResult structure
type ListDockerTagsResult struct {
	TagResponse []TagResponse `json:"tags"`
}

//TagResponse Structure
type TagResponse struct {
	Created time.Time `json:"created"`
	Tag     string    `json:"tag"`
}

//DockerRepo structure
type DockerRepo struct {
	gorm.Model
	Host     string `json:"host"`
	Username string `json:"username"`
	Password string `json:"password"`
}

//TagsResult structure
type TagsResult struct {
	Name string
	Tags []string
}

//ManifestResult structure
type ManifestResult struct {
	History []ManifestHistory `json:"history"`
}

//ManifestHistory structure
type ManifestHistory struct {
	V1Compatibility string `json:"v1Compatibility"`
}

//V1Compatibility  structure
type V1Compatibility struct {
	Created time.Time `json:"created"`
}

//ListDockerRepositoryResponse structure
type ListDockerRepositoryResponse struct {
	Repositories []DockerRepo `json:"repositories"`
}

// DockerVariablesPayload structure
type DockerVariablesPayload struct {
	ImageName string `json:"imageName"`
	ImageTag  string `json:"imageTag"`
}

type dockerManifestConfig struct {
	Digest string `json:"digest"`
}

// DockerManifestResponse structure
type DockerManifestResponse struct {
	Config dockerManifestConfig `json:"config"`
}

// DockerVariablesResponse structure
type DockerVariablesResponse struct {
	ConfigJSON       []dockerVariable `json:"configJson"`
	GlobalConfigJSON []dockerVariable `json:"globalConfigJson"`
}

type dockerVariable struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

// UnmarshalJSON implements the unmarshal for ListDockerVariablesResponse
func (response *DockerVariablesResponse) UnmarshalJSON(dockerManifest []byte) error {
	var jsonRootMap map[string]interface{}
	if err := json.Unmarshal(dockerManifest, &jsonRootMap); err != nil {
		return err
	}

	propContainerConfigMap := jsonRootMap["container_config"].(map[string]interface{})
	propLabelsMap := propContainerConfigMap["Labels"].(map[string]interface{})

	response.ConfigJSON, _ = getVarsFromProp(propLabelsMap["configJson"])
	response.GlobalConfigJSON, _ = getVarsFromProp(propLabelsMap["globalConfigJson"])

	return nil
}

func getVarsFromProp(labels interface{}) ([]dockerVariable, error) {
	if labels == nil {
		return make([]dockerVariable, 0), nil
	}

	varsJSONStr := []byte(labels.(string))
	var varsMap map[string]interface{}
	if err := json.Unmarshal(varsJSONStr, &varsMap); err != nil {
		return make([]dockerVariable, 0), err
	}

	// Extracts variables array from map
	variables := varsMap["variables"].([]interface{})
	return marshalDockerVars(variables), nil
}

func marshalDockerVars(vars []interface{}) []dockerVariable {
	dockerVars := make([]dockerVariable, 0)
	for _, v := range vars {
		mapStrI := v.(map[string]interface{})

		var dv dockerVariable
		dv.Name = fmt.Sprint(mapStrI["name"])
		dv.Description = fmt.Sprint(mapStrI["description"])
		required, _ := strconv.ParseBool(fmt.Sprint(mapStrI["required"]))
		dv.Required = required

		dockerVars = append(dockerVars, dv)
	}
	return dockerVars
}
