package plugins

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	dg "github.com/bwmarrin/discordgo"
	"gitlab.com/technonauts/akordo/plugins/mocks"
)

func TestNewAgreement(t *testing.T) {
	dgSess := &dg.Session{}

	type args struct {
		s *dg.Session
	}
	tests := []struct {
		name string
		args args
		want *Agreement
	}{
		{
			name: "create new agreement struct",
			args: args{s: dgSess},
			want: &Agreement{
				session:  dgSess,
				BaseRole: "",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAgreement(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAgreement() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAgreement_Handler(t *testing.T) {
	err := os.Setenv("BOT_OWNER", "1111")
	if err != nil {
		t.FailNow()
	}
	mockSession := new(mocks.AkSession)
	mockSession.On("GuildMemberRoleAdd", "", "22222", "12345").Return(nil).Once()
	mockSession.On("GuildMemberRoleAdd", "", "22222", "12345").Return(fmt.Errorf("some message")).Once()

	type fields struct {
		session  AkSession
		BaseRole string
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
			name: "Comand call with a length less than 2",
			fields: fields{
				session:  mockSession,
				BaseRole: "",
			},
			args: args{
				req: []string{"cmd"},
				msg: &dg.MessageCreate{},
			},
			want:    "Usage: <prefix>rules agreed",
			wantErr: false,
		},
		{
			name: "Command run with unrecognized parameter",
			fields: fields{
				session:  mockSession,
				BaseRole: "",
			},
			args: args{
				req: []string{"cmd", "boot"},
				msg: &dg.MessageCreate{},
			},
			want:    "Usage: <prefix>rules agreed",
			wantErr: false,
		},
		{
			name: "rules set command",
			fields: fields{
				session:  mockSession,
				BaseRole: "",
			},
			args: args{
				req: []string{"cmd", "set", "<@&12345>"},
				msg: &dg.MessageCreate{
					&dg.Message{
						Author: &dg.User{
							ID: "1111",
						},
					},
				},
			},
			want:    "Role added as as the base chat role",
			wantErr: false,
		},
		{
			name: "set command run by non-owner",
			fields: fields{
				session:  mockSession,
				BaseRole: "",
			},
			args: args{
				req: []string{"cmd", "set", "<@&12345>"},
				msg: &dg.MessageCreate{
					&dg.Message{
						Author: &dg.User{
							ID: "0000",
						},
					},
				},
			},
			want:    "This command is for the bot owner only :rage:",
			wantErr: false,
		},
		{
			name: "set command using non-@ role name",
			fields: fields{
				session:  mockSession,
				BaseRole: "",
			},
			args: args{
				req: []string{"cmd", "set", "someRoleName"},
				msg: &dg.MessageCreate{
					&dg.Message{
						Author: &dg.User{
							ID: "1111",
						},
					},
				},
			},
			want:    "Editor must be a role and formatted as: `@mod`",
			wantErr: false,
		},
		{
			name: "User self assigns base role",
			fields: fields{
				session:  mockSession,
				BaseRole: "12345",
			},
			args: args{
				req: []string{"cmd", "agreed"},
				msg: &dg.MessageCreate{
					&dg.Message{
						Author: &dg.User{
							ID: "22222",
						},
					},
				},
			},
			want:    "Added :ok_hand:",
			wantErr: false,
		},
		{
			name: "GuildMemberRoleAdd fails",
			fields: fields{
				session:  mockSession,
				BaseRole: "12345",
			},
			args: args{
				req: []string{"cmd", "agreed"},
				msg: &dg.MessageCreate{
					&dg.Message{
						Author: &dg.User{
							ID: "22222",
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
			a := &Agreement{
				session:  tt.fields.session,
				BaseRole: tt.fields.BaseRole,
			}
			got, err := a.Handler(tt.args.req, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Agreement.Handler() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Agreement.Handler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAgreement_LoadAgreementRole(t *testing.T) {
	mockSession := new(mocks.AkSession)
	type fields struct {
		session  AkSession
		BaseRole string
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
		{
			name: "load saved data",
			fields: fields{
				session:  mockSession,
				BaseRole: "12345",
			},
			args:    args{"data/ruleRole.json"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Agreement{
				session:  tt.fields.session,
				BaseRole: tt.fields.BaseRole,
			}
			if err := a.LoadAgreementRole(tt.args.filePath); (err != nil) != tt.wantErr {
				t.Errorf("Agreement.LoadAgreementRole() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
