package conf

import (
	"encoding/json"
	"os"
)

type Config struct {
	data map[string]any
}

func (c *Config) Load(filepath string) error {
	f, err := os.Open(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	var data map[string]any
	if err := json.NewDecoder(f).Decode(&data); err != nil {
		return err
	}

	c.data = data
	return nil
}

func (c *Config) Value(key string) (string, bool) {
	v, ok := c.data[key]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	if ok {
		return s, true
	}
	return "", true
}

func (c *Config) SubPath(prefix string) func(string) (string, bool) {
	return func(key string) (string, bool) {
		data, ok := c.data[prefix].(map[string]any)
		if !ok {
			return "", false
		}
		s, ok := data[key].(string)
		if ok {
			return s, true
		}
		return "", false
	}
}
