package model

type TenantService struct {
	ServiceID    string `json:"service_id" gorm:"column:service_id"`
	ServiceCName string `json:"service_cname" gorm:"column:service_cname"`
}

// 自定义表名
func (TenantService) TableName() string {
	return "tenant_service"
}
