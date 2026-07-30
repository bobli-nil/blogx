package conf

type Config struct {
	System System `yaml:"system"`
	Log    Log    `yaml:"log"`
	Jwt    Jwt    `yaml:"jwt"`
	Redis  Redis  `yaml:"redis"`
	DB     DB     `yaml:"db"`
	Site   Site   `yaml:"site"`
	Email  Email  `yaml:"email"`
	QQ     QQ     `yaml:"qq"`
	QiNiu  QiNiu  `yaml:"qiNiu"`
	Ai     Ai     `yaml:"ai"`
}
