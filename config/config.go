package config

type Config struct {
	Log         Log         `yaml:"log"`
	Metrics     Metrics     `yaml:"metrics"`
	Transponder Transponder `yaml:"transponder"`
	Server      Server      `yaml:"server"`
}
