package paginator

import "gorm.io/gorm/utils/tests"

var mockUsers = []tests.User{
	{
		Name:   "Jane",
		Active: false,
	},
	{
		Name:   "Jack",
		Active: true,
	},
	{
		Name:   "Jill",
		Active: true,
	},
	{
		Name:   "John",
		Active: true,
	},
	{
		Name:   "Julia",
		Active: true,
	},
}
