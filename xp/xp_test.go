package xp

import (
	"reflect"
	"sync"
	"testing"

	p "git.sr.ht/~jrswab/akordo/plugins"
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
				data:        &xpData{Users: make(map[string]float64)},
				mutex:       testMutex,
				dgs:         testSession,
				callRec:     p.NewRecorder(),
				defaultFile: DefaultFile,
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
