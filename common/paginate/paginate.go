/*
@Time : 2023/4/12 15:12
@Author : Hhx06
@File : paginate
@Description:
@Software: GoLand
*/

package paginate

import "math"

type paginate struct {
	Offset       int
	Limit        int
	Page         int
	LastPage     int
	TotalResults int
}

func NewPaginate(page, limit, totalResults int) *paginate {
	pager := &paginate{
		Offset:       0,
		Limit:        limit,
		Page:         page,
		TotalResults: totalResults,
	}
	pager.build()
	return pager
}

func (c *paginate) build() {
	c.setDefaults()

	c.ceilLastPage()
	c.doPaginate()
}

// 计算 offset
func (c *paginate) doPaginate() {
	if c.Page == 0 {
		c.Page = 1
	}
	c.Offset = (c.Page - 1) * c.Limit
}

func (c *paginate) ceilLastPage() {
	if c.TotalResults == 0 {
		c.LastPage = 0
		return
	}
	c.LastPage = int(math.Ceil(float64(c.TotalResults) / float64(c.Limit)))
	return
}

func (c *paginate) setDefaults() {
	if c.Page == 0 {
		c.Page = 1
	}

	if c.Limit == 0 {
		c.Limit = 10
	}
}
