package model

import (
	"github.com/oklog/ulid/v2"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Bind struct {
	gorm.Model
	ID       string         `json:"id" gorm:"primary_key;type:char(26)"`  // 绑定ID
	UID      string         `gorm:"not null;column:uid;index"`            // 对应User表的ID
	Platform string         `gorm:"not null;column:platform"`             // 平台类型。定义在常量中
	OpenID   string         `gorm:"not null;column:open_id;index;unique"` // 第三方平台唯一id
	Attr     datatypes.JSON `gorm:"column:attr"`                          // 更多属性 根据OAuth方不同自行设置
}

func (b *Bind) BeforeCreate(_ *gorm.DB) (err error) {
	b.ID = ulid.Make().String()
	return
}
