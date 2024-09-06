package configs

type ServerConfig struct {
	System System `mapstructure:"system"`
	Nodes  []Node `mapstructure:"nodes"`
	Docker Docker `mapstructure:"docker"`
}
