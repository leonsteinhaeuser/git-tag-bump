package branch

import (
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

// ReadConfig opens the config file at the given path.
func ReadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	cfg := &Config{}
	err = yaml.Unmarshal(file, cfg)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

type Config struct {
	Major Identifier `yaml:"major"`
	Minor Identifier `yaml:"minor"`
	Patch Identifier `yaml:"patch"`
}

type Identifier struct {
	Branch BranchIdentifier `yaml:"branch"`
}

func (i *Identifier) match(value string) bool {
	return i.Branch.match(value)
}

type BranchIdentifier struct {
	Name RegExIdentifier `yaml:"name"`
}

func (bi *BranchIdentifier) match(value string) bool {
	return bi.Name.match(value)
}

type RegExIdentifier struct {
	RegEx string `yaml:"regex"`
}

// match returns true if the given name matches the regex.
func (ri RegExIdentifier) match(value string) bool {
	return regexp.MustCompile(ri.RegEx).MatchString(value)
}
