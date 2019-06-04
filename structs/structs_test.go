package structs

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_structToMap(t *testing.T) {
	type args struct {
		in interface{}
	}
	tests := []struct {
		name    string
		wantErr bool
		args    *args
		want    map[string]interface{}
	}{
		{
			name: "normal struct",
			args: &args{
				in: struct {
					ID string
				}{
					ID: "123",
				},
			},
			want: map[string]interface{}{
				"ID": "123",
			},
		},
		{
			name: "allow struct tags",
			args: &args{
				in: struct {
					ID string `json:"my_id"`
				}{
					ID: "123",
				},
			},
			want: map[string]interface{}{
				"my_id": "123",
			},
		},
		{
			name: "allow passing a map",
			args: &args{
				in: map[string]interface{}{
					"id": "123",
				},
			},
			want: map[string]interface{}{
				"id": "123",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := StructToMap(tt.args.in)
			if (gotErr != nil) != tt.wantErr {
				t.Errorf("StructToMap() error = %v, wantErr %v", gotErr, tt.wantErr)
				return
			}
			require.Equal(t, tt.want, got)
		})
	}
}
