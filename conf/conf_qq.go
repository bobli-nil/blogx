package conf

type QQ struct {
	AppID    string `yaml:"appID" json:"appID"`
	AppKey   string `yaml:"appKey" json:"appKey"`
	Redirect string `yaml:"redirect" json:"redirect"`
}

// TODO 后面完善
func (q QQ) Url() string {
	return "QQ的url"
}
