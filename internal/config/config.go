package config

type Config struct {
	WorkDir  string `yaml:"workDir"`
	MaxDepth int    `yaml:"maxDepth"`
	Base     string `yaml:"base"`
}
