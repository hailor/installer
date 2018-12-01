package config

type Node struct{
	Ipv4Address string `yaml:"ipv4Address"`
}

type Config struct{
	Nodes map[string]Node `yaml:"nodes"`
	Masters []string `yaml:"masters"`
	Workers []string `yaml:"workers"`
	Etcds []string `yaml:"etcds"`
}