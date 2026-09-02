package models

type SiteFlowModel struct {
	Model
	Count int `json:"count"`
}

func (SiteFlowModel) TableName() string {
	return "site_flow"
}
