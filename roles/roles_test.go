package roles

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	"git.sr.ht/~jrswab/akordo/roles/mocks"
	dg "github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
)

func TestNewRoleStorage(t *testing.T) {
	dgSess := &dg.Session{}
	type args struct {
		s *dg.Session
	}
	tests := []struct {
		name string
		args args
		want Assigner
	}{
		{
			name: "Default role storage creation",
			args: args{s: dgSess},
			want: &roleSystem{
				dgs: dgSess,
				sar: &roleStorage{
					SelfRoles: make(map[string]string),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRoleStorage(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRoleStorage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_roleSystem_LoadSelfAssignRoles(t *testing.T) {
	type fields struct {
		dgs *dg.Session
		sar *roleStorage
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
		{
			name:    "Load test json file",
			fields:  fields{dgs: &dg.Session{}, sar: &roleStorage{}},
			args:    args{file: "test.json"},
			wantErr: false,
		},
		{
			name:    "Incorrect file path",
			fields:  fields{dgs: &dg.Session{}, sar: &roleStorage{}},
			args:    args{file: "wrongTest.json"},
			wantErr: true,
		},
		{
			name:    "Invalid JSON file",
			fields:  fields{dgs: &dg.Session{}, sar: &roleStorage{}},
			args:    args{file: "testInvalid.json"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &roleSystem{
				dgs: tt.fields.dgs,
				sar: tt.fields.sar,
			}
			if err := r.LoadSelfAssignRoles(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("roleSystem.LoadSelfAssignRoles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_roleSystem_ExecuteRoleCommands(t *testing.T) {
	testRoleMap := make(map[string]string)
	testRoleMap["Golang"] = "12345"
	var testRoleOutput string

	for role := range testRoleMap {
		testRoleOutput = fmt.Sprintf("%s\n%s", testRoleOutput, role)
	}

	notListed := "I don't know what to do :thinking:"
	notListed = fmt.Sprintf("%s\nPlease check the command and try again", notListed)

	err := os.Setenv("BOT_OWNER", "1111")
	if err != nil {
		t.Fatalf("Exporting test envar failed")
	}

	testDGS := new(mocks.DgSession)
	roleSlice := []*dg.Role{&dg.Role{ID: "12345", Name: "testRole"}}
	testDGS.On("GuildRoles", mock.Anything).Return(roleSlice, nil).Twice()
	testDGS.On("GuildMemberRoleAdd", "0", "0000", "12345").Return(nil).Once()
	testDGS.On("GuildMemberRoleRemove", "0", "0000", "12345").Return(nil).Once()

	type fields struct {
		dgs DgSession
		sar *roleStorage
	}
	type args struct {
		req []string
		msg *dg.MessageCreate
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *dg.MessageEmbed
		wantErr bool
	}{
		{
			name: "asar: Add new self assign role",
			fields: fields{
				dgs: testDGS,
				sar: &roleStorage{
					SelfRoles: make(map[string]string),
				},
			},
			args: args{
				req: []string{"cmd", "asar", "testRole"},
				msg: &dg.MessageCreate{
					&dg.Message{Author: &dg.User{ID: "1111"}},
				},
			},
			want:    &dg.MessageEmbed{Description: "Added testRole to the self assign role list"},
			wantErr: false,
		},
		{
			name: "asar: Executor not the bot owner",
			fields: fields{
				dgs: testDGS,
				sar: &roleStorage{
					SelfRoles: make(map[string]string),
				},
			},
			args: args{
				req: []string{"cmd", "asar", "testRole2"},
				msg: &dg.MessageCreate{
					&dg.Message{Author: &dg.User{ID: "0000"}},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "asar: Command length less than 3",
			fields: fields{
				dgs: testDGS,
				sar: &roleStorage{
					SelfRoles: make(map[string]string),
				},
			},
			args: args{
				req: []string{"cmd", "asar"},
				msg: &dg.MessageCreate{
					&dg.Message{Author: &dg.User{ID: "1111"}},
				},
			},
			want:    &dg.MessageEmbed{Description: "Usage: `<prefix>roles asar [role name]`"},
			wantErr: false,
		},
		{
			name: "asar: Role name with spaces",
			fields: fields{
				dgs: testDGS,
				sar: &roleStorage{
					SelfRoles: make(map[string]string),
				},
			},
			args: args{
				req: []string{"cmd", "asar", "test", "role"},
				msg: &dg.MessageCreate{
					&dg.Message{Author: &dg.User{ID: "1111"}},
				},
			},
			want:    &dg.MessageEmbed{Description: "Added test role to the self assign role list"},
			wantErr: false,
		},
		{
			name: "lsar: List roles",
			fields: fields{
				dgs: testDGS,
				sar: &roleStorage{
					SelfRoles: testRoleMap,
				},
			},
			args: args{
				req: []string{"cmd", "lsar"},
				msg: &dg.MessageCreate{
					&dg.Message{Author: &dg.User{ID: "0000"}},
				},
			},
			want:    &dg.MessageEmbed{Description: testRoleOutput},
			wantErr: false,
		},
		{
			name: "sar: Add role to user",
			fields: fields{
				dgs: testDGS,
				sar: &roleStorage{
					SelfRoles: testRoleMap,
				},
			},
			args: args{
				req: []string{"cmd", "sar", "Golang"},
				msg: &dg.MessageCreate{
					&dg.Message{GuildID: "0", Author: &dg.User{Username: "User1", ID: "0000"}},
				},
			},
			want:    &dg.MessageEmbed{Description: "Added role, Golang, to User1"},
			wantErr: false,
		},
		{
			name: "usar: Remove role to user",
			fields: fields{
				dgs: testDGS,
				sar: &roleStorage{
					SelfRoles: testRoleMap,
				},
			},
			args: args{
				req: []string{"cmd", "uar", "Golang"},
				msg: &dg.MessageCreate{
					&dg.Message{GuildID: "0", Author: &dg.User{Username: "User1", ID: "0000"}},
				},
			},
			want:    &dg.MessageEmbed{Description: "Removed role, Golang, from User1"},
			wantErr: false,
		},
		{
			name: "Incorrect parameter passed with command",
			fields: fields{
				dgs: testDGS,
				sar: &roleStorage{
					SelfRoles: testRoleMap,
				},
			},
			args: args{
				req: []string{"cmd", "blah", "Golang"},
				msg: &dg.MessageCreate{
					&dg.Message{GuildID: "0", Author: &dg.User{Username: "User1", ID: "0000"}},
				},
			},
			want:    &dg.MessageEmbed{Description: notListed},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &roleSystem{
				dgs: tt.fields.dgs,
				sar: tt.fields.sar,
			}
			got, err := r.ExecuteRoleCommands(tt.args.req, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("roleSystem.ExecuteRoleCommands() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("roleSystem.ExecuteRoleCommands() = %v, want %v", got, tt.want)
			}
		})
	}
}
