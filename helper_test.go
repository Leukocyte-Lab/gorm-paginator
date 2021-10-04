package paginator

import (
	"reflect"
	"testing"
)

func Test_GenPage(t *testing.T) {
	type args struct {
		pageNo   int
		pageSize int
	}
	tests := []struct {
		name string
		args args
		want Page
	}{
		// test case
		{
			name: "General",
			args: args{
				pageNo:   1,
				pageSize: 10,
			},
			want: Page{
				Number: 1,
				Size:   10,
				Total:  0,
			},
		},
		{
			name: "Large PageSize",
			args: args{
				pageNo:   1,
				pageSize: 999999,
			},
			want: Page{
				Number: 1,
				Size:   100,
				Total:  0,
			},
		},
		{
			name: "Zero PageSize",
			args: args{
				pageNo:   1,
				pageSize: 0,
			},
			want: Page{
				Number: 1,
				Size:   20,
				Total:  0,
			},
		},
		{
			name: "Nagtive PageSize",
			args: args{
				pageNo:   1,
				pageSize: -1,
			},
			want: Page{
				Number: 1,
				Size:   20,
				Total:  0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenPage(tt.args.pageNo, tt.args.pageSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GenPage() = %v, want %v", got, tt.want)
			}
		})
	}
}
