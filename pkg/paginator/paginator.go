package paginator

import (
	pb "github.com/Leukocyte-Lab/AGH2-Proto/go/pagination/v1"
	"github.com/kataras/iris/v12"
	"google.golang.org/protobuf/encoding/protojson"
)

const (
	DefaultPageSize = 25
)

func Paginator(ctx iris.Context) (number int, size int, orders []*pb.Order, err error) {
	// querystring /?page={uint}
	number = ctx.URLParamIntDefault("page", 1)

	// querystring /?size={uint}
	size = ctx.URLParamIntDefault("size", DefaultPageSize)

	// querystring /?order={"column_name":"{ColumnName}", "direction":"{DIRECTION_ASC || DIRECTION_DESC}"}
	for _, each := range ctx.Request().URL.Query()["order"] {
		order := pb.Order{}
		err = protojson.Unmarshal([]byte(each), &order)
		if err != nil {
			return -1, -1, nil, err
		}
		orders = append(orders, &order)
	}

	return number, size, orders, nil
}
