package xp

import (
	"fmt"
	"reflect"
	"sync"
	"testing"

	p "git.sr.ht/~jrswab/akordo/plugins"
	"git.sr.ht/~jrswab/akordo/xp/mocks"
	dg "github.com/bwmarrin/discordgo"
)

func TestNewXpStore(t *testing.T) {
	testMutex := &sync.Mutex{}
	testSession := &dg.Session{}

	type args struct {
		mtx *sync.Mutex
		dgs *dg.Session
	}
	tests := []struct {
		name string
		args args
		want Exp
	}{
		{
			name: "Create default xpStore",
			args: args{testMutex, testSession},
			want: &System{
				data:    &xpData{Users: make(map[string]float64)},
				mutex:   testMutex,
				dgs:     testSession,
				callRec: p.NewRecorder(),
				tiers:   &autoRanks{Tiers: make(map[string]float64)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewXpStore(tt.args.mtx, tt.args.dgs); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewXpStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSystem_LoadXP(t *testing.T) {
	type fields struct {
		data  *xpData
		mutex *sync.Mutex
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
			name: "file exists",
			fields: fields{
				data:  &xpData{},
				mutex: &sync.Mutex{},
			},
			args: args{
				file: "testXP.json",
			},
			wantErr: false,
		},
		{
			name: "file does not exist",
			fields: fields{
				data:  &xpData{},
				mutex: &sync.Mutex{},
			},
			args: args{
				file: "testMissingXp.json",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &System{
				data:  tt.fields.data,
				mutex: tt.fields.mutex,
			}
			if err := x.LoadXP(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("System.LoadXP() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSystem_ManipulateXP(t *testing.T) {
	data := &xpData{Users: make(map[string]float64)}
	data.Users["2222"] = 0.10

	type fields struct {
		data  *xpData
		mutex *sync.Mutex
	}
	type args struct {
		action string
		msg    *dg.MessageCreate
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
	}{
		{
			name: "give xp to new user",
			fields: fields{
				data:  data,
				mutex: &sync.Mutex{},
			},
			args: args{
				action: "addMessagePoints",
				msg: &dg.MessageCreate{
					&dg.Message{
						Content: "0123456789",
						Author: &dg.User{
							ID: "1111",
						},
					},
				},
			},
			want: 0.10,
		},
		{
			name: "give xp to existing user",
			fields: fields{
				data:  data,
				mutex: &sync.Mutex{},
			},
			args: args{
				action: "addMessagePoints",
				msg: &dg.MessageCreate{
					&dg.Message{
						Content: "0123456789",
						Author: &dg.User{
							ID: "2222",
						},
					},
				},
			},
			want: 0.20,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &System{
				data:  tt.fields.data,
				mutex: tt.fields.mutex,
			}
			x.ManipulateXP(tt.args.action, tt.args.msg)
			if data.Users[tt.args.msg.Author.ID] != tt.want {
				t.Errorf("got %.2f, want %.2f", data.Users[tt.args.msg.Author.ID], tt.want)
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
	testUsers["4444"] = 50

	role1 := &dg.Role{ID: "0000", Name: "Crew"}
	role2 := &dg.Role{ID: "1111", Name: "Ensign"}
	roleSlice := []*dg.Role{role1, role2}

	mockSess := new(mocks.AkSession)
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
		data        *xpData
		tiers       *autoRanks
		defaultFile string
		callRec     *p.Record
		mutex       *sync.Mutex
		dgs         AkSession
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
				data:  &xpData{Users: testUsers},
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
				data:  &xpData{Users: testUsers},
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
				data:  &xpData{Users: testUsers},
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
				data:  &xpData{Users: testUsers},
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
			x := &System{
				data:    tt.fields.data,
				tiers:   tt.fields.tiers,
				callRec: tt.fields.callRec,
				mutex:   tt.fields.mutex,
				dgs:     tt.fields.dgs,
			}
			if err := x.AutoPromote(tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("System.AutoPromote() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
