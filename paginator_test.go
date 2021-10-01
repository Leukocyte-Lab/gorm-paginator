package paginator

import (
	"testing"

	"gorm.io/gorm"
	"gorm.io/gorm/utils/tests"
)

func TestPaginator_GenGormTransaction(t *testing.T) {
	type fields struct {
		Page   Page
		Order  []Order
		Filter map[string]string
	}
	type args struct {
		db   *gorm.DB
		dest interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		// test case
		{
			name: "General",
			fields: fields{
				Page: Page{
					Number: 1,
					Size:   25,
				},
				Order: []Order{
					{
						Column:    "id",
						Direction: SortASC,
					},
				},
			},
			args: args{
				db:   setupMockDB("postgres"),
				dest: &[]tests.User{},
			},
			want: `SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL ORDER BY "id" LIMIT 25`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := New(
				tt.fields.Page,
				tt.fields.Order,
				tt.fields.Filter,
			).GenGormTransaction(tt.args.db).Find(tt.args.dest)

			got := tx.Statement.SQL.String()
			if got != tt.want {
				t.Errorf("Paginator.GenGormTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}
