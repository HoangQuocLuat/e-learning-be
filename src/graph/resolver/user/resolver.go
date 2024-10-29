package resolver_user

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{}

func CalculateTotalPage(total, limit int) int {
	if limit == 0 {
		limit = 10
	}
	totalPage := total / limit
	if total%limit != 0 {
		totalPage++
	}

	return totalPage
}
