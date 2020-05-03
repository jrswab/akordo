package plugins

import (
	"fmt"
	"testing"
	"time"

	dg "github.com/bwmarrin/discordgo"
)

func TestRecord_Pong(t *testing.T) {
	userMap := make(map[string]time.Time)
	userMap["1111"] = time.Now()
	type fields struct {
		MinWaitTime time.Duration
		LastReq     map[string]time.Time
	}
	type args struct {
		msg *dg.MessageCreate
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "Requested too soon",
			fields: fields{
				LastReq:     userMap,
				MinWaitTime: (botDelay),
			},
			args: args{
				msg: &dg.MessageCreate{
					&dg.Message{
						Author: &dg.User{
							ID:       "1111",
							Username: "user1",
						},
					},
				},
			},
			want: fmt.Sprintf("user1 please wait %d seconds before requesting the same command.", CommandDelay),
		},
		{
			name: "Return pong",
			fields: fields{
				LastReq:     userMap,
				MinWaitTime: (botDelay),
			},
			args: args{
				msg: &dg.MessageCreate{
					&dg.Message{
						Author: &dg.User{
							ID:       "2222",
							Username: "user2",
						},
					},
				},
			},
			want: "pong",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Record{
				MinWaitTime: tt.fields.MinWaitTime,
				LastReq:     tt.fields.LastReq,
			}
			if got := r.Pong(tt.args.msg); got != tt.want {
				t.Errorf("Record.Pong() = %v, want %v", got, tt.want)
			}
		})
	}
}
