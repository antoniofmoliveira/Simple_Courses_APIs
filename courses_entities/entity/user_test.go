package entity

import (
	"testing"
)

func TestNewUser(t *testing.T) {
	type args struct {
		name     string
		email    string
		password string
	}
	tests := []struct {
		name    string
		args    args
		want    *User
		wantErr bool
	}{
		{
			name: "test",
			args: args{
				name:     "test",
				email:    "test@test.com",
				password: "test",
			},
			want: &User{
				ID:       "",
				Name:     "test",
				Email:    "test@test.com",
				Password: "test",
			},
			wantErr: false,
		},
		{
			name: "test1",
			args: args{
				name:     "",
				email:    "test@test.com",
				password: "test",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test2",
			args: args{
				name:     "test",
				email:    "",
				password: "test",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "test3",
			args: args{
				name:     "test",
				email:    "test@test.com",
				password: "",
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewUser(tt.args.name, tt.args.email, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewUser() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err == nil {
				if tt.args.name != got.Name || tt.args.email != got.Email || got.ID == "" || got.ValidatePassword(tt.args.password) == false {
					t.Errorf("NewUser() got = %v, want %v", got, tt.want)
				}
			}
		})
	}
}
