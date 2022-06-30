package jackett

import (
	"context"
	"net/http"
	"net/url"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func Test_client_makeURL(t *testing.T) {
	type fields struct {
		apiURL     string
		apiKey     string
		httpClient *http.Client
	}
	type args struct {
		ctx     context.Context
		options []Option
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *http.Request
		wantErr bool
	}{
		{
			name: "query",
			fields: fields{
				apiURL:     "http://foo.bar",
				apiKey:     "bar",
				httpClient: nil,
			},
			args: args{
				ctx: context.Background(),
				options: []Option{
					WithQuery("hello"),
				},
			},
			want: &http.Request{
				Method: "GET",
				URL: func() *url.URL {
					u, _ := url.Parse("http://foo.bar/api/v2.0/indexers/all/results?Query=hello&apikey=bar")
					return u
				}(),
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header:     make(http.Header),
				Body:       http.NoBody,
				Host:       "foo.bar",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &client{
				apiURL:     tt.fields.apiURL,
				apiKey:     tt.fields.apiKey,
				httpClient: tt.fields.httpClient,
			}
			got, err := c.makeURL(tt.args.ctx, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("makeURL() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if diff := cmp.Diff(tt.want, got, cmpopts.IgnoreUnexported(http.Request{})); diff != "" {
				t.Errorf("makeURL() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
