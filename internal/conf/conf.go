package conf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type JSONFile struct {
	Path   string
	values map[string]string
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
	err = ioutil.WriteFile(j.Path, data, 0644)
	if err != nil {
		return fmt.Errorf("could not write to configuration file \"%s\": %w", j.Path, err)
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

// reload reads the JSON file at j.Path and updates j.values to match it. If
// j.Path doesn't exist, j.values becomes just an empty map.
func (j *JSONFile) reload() error {
	var data []byte
	if _, err := os.Stat(j.Path); !os.IsNotExist(err) {
		data, err = ioutil.ReadFile(j.Path)
		if err != nil {
			return fmt.Errorf("could not read configuration file \"%s\": %w", j.Path, err)
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
			j.Path, err,
		)
	}
	j.values = values
	return nil
}
