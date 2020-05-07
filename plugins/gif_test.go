package plugins

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	dg "github.com/bwmarrin/discordgo"
)

func TestGifRequest_Gif(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w,
			`{"data":[{"embed_url": "somef.ake/url/withImage.gif"},`+
				`{"embed_url": "somef.ake/url/withImage.gif"}]}`)
	}))
	defer ts.Close()

	userMap := make(map[string]time.Time)
	userMap["1111"] = time.Now()

	type fields struct {
		record  *Record
		baseURL string
	}
	type args struct {
		req []string
		s   *dg.Session
		msg *dg.MessageCreate
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Requested too soon",
			fields: fields{
				record:  &Record{LastReq: userMap, MinWaitTime: (botDelay)},
				baseURL: ts.URL,
			},
			args: args{
				req: []string{"~testCMD", "test"},
				s:   &dg.Session{},
				msg: &dg.MessageCreate{
					&dg.Message{
						Author: &dg.User{
							ID:       "1111",
							Username: "user",
						},
					},
				},
			},
			want:    fmt.Sprintf("user please wait %d seconds before requesting the same command.", CommandDelay),
			wantErr: false,
		},
		{
			name: "Return message on incorrect formatting",
			fields: fields{
				record:  &Record{LastReq: userMap, MinWaitTime: (botDelay)},
				baseURL: ts.URL,
			},
			args: args{
				req: []string{"~testCMD"},
				s:   &dg.Session{},
				msg: &dg.MessageCreate{
					&dg.Message{
						Author: &dg.User{
							ID:       "2222",
							Username: "user2",
						},
					},
				},
			},
			want:    "Usage: `<prefix>gif word`",
			wantErr: false,
		},
		{
			name: "Invalid URL",
			fields: fields{
				record:  &Record{LastReq: userMap, MinWaitTime: (botDelay)},
				baseURL: "fake.url",
			},
			args: args{
				req: []string{"~testCMD", "test"},
				s:   &dg.Session{},
				msg: &dg.MessageCreate{
					&dg.Message{
						Author: &dg.User{
							ID:       "3333",
							Username: "user3",
						},
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Receive URL",
			fields: fields{
				record:  &Record{LastReq: userMap, MinWaitTime: (botDelay)},
				baseURL: fmt.Sprintf("%s/?rating=pg", ts.URL),
			},
			args: args{
				req: []string{"~testCMD", "test"},
				s:   &dg.Session{},
				msg: &dg.MessageCreate{
					&dg.Message{
						Author: &dg.User{
							ID:       "4444",
							Username: "user4",
						},
					},
				},
			},
			want:    "Here's what I found for *test* :wink: somef.ake/url/withImage.gif",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GifRequest{
				Record:  tt.fields.record,
				BaseURL: tt.fields.baseURL,
			}
			got, err := g.Gif(tt.args.req, tt.args.s, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("GifRequest.Gif() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GifRequest.Gif() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewGifRequest(t *testing.T) {
	type args struct {
		url []string
	}
	tests := []struct {
		name string
		args args
		want *GifRequest
	}{
		{
			name: "Called with no params",
			args: args{},
			want: &GifRequest{
				Record:  NewRecorder(),
				BaseURL: fmt.Sprintf("http://api.giphy.com/v1/gifs/search?rating=pg"),
			},
		},
		{
			name: "Called with one params",
			args: args{url: []string{"fake.url"}},
			want: &GifRequest{
				Record:  NewRecorder(),
				BaseURL: "fake.url",
			},
		},
		{
			name: "Called with more than one params",
			args: args{url: []string{"fake.url", "another.url"}},
			want: &GifRequest{
				Record:  NewRecorder(),
				BaseURL: "fake.url",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewGifRequest(tt.args.url...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewGifRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}
