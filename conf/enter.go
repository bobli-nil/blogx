package conf

type Config struct {
	System System `yaml:"system"`
	Log    Log    `yaml:"log"`
	Jwt    Jwt    `yaml:"jwt"`
	Redis  Redis  `yaml:"redis"`
	DB     DB     `yaml:"db"`
}
