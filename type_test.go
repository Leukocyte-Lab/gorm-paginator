package paginator

import "testing"

func TestDirection_String(t *testing.T) {
	tests := []struct {
		name string
		enum Direction
		want string
	}{
		{
			name: "SortASC",
			enum: SortASC,
			want: "ASC",
		},
		{
			name: "SortDESC",
			enum: SortDESC,
			want: "DESC",
		},
		{
			name: "Unknow",
			enum: Direction(0),
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.enum.String(); got != tt.want {
				t.Errorf("Direction.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
