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

func TestNewMemeRequest(t *testing.T) {
	type args struct {
		url []string
	}
	tests := []struct {
		name string
		args args
		want *MemeRequest
	}{
		{
			name: "Called with no params",
			args: args{},
			want: &MemeRequest{
				record:  NewRecorder(),
				baseURL: fmt.Sprintf("https://memegen.link/api/templates/"),
			},
		},
		{
			name: "Called with one params",
			args: args{url: []string{"fake.url"}},
			want: &MemeRequest{
				record:  NewRecorder(),
				baseURL: "fake.url",
			},
		},
		{
			name: "Called with more than one params",
			args: args{url: []string{"fake.url", "another.url"}},
			want: &MemeRequest{
				record:  NewRecorder(),
				baseURL: "fake.url",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMemeRequest(tt.args.url...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMemeRequest() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemeRequest_RequestMeme(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{"direct":{"masked": "somef.ake/url/withImage.gif"}}`)
	}))
	defer ts.Close()

	userMap := make(map[string]time.Time)
	userMap["3333"] = time.Now()

	listMsg0 := "To see all available memes head to https://memegen.link/api/templates/\n"
	listMsg1 := "Use the name at the end of the URLs that are displayed."
	listOutput := fmt.Sprintf("%s%s", listMsg0, listMsg1)
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
			name: "meme list",
			fields: fields{
				record:  &Record{LastReq: userMap, MinWaitTime: (botDelay)},
				baseURL: ts.URL,
			},
			args: args{
				req: []string{"~testCMD", "list"},
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
			want:    listOutput,
			wantErr: false,
		},
		{
			name: "les than 3 args and not list",
			fields: fields{
				record:  &Record{LastReq: userMap, MinWaitTime: (botDelay)},
				baseURL: ts.URL,
			},
			args: args{
				req: []string{"~testCMD", "hmmm"},
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
			want:    "Usage: `<prefix>meme name top_text <bottom_text>`",
			wantErr: false,
		},
		{
			name: "Command with more than 3 arguments",
			fields: fields{
				record:  &Record{LastReq: userMap, MinWaitTime: (botDelay)},
				baseURL: ts.URL,
			},
			args: args{
				req: []string{"~testCMD", "one", "two", "three", "four"},
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
			want:    "Usage: `<prefix>meme name top_text <bottom_text>`",
			wantErr: false,
		},
		{
			name: "Requested too soon",
			fields: fields{
				record:  &Record{LastReq: userMap, MinWaitTime: (botDelay)},
				baseURL: ts.URL,
			},
			args: args{
				req: []string{"~testCMD", "test", "one", "two"},
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
			want:    fmt.Sprintf("user3 please wait %d seconds before requesting the same command.", CommandDelay),
			wantErr: false,
		},
		{
			name: "Correct URL",
			fields: fields{
				record:  &Record{LastReq: userMap, MinWaitTime: (botDelay)},
				baseURL: fmt.Sprintf("%s/", ts.URL),
			},
			args: args{
				req: []string{"~testCMD", "test", "two", "three"},
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
			want:    "somef.ake/url/withImage.gif",
			wantErr: false,
		},
		{
			name: "Bad URL",
			fields: fields{
				record:  &Record{LastReq: userMap, MinWaitTime: (botDelay)},
				baseURL: "wron.gurl",
			},
			args: args{
				req: []string{"~testCMD", "test", "two", "three"},
				s:   &dg.Session{},
				msg: &dg.MessageCreate{
					&dg.Message{
						Author: &dg.User{
							ID:       "5555",
							Username: "user5",
						},
					},
				},
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemeRequest{
				record:  tt.fields.record,
				baseURL: tt.fields.baseURL,
			}
			got, err := m.RequestMeme(tt.args.req, tt.args.s, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemeRequest.RequestMeme() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MemeRequest.RequestMeme() = %v, want %v", got, tt.want)
			}
		})
	}
}
