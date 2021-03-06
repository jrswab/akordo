package plugins

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	dg "github.com/bwmarrin/discordgo"
)

func TestNewRecorder(t *testing.T) {
	userMap := make(map[string]time.Time)
	tests := []struct {
		name string
		want *Record
	}{
		{
			name: "Create the Record struct",
			want: &Record{LastReq: userMap, MinWaitTime: (botDelay)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRecorder(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRecorder() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestRecord_checkLastAsk(t *testing.T) {
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
		name     string
		fields   fields
		args     args
		want     string
		wantBool bool
	}{
		{
			name: "User executed same command too soon",
			fields: fields{
				MinWaitTime: botDelay,
				LastReq:     userMap,
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
			want:     fmt.Sprintf("user1 please wait %d seconds before requesting the same command.", CommandDelay),
			wantBool: true,
		},
		{
			name: "User executed same command after timeout",
			fields: fields{
				MinWaitTime: botDelay,
				LastReq:     userMap,
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
			want:     "",
			wantBool: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Record{
				MinWaitTime: tt.fields.MinWaitTime,
				LastReq:     tt.fields.LastReq,
			}
			got, gotBool := r.CheckLastAsk(tt.args.msg)
			if got != tt.want {
				t.Errorf("Record.checkLastAsk() got = %v, want %v", got, tt.want)
			}
			if gotBool != tt.wantBool {
				t.Errorf("Record.checkLastAsk() got1 = %v, want %v", gotBool, tt.wantBool)
			}
		})
	}
}
