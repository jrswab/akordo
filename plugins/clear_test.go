package plugins

import (
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"

	dg "github.com/bwmarrin/discordgo"
	"gitlab.com/technonauts/akordo/plugins/mocks"
)

func TestNewEraser(t *testing.T) {
	mockSession := new(mocks.AkSession)
	type args struct {
		s AkSession
	}
	tests := []struct {
		name string
		args args
		want Eraser
	}{
		{
			name: "Create new struct",
			args: args{s: mockSession},
			want: &clear{
				dgs: mockSession,
				Authd: &authorizedRoles{
					roleID: make(map[string]bool),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEraser(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEraser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_clear_ClearHandler(t *testing.T) {
	err := os.Setenv("BOT_ID", "0000")
	if err != nil {
		t.FailNow()
	}
	err = os.Setenv("BOT_OWNER", "1111")
	if err != nil {
		t.FailNow()
	}

	mockSession := new(mocks.AkSession)
	msg1 := &dg.Message{ChannelID: "1111", Author: &dg.User{ID: "0000"}}
	msgSlice := []*dg.Message{msg1}
	msgIDSlice := []string{}

	mockMemb := &dg.Member{
		User: &dg.User{
			ID:       "165899680323076097",
			Username: "jrswab",
		},
	}
	mockMembSlice := []*dg.Member{}
	// create a []*dg.Member with 1000 "members"
	for i := 0; i < 999; i++ {
		if i == 0 {
			mockMembSlice = append(mockMembSlice, mockMemb)
		}

		mockMembInc := &dg.Member{
			User: &dg.User{
				ID:       strconv.Itoa(i),
				Username: fmt.Sprintf("User%d", i),
			},
		}

		mockMembSlice = append(mockMembSlice, mockMembInc)
	}

	mockMembSlice2 := []*dg.Member{
		{
			User: &dg.User{
				ID:       "1000",
				Username: "User1000",
			},
		},
	}

	mockSession.On("ChannelMessages", "1111", 100, "", "", "").Return(msgSlice, nil)
	mockSession.On("ChannelMessagesBulkDelete", "1111", msgIDSlice).Return(nil)
	mockSession.On("GuildMembers", "", "", 1000).Return(mockMembSlice, nil)
	mockSession.On("GuildMembers", "", "998", 1000).Return(mockMembSlice2, nil)

	type fields struct {
		dgs   AkSession
		Authd *authorizedRoles
	}
	type args struct {
		request []string
		msg     *dg.MessageCreate
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Set authorized role for clearing messages",
			fields: fields{
				dgs:   mockSession,
				Authd: &authorizedRoles{roleID: make(map[string]bool)},
			},
			args: args{
				request: []string{"=clear", "set", "<@&123456879>"},
				msg: &dg.MessageCreate{
					&dg.Message{
						ChannelID: "1111",
						Author: &dg.User{
							ID: "1111",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "clear a user's messages",
			fields: fields{
				dgs:   mockSession,
				Authd: &authorizedRoles{},
			},
			args: args{
				request: []string{"=clear", "testUser"},
				msg: &dg.MessageCreate{
					&dg.Message{
						ChannelID: "1111",
						Author: &dg.User{
							ID: "0000",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "clear a @user messages",
			fields: fields{
				dgs:   mockSession,
				Authd: &authorizedRoles{},
			},
			args: args{
				request: []string{"=clear", "<@!testUser>"},
				msg: &dg.MessageCreate{
					&dg.Message{
						ChannelID: "1111",
						Author: &dg.User{
							ID: "0000",
						},
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &clear{
				dgs:   tt.fields.dgs,
				Authd: tt.fields.Authd,
			}
			if err := c.ClearHandler(tt.args.request, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("clear.ClearHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_clear_LoadAuthList(t *testing.T) {
	type fields struct {
		dgs   AkSession
		Authd *authorizedRoles
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
			c := &clear{
				dgs:   tt.fields.dgs,
				Authd: tt.fields.Authd,
			}
			if err := c.LoadAuthList(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("clear.LoadAuthList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
