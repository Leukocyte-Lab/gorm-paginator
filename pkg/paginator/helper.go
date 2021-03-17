package paginator

import (
	pb "github.com/Leukocyte-Lab/AGH2-Proto/go/pagination/v1"
)

func GenPage(pageNo int, pageSize int) Page {
	if pageNo < MinPageNumber {
		pageNo = MinPageNumber
	}
	switch {
	case pageSize > MaxPageSize:
		pageSize = MaxPageSize
	case pageSize < MinPageSize:
		pageSize = DefaultPageSize
	}

	return Page{
		Number: pageNo,
		Size:   pageSize,
	}
}

func CastPbOrders2Order(pbOrder []*pb.Order) []Order {
	var orders []Order
	for _, order := range pbOrder {
		var direction string
		switch order.GetDirection() {
		case pb.Order_DIRECTION_ASC:
			direction = SortASC
		case pb.Order_DIRECTION_DESC:
			direction = SortDESC
		default:
			continue
		}

		orders = append(orders, Order{
			Column:    order.GetColumnName(),
			Direction: direction,
		})
	}

	return orders
}
