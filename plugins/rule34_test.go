package plugins

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"git.sr.ht/~jrswab/akordo/plugins/mocks"
	dg "github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
)

func TestRule34Request_Rule34(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w,
			`<?xml version="1.0" encoding="UTF-8"?><posts count="1414775" offset="0">`+
				`<post sample_url="https://us.fake.url/samples/3324/sample_000.jpg"/>`+
				`<post sample_url="https://us.fake.url/samples/3324/sample_000.jpg"/></posts>`)
	}))
	defer ts.Close()

	sfwChannel := &dg.Channel{
		ID:   "sfw",
		NSFW: false,
	}
	nsfwChannel := &dg.Channel{
		ID:   "nsfw",
		NSFW: true,
	}

	mockChan := new(mocks.AkSession)
	mockChan.On("Channel", mock.Anything).Return(sfwChannel, nil).Once()
	mockChan.On("Channel", mock.Anything).Return(nsfwChannel, nil)

	userMap := make(map[string]time.Time)
	userMap["2222"] = time.Now()

	type fields struct {
		record  *Record
		baseURL string
	}
	type args struct {
		req []string
		s   AkSession
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
			name: "Catch sfw channel",
			fields: fields{
				record:  &Record{LastReq: userMap, MinWaitTime: (2 * time.Minute)},
				baseURL: fmt.Sprintf("%s?page=dapi&s=post&q=index&tags=", ts.URL),
			},
			args: args{
				req: []string{"~testCMD", "test"},
				s:   mockChan,
				msg: &dg.MessageCreate{
					&dg.Message{
						ChannelID: "sfw",
						Author: &dg.User{
							ID:       "1111",
							Username: "user",
						},
					},
				},
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "Catch improper command format",
			fields: fields{
				record:  &Record{LastReq: userMap, MinWaitTime: (2 * time.Minute)},
				baseURL: fmt.Sprintf("%s?page=dapi&s=post&q=index&tags=", ts.URL),
			},
			args: args{
				req: []string{"~testCMD"},
				s:   mockChan,
				msg: &dg.MessageCreate{
					&dg.Message{
						ChannelID: "nsfw",
						Author: &dg.User{
							ID:       "2222",
							Username: "user2",
						},
					},
				},
			},
			want:    "Usage: `<prefix>rule34 tag`",
			wantErr: false,
		},
		{
			name: "Requested too soon",
			fields: fields{
				record:  &Record{LastReq: userMap, MinWaitTime: (2 * time.Minute)},
				baseURL: fmt.Sprintf("%s?page=dapi&s=post&q=index&tags=", ts.URL),
			},
			args: args{
				req: []string{"~testCMD", "test"},
				s:   mockChan,
				msg: &dg.MessageCreate{
					&dg.Message{
						ChannelID: "nsfw",
						Author: &dg.User{
							ID:       "2222",
							Username: "user2",
						},
					},
				},
			},
			want:    "user2 please wait 120 seconds before requesting the same command.",
			wantErr: false,
		},
		{
			name: "Return message on incorrect formatting",
			fields: fields{
				record:  &Record{LastReq: userMap, MinWaitTime: (2 * time.Minute)},
				baseURL: fmt.Sprintf("%s?page=dapi&s=post&q=index&tags=", ts.URL),
			},
			args: args{
				req: []string{"~testCMD"},
				s:   mockChan,
				msg: &dg.MessageCreate{
					&dg.Message{
						Author: &dg.User{
							ID:       "3333",
							Username: "user3",
						},
					},
				},
			},
			want:    "Usage: `<prefix>rule34 tag`",
			wantErr: false,
		},
		{
			name: "Invalid URL",
			fields: fields{
				record:  &Record{LastReq: userMap, MinWaitTime: (2 * time.Minute)},
				baseURL: "fake.url",
			},
			args: args{
				req: []string{"~testCMD", "test"},
				s:   mockChan,
				msg: &dg.MessageCreate{
					&dg.Message{
						Author: &dg.User{
							ID:       "4444",
							Username: "user4",
						},
					},
				},
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Valid URL",
			fields: fields{
				record:  &Record{LastReq: userMap, MinWaitTime: (2 * time.Minute)},
				baseURL: fmt.Sprintf("%s?page=dapi&s=post&q=index&tags=", ts.URL),
			},
			args: args{
				req: []string{"~testCMD", "test"},
				s:   mockChan,
				msg: &dg.MessageCreate{
					&dg.Message{
						Author: &dg.User{
							ID:       "5555",
							Username: "user5",
						},
					},
				},
			},
			want:    "https://us.fake.url/samples/3324/sample_000.jpg",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rr := &Rule34Request{
				record:  tt.fields.record,
				baseURL: tt.fields.baseURL,
			}
			got, err := rr.Rule34(tt.args.req, tt.args.s, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Rule34Request.Rule34() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Rule34Request.Rule34() = %v, want %v", got, tt.want)
			}
		})
	}
}
