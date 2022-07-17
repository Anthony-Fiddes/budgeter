package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// JSONFile represents a configuration file that is stored locally on disk and only
// contains string values. I.e. no key has an array value or an embedded object.
type JSONFile struct {
	path   string
	values map[string]string
}

func NewJSONFile(path string) *JSONFile {
	return &JSONFile{path: path}
}

func (j *JSONFile) Put(key, value string) error {
	if err := j.reload(); err != nil {
		return err
	}
	j.values[key] = value
	data, err := json.Marshal(j.values)
	if err != nil {
		return fmt.Errorf("could not marshal config to JSON: %w", err)
	}
	err = ioutil.WriteFile(j.path, data, 0644)
	if err != nil {
		return fmt.Errorf("could not write to configuration file \"%s\": %w", j.path, err)
	}
	return nil
}

func (j *JSONFile) Get(key string) (string, error) {
	if err := j.reload(); err != nil {
		return "", err
	}
	val, ok := j.values[key]
	if ok {
		return val, nil
	}
	return "", nil
}

func (j *JSONFile) GetAll() (map[string]string, error) {
	if err := j.reload(); err != nil {
		return nil, err
	}
	result := make(map[string]string, len(j.values))
	for k, v := range j.values {
		result[k] = v
	}
	return result, nil
}

// reload reads the JSON file at j.Path and updates j.values to match it. If
// j.Path doesn't exist, j.values becomes just an empty map.
func (j *JSONFile) reload() error {
	var data []byte
	if _, err := os.Stat(j.path); !os.IsNotExist(err) {
		data, err = ioutil.ReadFile(j.path)
		if err != nil {
			return fmt.Errorf("configuration file \"%s\" does not exist", j.path)
		}
	} else {
		j.values = map[string]string{}
		return nil
	}
	values := make(map[string]string)
	err := json.Unmarshal(data, &values)
	if err != nil {
		return fmt.Errorf(
			"configuration file \"%s\" does not contain JSON with only string values: %w",
			j.path, err,
		)
	}
	j.values = values
	return nil
}
