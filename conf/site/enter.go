package site

type SiteInfo struct {
	Title string `yaml:"title"`
	Logo  string `yaml:"logo"`
	Beian string `yaml:"beian"`
	Mode  int8   `yaml:"mode"` // 1社区模式 2博客模式
}

type Project struct {
	Title   string `yaml:"title"`
	Icon    string `yaml:"icon"`
	WebPath string `yaml:"webPath"`
}

type Seo struct {
	Keywords    string `yaml:"keywords"`
	Description string `yaml:"description"`
}

type About struct {
	Version   string `yaml:"-" json:"version"`
	SiteAbout string `yaml:"siteAbout" json:"siteAbout"`
	QQ        string `yaml:"qq" json:"qq"`
	Wechat    string `yaml:"wechat" json:"wechat"`
	Gitee     string `yaml:"gitee" json:"gitee"`
	Bilibili  string `yaml:"bilibili" json:"bilibili"`
	Github    string `yaml:"github" json:"github"`
}

type Login struct {
	QQLogin          bool `yaml:"qqLogin" json:"qqLogin"`
	UsernamePwdLogin bool `yaml:"usernamePwdLogin" json:"usernamePwdLogin"`
	EmailPwdLogin    bool `yaml:"emailPwdLogin" json:"emailPwdLogin"`
	Captcha          bool `yaml:"captcha" json:"captcha"`
}

type ComponentInfo struct {
	Title  string `yaml:"title" json:"title"`
	Enable bool   `yaml:"enable" json:"enable"`
}
type IndexRight struct {
	List []ComponentInfo `yaml:"list" json:"list"`
}

type Article struct {
	NoExamine bool `yaml:"noExamine" json:"noExamine"`
}
