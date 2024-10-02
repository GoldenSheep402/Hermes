package model

type Subnet struct {
	CIDR    string `gorm:"column:cidr"`
	IsAllow bool   `gorm:"column:is_allow"`
}
