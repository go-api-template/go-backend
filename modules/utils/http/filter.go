package httputil

import (
	"github.com/gin-gonic/gin"
	"github.com/go-mods/convert"
	"gorm.io/gorm"
)

type Filter struct {
	Page    *int
	Limit   *int
	Offset  *int
	SortBy  []string
	OrderBy []string
	Search  string
}

// GetFilter creates a filter from the query parameters
func GetFilter(c *gin.Context) Filter {
	var filter Filter

	// Get the page number
	if page, ok := c.GetQuery("page"); ok {
		if page, err := convert.ToInt(page); err == nil {
			filter.Page = convert.ToPtr(page)
		}
	}

	// Get the limit
	if limit, ok := c.GetQuery("limit"); ok {
		if limit, err := convert.ToInt(limit); err == nil {
			filter.Limit = convert.ToPtr(limit)
		}
	}

	// Get the offset
	if filter.Page != nil && *filter.Page > 1 && filter.Limit != nil {
		page := *filter.Page
		limit := *filter.Limit
		filter.Offset = convert.ToPtr((page - 1) * limit)
	}

	// Get the sort by
	if sortBy, ok := c.GetQueryArray("sort_by"); ok {
		filter.SortBy = sortBy
	}

	// Get the order by
	if orderBy, ok := c.GetQueryArray("order_by"); ok {
		filter.OrderBy = orderBy
	}

	// Get the search
	if search, ok := c.GetQuery("search"); ok {
		filter.Search = search
	}

	return filter
}

// Apply filter to the query
func (f Filter) Apply(query *gorm.DB) *gorm.DB {
	if f.Limit != nil {
		query = query.Limit(*f.Limit)
	}
	if f.Offset != nil {
		query = query.Offset(*f.Offset)
	}
	if len(f.SortBy) > 0 {
		for i, sortBy := range f.SortBy {
			if len(f.OrderBy) > i {
				query = query.Order(sortBy + " " + f.OrderBy[i])
			} else {
				query = query.Order(sortBy)
			}
		}
	}
	if f.Search != "" {
		query = query.Where(f.Search)
	}
	return query
}
