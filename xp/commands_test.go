package xp

import (
	"fmt"
	"strconv"
	"sync"
	"testing"

	p "git.sr.ht/~jrswab/akordo/plugins"
	"git.sr.ht/~jrswab/akordo/xp/mocks"
	dg "github.com/bwmarrin/discordgo"
)

func TestSystem_Execute(t *testing.T) {
	testUsers := make(map[string]float64)
	testUsers["165899680323076097"] = 5.56

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

	mockSess := new(mocks.AkSession)
	mockSess.On("GuildMember", "1111", "165899680323076097").Return(mockMemb, nil).Once()
	mockSess.On("GuildMembers", "1111", "", 1000).Return(mockMembSlice, nil)
	mockSess.On("GuildMembers", "1111", "998", 1000).Return(mockMembSlice2, nil)

	type fields struct {
		data    *xpData
		callRec *p.Record
		mutex   *sync.Mutex
		dgs     AkSession
	}
	type args struct {
		req []string
		msg *dg.MessageCreate
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    MsgEmbed
		wantErr bool
	}{
		{
			name: "save comand",
			fields: fields{
				data:    &xpData{},
				callRec: p.NewRecorder(),
				mutex:   &sync.Mutex{},
				dgs:     &dg.Session{},
			},
			args: args{
				req: []string{"=xp", "save"},
				msg: &dg.MessageCreate{},
			},
			want:    &dg.MessageEmbed{Description: "XP data saved!"},
			wantErr: false,
		},
		{
			name: "xp without params",
			fields: fields{
				data:    &xpData{Users: testUsers},
				callRec: p.NewRecorder(),
				mutex:   &sync.Mutex{},
				dgs:     &dg.Session{},
			},
			args: args{
				req: []string{"=xp"},
				msg: &dg.MessageCreate{
					&dg.Message{
						Author: &dg.User{
							ID:       "165899680323076097",
							Username: "jrswab",
						},
					},
				},
			},
			want:    &dg.MessageEmbed{Description: fmt.Sprintf("%s has a total of %.2f xp", "jrswab", 5.56)},
			wantErr: false,
		},
		{
			name: "xp with @username",
			fields: fields{
				data:    &xpData{Users: testUsers},
				callRec: p.NewRecorder(),
				mutex:   &sync.Mutex{},
				dgs:     mockSess,
			},
			args: args{
				req: []string{"=xp", "<@!165899680323076097>"},
				msg: &dg.MessageCreate{
					&dg.Message{
						GuildID: "1111",
						Author: &dg.User{
							ID:       "1",
							Username: "AnotherUser",
						},
					},
				},
			},
			want:    &dg.MessageEmbed{Description: fmt.Sprintf("%s has a total of %.2f xp", "jrswab", 5.56)},
			wantErr: false,
		},
		{
			name: "xp with username and no @",
			fields: fields{
				data:    &xpData{Users: testUsers},
				callRec: p.NewRecorder(),
				mutex:   &sync.Mutex{},
				dgs:     mockSess,
			},
			args: args{
				req: []string{"=xp", "jrswab"},
				msg: &dg.MessageCreate{
					&dg.Message{
						GuildID: "1111",
						Author: &dg.User{
							ID:       "165899680323076097",
							Username: "jrswab",
						},
					},
				},
			},
			want:    &dg.MessageEmbed{Description: fmt.Sprintf("%s has a total of %.2f xp", "jrswab", 5.56)},
			wantErr: false,
		},
		{
			name: "user not found in xp.json file",
			fields: fields{
				data:    &xpData{Users: testUsers},
				callRec: p.NewRecorder(),
				mutex:   &sync.Mutex{},
				dgs:     mockSess,
			},
			args: args{
				req: []string{"=xp", "User3"},
				msg: &dg.MessageCreate{
					&dg.Message{
						GuildID: "1111",
						Author: &dg.User{
							ID:       "3",
							Username: "User3",
						},
					},
				},
			},
			want:    &dg.MessageEmbed{Description: "User3 has not earned any XP"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x := &System{
				data:    tt.fields.data,
				callRec: tt.fields.callRec,
				mutex:   tt.fields.mutex,
				dgs:     tt.fields.dgs,
			}
			got, err := x.Execute(tt.args.req, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("System.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("System.Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
