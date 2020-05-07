package antispam

import (
	"os"
	"reflect"
	"testing"

	dg "github.com/bwmarrin/discordgo"
)

func TestNewSpamTracker(t *testing.T) {
	tests := []struct {
		name string
		want *SpamTracker
	}{
		{
			name: "create new spam tracker struct",
			want: &SpamTracker{
				max:      0,
				messages: make(map[string][]string),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSpamTracker(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSpamTracker() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpamTracker_Handler(t *testing.T) {
	err := os.Setenv("BOT_OWNER", "11111")
	if err != nil {
		t.FailNow()
	}

	mockMessages := make(map[string][]string)

	type fields struct {
		max      int
		messages map[string][]string
	}
	type args struct {
		request []string
		msg     *dg.MessageCreate
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Request without parameters",
			fields: fields{
				max:      5,
				messages: mockMessages,
			},
			args: args{
				request: []string{"cmd"},
				msg:     &dg.MessageCreate{},
			},
			want:    "Usage: <prefix>antispam set [number]",
			wantErr: false,
		},
		{
			name: "Request parameter is not 'set'",
			fields: fields{
				max:      5,
				messages: mockMessages,
			},
			args: args{
				request: []string{"cmd", "something"},
				msg:     &dg.MessageCreate{},
			},
			want:    "Usage: <prefix>antispam set [number]",
			wantErr: false,
		},
		{
			name: "Not bot owner executing 'set'",
			fields: fields{
				max:      5,
				messages: mockMessages,
			},
			args: args{
				request: []string{"cmd", "set"},
				msg: &dg.MessageCreate{
					&dg.Message{
						Author: &dg.User{
							ID: "00000",
						},
					},
				},
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "set max",
			fields: fields{
				messages: mockMessages,
			},
			args: args{
				request: []string{"cmd", "set", "5"},
				msg: &dg.MessageCreate{
					&dg.Message{
						Author: &dg.User{
							ID: "11111",
						},
					},
				},
			},
			want:    "Max repeated messages for kick set.",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SpamTracker{
				max:      tt.fields.max,
				messages: tt.fields.messages,
			}
			got, err := s.Handler(tt.args.request, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("SpamTracker.Handler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SpamTracker.Handler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpamTracker_CheckForSpam(t *testing.T) {
	err := os.Setenv("BOT_OWNER", "11111")
	if err != nil {
		t.FailNow()
	}

	mockSpam := make(map[string][]string)
	mockSpam["22222"] = []string{"test", "test", "test", "test"}

	type fields struct {
		max      int
		messages map[string][]string
	}
	type args struct {
		msg *dg.MessageCreate
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "Spam!",
			fields: fields{
				max:      5,
				messages: mockSpam,
			},
			args: args{
				msg: &dg.MessageCreate{
					&dg.Message{
						Content: "test",
						Author: &dg.User{
							ID: "22222",
						},
					},
				},
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "Not Spam!",
			fields: fields{
				max:      5,
				messages: mockSpam,
			},
			args: args{
				msg: &dg.MessageCreate{
					&dg.Message{
						Content: "No Spam!",
						Author: &dg.User{
							ID: "22222",
						},
					},
				},
			},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SpamTracker{
				max:      tt.fields.max,
				messages: tt.fields.messages,
			}
			got, err := s.CheckForSpam(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("SpamTracker.CheckForSpam() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SpamTracker.CheckForSpam() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSpamTracker_LoadMax(t *testing.T) {
	type fields struct {
		max      int
		messages map[string][]string
	}
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &SpamTracker{
				max:      tt.fields.max,
				messages: tt.fields.messages,
			}
			if err := s.LoadMax(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("SpamTracker.LoadMax() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
