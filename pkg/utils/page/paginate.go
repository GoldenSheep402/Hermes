package page

import "gorm.io/gorm"

type Paginate struct {
	Current  int `query:"current" json:"current"`
	PageSize int `query:"pageSize" json:"pageSize"`
}

func (p *Paginate) Paginate() func(db *gorm.DB) *gorm.DB {
	page := p.Current
	limit := p.PageSize
	return func(db *gorm.DB) *gorm.DB {
		if page == 0 {
			page = 1
		}
		switch {
		case limit > 100:
			limit = 100
		case limit <= 0:
			limit = 10
		}

		offset := (page - 1) * limit
		return db.Offset(offset).Limit(limit)
	}
}
