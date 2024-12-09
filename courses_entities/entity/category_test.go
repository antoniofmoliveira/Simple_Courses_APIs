package entity

import (
	"reflect"
	"testing"
)

func TestNewCategory(t *testing.T) {
	type args struct {
		id          string
		name        string
		description string
	}
	tests := []struct {
		name string
		args args
		want *Category
	}{
		{
			name: "test",
			args: args{
				id:          "1",
				name:        "test",
				description: "test",
			},
			want: &Category{
				ID:          "1",
				Name:        "test",
				Description: "test",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCategory(tt.args.id, tt.args.name, tt.args.description); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCategory() = %v, want %v", got, tt.want)
			}
		})
	}
}
