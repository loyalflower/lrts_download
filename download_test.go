package lrts_download

import (
	"net/url"
	"testing"
)

func Test_paramsMd5(t *testing.T) {
	type args struct {
		query   url.Values
		urlPath string
		secret  string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "http params md5",
			args: args{
				query: url.Values{
					"bookId":   []string{"5201"},
					"pageNum":  []string{"1"},
					"pageSize": []string{"50"},
					"sortType": []string{"0"},
					"token":    []string{"OqzlvCxt2i_P1SZKF6GjFg**_lK0uCQpm5tN-P6XdFZYawCDKSgeC4anU"},
					"imei":     []string{"MDI6MDA6MDA6MDA6MDA6MDA="},
					"nwt":      []string{"1"},
					"q":        []string{"1930"},
				},
				urlPath: urlPathBookList,
				secret:  secret,
			},
			want: "c29051c00df0f8aa40dbac91285b8cf9",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := paramsMd5(tt.args.query, tt.args.urlPath, tt.args.secret); got != tt.want {
				t.Errorf("paramsMd5() = %v, want %v", got, tt.want)
			}
		})
	}
}
