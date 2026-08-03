package helper

import "github.com/tubagusmf/ivlolitas-be/internal/model"

func NormalizePagination(filter *model.ProductFilter) {
	if filter.Page <= 0 {
		filter.Page = 1
	}

	if filter.Limit <= 0 {
		filter.Limit = 10
	}

	if filter.Limit > 100 {
		filter.Limit = 100
	}
}
