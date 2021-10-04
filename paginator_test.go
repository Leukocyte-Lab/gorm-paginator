package paginator

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/utils/tests"
)

func TestPaginator_New(t *testing.T) {
	type args struct {
		page   Page
		orders []Order
		filter map[string]string
	}
	tests := []struct {
		name string
		args args
		want *Paginator
	}{
		{
			name: "General",
			args: args{
				page: Page{
					Number: 1,
					Size:   25,
				},
				orders: []Order{
					{
						Column:    "id",
						Direction: SortASC,
					},
				},
				filter: map[string]string{
					"id": "excludeID",
				},
			},
			want: &Paginator{
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
				Filter: map[string]string{
					"id": "excludeID",
				},
			},
		},
		{
			name: "Nagtive Page",
			args: args{
				page: Page{
					Number: -1,
					Size:   25,
				},
			},
			want: &Paginator{
				Page: Page{
					Number: MinPageNumber,
					Size:   25,
				},
			},
		},
		{
			name: "Zero Page",
			args: args{
				page: Page{
					Number: 0,
					Size:   25,
				},
			},
			want: &Paginator{
				Page: Page{
					Number: MinPageNumber,
					Size:   25,
				},
			},
		},
		{
			name: "Nagtive Page Size",
			args: args{
				page: Page{
					Number: 1,
					Size:   -1,
				},
			},
			want: &Paginator{
				Page: Page{
					Number: 1,
					Size:   MinPageSize,
				},
			},
		},
		{
			name: "Zero Page Size",
			args: args{
				page: Page{
					Number: 1,
					Size:   0,
				},
			},
			want: &Paginator{
				Page: Page{
					Number: 1,
					Size:   MinPageSize,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.page, tt.args.orders, tt.args.filter); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		{
			name: "Page2",
			fields: fields{
				Page: Page{
					Number: 2,
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
			want: `SELECT * FROM "users" WHERE "users"."deleted_at" IS NULL ORDER BY "id" LIMIT 25 OFFSET 25`,
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

func TestPaginator_CountPageTotal(t *testing.T) {
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
		name    string
		fields  fields
		args    args
		wantErr bool
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
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pgntr := New(
				tt.fields.Page,
				tt.fields.Order,
				tt.fields.Filter,
			)
			tx := pgntr.GenGormTransaction(tt.args.db).Find(tt.args.dest)

			if err := pgntr.CountPageTotal(tx); (err != nil) != tt.wantErr {
				t.Errorf("Paginator.CountPageTotal() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPaginator_Paginator(t *testing.T) {
	resource := setupResource(t)
	db := setupTestDB(t, resource)

	type fields struct {
		Page   Page
		Order  []Order
		Filter map[string]string
	}
	type args struct {
		dest []*tests.User
	}
	tests := []struct {
		name     string
		fields   fields
		args     args
		wantPage Page
		wantDest []*tests.User
		wantErr  bool
	}{
		// test case
		{
			name: "General",
			fields: fields{
				Page: Page{
					Number: 1,
					Size:   2,
				},
			},
			args: args{
				dest: []*tests.User{},
			},
			wantPage: Page{
				Number: 1,
				Size:   2,
				Total:  3,
			},
			wantDest: []*tests.User{
				{
					Name: "Jane",
				},
				{
					Name: "Jack",
				},
			},
			wantErr: false,
		},
		{
			name: "Page 2",
			fields: fields{
				Page: Page{
					Number: 2,
					Size:   2,
				},
			},
			args: args{
				dest: []*tests.User{},
			},
			wantPage: Page{
				Number: 2,
				Size:   2,
				Total:  3,
			},
			wantDest: []*tests.User{
				{
					Name: "Jill",
				},
				{
					Name: "John",
				},
			},
			wantErr: false,
		},
		{
			name: "DESC Page 2",
			fields: fields{
				Page: Page{
					Number: 2,
					Size:   2,
				},
				Order: []Order{
					{
						Column:    "id",
						Direction: SortDESC,
					},
				},
			},
			args: args{
				dest: []*tests.User{},
			},
			wantPage: Page{
				Number: 2,
				Size:   2,
				Total:  3,
			},
			wantDest: []*tests.User{
				{
					Name: "Jill",
				},
				{
					Name: "Jack",
				},
			},
			wantErr: false,
		},
		{
			name: "Filter Page 2",
			fields: fields{
				Page: Page{
					Number: 2,
					Size:   2,
				},
				Filter: map[string]string{
					"active": "true",
				},
			},
			args: args{
				dest: []*tests.User{},
			},
			wantPage: Page{
				Number: 2,
				Size:   2,
				Total:  2,
			},
			wantDest: []*tests.User{
				{
					Name: "John",
				},
				{
					Name: "Julia",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := assert.New(t)

			pgntr := New(
				tt.fields.Page,
				tt.fields.Order,
				tt.fields.Filter,
			)
			tx := pgntr.GenGormTransaction(db).Find(&tt.args.dest)
			err := pgntr.CountPageTotal(tx)
			gotDest := tt.args.dest
			gotPage := pgntr.Page

			if (err != nil) != tt.wantErr {
				t.Errorf("Paginator error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil {
				assert.NotEmpty(gotPage)
				assert.Equal(tt.wantPage, gotPage)
				assert.NotEmpty(gotDest)
				assert.Len(gotDest, len(tt.wantDest))
				for index, user := range tt.wantDest {
					assert.Equal(user.Name, gotDest[index].Name)
					// TODO: other fileds
				}
			}
		})
	}

	cleanResource(t, resource)
}
