package roles

import (
	"fmt"
	"os"
	"reflect"
	"testing"

	dg "github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
	"gitlab.com/technonauts/akordo/roles/mocks"
	"gitlab.com/technonauts/akordo/xp"
)

func TestNewRoleStorage(t *testing.T) {
	dgSess := &dg.Session{}
	xpSys := &xp.System{}
	type args struct {
		s  *dg.Session
		xp *xp.System
	}
	tests := []struct {
		name string
		args args
		want Assigner
	}{
		{
			name: "Default role storage creation",
			args: args{s: dgSess, xp: xpSys},
			want: &roleSystem{
				dgs:   dgSess,
				xp:    xpSys,
				tiers: &autoRanks{Tiers: make(map[string]float64)},
				sar: &roleStorage{
					SelfRoles: make(map[string]string),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRoleStorage(tt.args.s, tt.args.xp); !reflect.DeepEqual(got, tt.want) {
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

	testRoleMap2 := make(map[string]float64)
	testRoleMap2["Golang"] = 12345.00
	testRoleMap2["Rust"] = 54321.00

	var sortedList string
	sortedList = fmt.Sprintf("%s\n%s: %.2f", sortedList, "Golang", float64(12345))
	sortedList = fmt.Sprintf("%s\n%s: %.2f", sortedList, "Rust", float64(54321))

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
		dgs   DgSession
		sar   *roleStorage
		tiers *autoRanks
		xp    *xp.System
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
		{
			name: "lar: List ranks in sorted order",
			fields: fields{
				dgs: testDGS,
				sar: &roleStorage{
					SelfRoles: testRoleMap,
				},
				tiers: &autoRanks{Tiers: testRoleMap2},
			},
			args: args{
				req: []string{"cmd", "lar"},
				msg: &dg.MessageCreate{
					&dg.Message{GuildID: "0", Author: &dg.User{Username: "User1", ID: "0000"}},
				},
			},
			want:    &dg.MessageEmbed{Title: "Auto Rank XP", Description: sortedList},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &roleSystem{
				dgs:   tt.fields.dgs,
				sar:   tt.fields.sar,
				tiers: tt.fields.tiers,
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

func TestSystem_AutoPromote(t *testing.T) {
	testTiers := make(map[string]float64)
	testTiers["Crew"] = 100
	testTiers["Ensign"] = 200

	testUsers := make(map[string]float64)
	testUsers["1111"] = 100
	testUsers["2222"] = 10

	role1 := &dg.Role{ID: "0000", Name: "Crew"}
	role2 := &dg.Role{ID: "1111", Name: "Ensign"}
	roleSlice := []*dg.Role{role1, role2}

	mockSess := new(mocks.DgSession)
	// For "user found and promoted"
	mockSess.On("GuildRoles", "1111").Return(roleSlice, nil).Once()
	mockSess.On("GuildMemberRoleAdd", "1111", "1111", "0000").Return(nil).Once()

	// For "bad guild id"
	mockSess.On("GuildRoles", "0000").Return(nil, fmt.Errorf("some fake error")).Once()

	// For "user not found"
	mockSess.On("GuildRoles", "1111").Return(roleSlice, nil).Once()

	// For "fail to update user roles"
	mockSess.On("GuildRoles", "1111").Return(roleSlice, nil).Once()
	mockSess.On("GuildMemberRoleAdd", "1111", "4444", "").
		Return(fmt.Errorf("some fake error")).Once()

	type fields struct {
		dgs   DgSession
		xp    *xp.System
		sar   *roleStorage
		tiers *autoRanks
	}
	type args struct {
		msg *dg.MessageCreate
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "User found and promoted",
			fields: fields{
				xp: &xp.System{
					Data: &xp.Data{Users: testUsers},
				},
				tiers: &autoRanks{Tiers: testTiers},
				dgs:   mockSess,
			},
			args: args{
				msg: &dg.MessageCreate{
					&dg.Message{
						GuildID: "1111",
						Author: &dg.User{
							ID: "1111",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "Bad Guild ID",
			fields: fields{
				xp: &xp.System{
					Data: &xp.Data{Users: testUsers},
				},
				tiers: &autoRanks{Tiers: testTiers},
				dgs:   mockSess,
			},
			args: args{
				msg: &dg.MessageCreate{
					&dg.Message{
						GuildID: "0000",
						Author: &dg.User{
							ID: "2222",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "User not found",
			fields: fields{
				xp: &xp.System{
					Data: &xp.Data{Users: testUsers},
				},
				tiers: &autoRanks{Tiers: testTiers},
				dgs:   mockSess,
			},
			args: args{
				msg: &dg.MessageCreate{
					&dg.Message{
						GuildID: "1111",
						Author: &dg.User{
							ID: "3333",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "Fail to update user roles",
			fields: fields{
				xp: &xp.System{
					Data: &xp.Data{Users: testUsers},
				},
				tiers: &autoRanks{Tiers: testTiers},
				dgs:   mockSess,
			},
			args: args{
				msg: &dg.MessageCreate{
					&dg.Message{
						GuildID: "1111",
						Author: &dg.User{
							ID: "4444",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &roleSystem{
				xp:    tt.fields.xp,
				tiers: tt.fields.tiers,
				dgs:   tt.fields.dgs,
			}
			if err := r.AutoPromote(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("System.AutoPromote() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
