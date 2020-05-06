package plugins

import (
	"os"
	"reflect"
	"testing"

	dg "github.com/bwmarrin/discordgo"
	"gitlab.com/technonauts/akordo/plugins/mocks"
)

func TestNewBlacklist(t *testing.T) {
	mockSession := new(mocks.AkSession)
	type args struct {
		s AkSession
	}
	tests := []struct {
		name string
		args args
		want *Blacklist
	}{
		{
			name: "New Blacklist struct",
			args: args{s: mockSession},
			want: &Blacklist{
				session: mockSession,
				data: &data{
					Editors: make(map[string]bool),
					Banned:  make(map[string]bool),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBlacklist(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBlacklist() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlacklist_Handler(t *testing.T) {
	mockSession := new(mocks.AkSession)
	mockData := &data{Editors: make(map[string]bool), Banned: make(map[string]bool)}
	mockData.Editors["11111"] = true
	mockData.Banned["badWord"] = true

	err := os.Setenv("BOT_OWNER", "00000")
	if err != nil {
		t.FailNow()
	}

	type fields struct {
		session AkSession
		data    *data
	}
	type args struct {
		req []string
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
			name:   "Add Blacklist Word",
			fields: fields{session: mockSession, data: mockData},
			args: args{req: []string{"cmd", "add", "test"},
				msg: &dg.MessageCreate{
					&dg.Message{Author: &dg.User{ID: "1111"}},
				},
			},
			want:    "Word(s) added to the blacklist",
			wantErr: false,
		},
		{
			name:   "Remove Blacklist Word",
			fields: fields{session: mockSession, data: mockData},
			args: args{req: []string{"cmd", "remove", "test"},
				msg: &dg.MessageCreate{
					&dg.Message{Author: &dg.User{ID: "1111"}},
				},
			},
			want:    "Word(s) removed from the blacklist",
			wantErr: false,
		},
		{
			name:   "Add Role As Editor (successful)",
			fields: fields{session: mockSession, data: mockData},
			args: args{req: []string{"cmd", "editor", "<@&12345>"},
				msg: &dg.MessageCreate{
					&dg.Message{Author: &dg.User{ID: "00000"}},
				},
			},
			want:    "Role added as an editor",
			wantErr: false,
		},
		{
			name:   "Add Role As Editor (not Owner ID)",
			fields: fields{session: mockSession, data: mockData},
			args: args{req: []string{"cmd", "editor", "<@&12345>"},
				msg: &dg.MessageCreate{
					&dg.Message{Author: &dg.User{ID: "22222"}},
				},
			},
			want:    "This command is for the bot owner only :rage:",
			wantErr: false,
		},
		{
			name:   "Add Role As Editor (incorrect formatting)",
			fields: fields{session: mockSession, data: mockData},
			args: args{req: []string{"cmd", "editor", "someRole"},
				msg: &dg.MessageCreate{
					&dg.Message{Author: &dg.User{ID: "00000"}},
				},
			},
			want:    "Editor must be a role and formatted as: `@mod`",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Blacklist{
				session: tt.fields.session,
				data:    tt.fields.data,
			}
			got, err := b.Handler(tt.args.req, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Blacklist.Handler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Blacklist.Handler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlacklist_CheckBannedWords(t *testing.T) {
	mockSession := new(mocks.AkSession)
	mockData := &data{Editors: make(map[string]bool), Banned: make(map[string]bool)}
	mockData.Editors["11111"] = true
	mockData.Banned["badWord"] = true

	sfwChan := &dg.Channel{ID: "00000", NSFW: false}
	mockSession.On("Channel", "00000").Return(sfwChan, nil)

	type fields struct {
		session AkSession
		data    *data
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
			name: "Word not in banned word list",
			fields: fields{
				session: mockSession,
				data:    mockData,
			},
			args: args{
				msg: &dg.MessageCreate{
					&dg.Message{
						ChannelID: "00000",
						Content:   "goodWord",
						Author:    &dg.User{ID: "1111"},
						Member:    &dg.Member{Roles: []string{"00000", "22222"}},
					},
				},
			},
			want:    false,
			wantErr: false,
		},
		{
			name: "Word is in banned word list",
			fields: fields{
				session: mockSession,
				data:    mockData,
			},
			args: args{
				msg: &dg.MessageCreate{
					&dg.Message{
						ChannelID: "00000",
						Content:   "badWord",
						Author:    &dg.User{ID: "1111"},
						Member:    &dg.Member{Roles: []string{"00000", "22222"}},
					},
				},
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Blacklist{
				session: tt.fields.session,
				data:    tt.fields.data,
			}
			got, err := b.CheckBannedWords(tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Blacklist.CheckBannedWords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Blacklist.CheckBannedWords() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBlacklist_LoadBannedWordList(t *testing.T) {
	type fields struct {
		session AkSession
		data    *data
	}
	type args struct {
		filePath string
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
			b := &Blacklist{
				session: tt.fields.session,
				data:    tt.fields.data,
			}
			if err := b.LoadBannedWordList(tt.args.filePath); (err != nil) != tt.wantErr {
				t.Errorf("Blacklist.LoadBannedWordList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
