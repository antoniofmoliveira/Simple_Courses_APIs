package entity

import (
	"reflect"
	"testing"
)

func TestNewCourse(t *testing.T) {
	type args struct {
		id          string
		name        string
		description string
		categoryID  string
	}
	tests := []struct {
		name string
		args args
		want *Course
	}{
		{
			name: "test",
			args: args{
				id:          "1",
				name:        "test",
				description: "test",
				categoryID:  "1",
			},
			want: &Course{
				ID:          "1",
				Name:        "test",
				Description: "test",
				CategoryID:  "1",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCourse(tt.args.id, tt.args.name, tt.args.description, tt.args.categoryID); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCourse() = %v, want %v", got, tt.want)
			}
		})
	}
}
