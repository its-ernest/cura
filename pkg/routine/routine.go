package routine

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Trigger struct {
	Type   string `yaml:"type" json:"type"`
	Target string `yaml:"target" json:"target"`
}

type Action struct {
	Type   string      `yaml:"type" json:"type"`
	Target string      `yaml:"target,omitempty" json:"target,omitempty"`
	Value  interface{} `yaml:"value,omitempty" json:"value,omitempty"`
}

type StopCondition struct {
	Type   string `yaml:"type" json:"type"`
	Target string `yaml:"target" json:"target"`
}

type Routine struct {
	Name          string        `yaml:"name" json:"name"`
	Enabled       bool          `yaml:"enabled" json:"enabled"`
	Path          string        `yaml:"-" json:"-"`
	Trigger       Trigger       `yaml:"trigger" json:"trigger"`
	Actions       []Action      `yaml:"actions" json:"actions"`
	StopCondition StopCondition `yaml:"stop_condition" json:"stop_condition"`
	IsActive      bool          `yaml:"-" json:"isActive"` // camelCase for JS
}

// LoadRoutine reads a YAML file and returns a Routine pointer
func LoadRoutine(path string) (*Routine, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var r Routine
	err = yaml.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}
	return &r, nil
}
