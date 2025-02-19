package class

import "testing"

func TestClassHandler_ConvertWeek(t *testing.T) {
	type args struct {
		weeks []int
	}
	tests := []struct {
		name string
		c    *ClassHandler
		args args
		want int64
	}{
		{
			name: "Single week",
			c:    &ClassHandler{},
			args: args{weeks: []int{1}},
			want: 1,
		},
		{
			name: "Multiple weeks",
			c:    &ClassHandler{},
			args: args{weeks: []int{1, 2, 3}},
			want: 7,
		},
		{
			name: "Non-consecutive weeks",
			c:    &ClassHandler{},
			args: args{weeks: []int{1, 3, 5}},
			want: 21,
		},
		{
			name: "Weeks out of range",
			c:    &ClassHandler{},
			args: args{weeks: []int{0, 31}},
			want: 0,
		},
		{
			name: "Mixed valid and invalid weeks",
			c:    &ClassHandler{},
			args: args{weeks: []int{1, 0, 3, 31}},
			want: 5,
		},
		{
			name: "Empty weeks",
			c:    &ClassHandler{},
			args: args{weeks: []int{}},
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.c.ConvertWeek(tt.args.weeks); got != tt.want {
				t.Errorf("ClassHandler.ConvertWeek() = %v, want %v", got, tt.want)
			}
		})
	}
}
