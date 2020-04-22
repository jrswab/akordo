package plugins

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"git.sr.ht/~jrswab/akordo/xp"
	dg "github.com/bwmarrin/discordgo"
)

func TestNewCrypto(t *testing.T) {
	tests := []struct {
		name string
		want *Crypto
	}{
		{
			name: "Create blank struct",
			want: &Crypto{waitTime: 2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewCrypto(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCrypto() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCrypto_Game(t *testing.T) {
	type args struct {
		req []string
		msg *dg.MessageCreate
	}
	tests := []struct {
		name    string
		fields  *Crypto
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "catch length less than 2",
			fields: &Crypto{},
			args: args{
				req: []string{"=test"},
				msg: &dg.MessageCreate{},
			},
			want:    "Usage: `<prefix>crypto init` to start a game",
			wantErr: false,
		},
		{
			name: "catch init during gameplay",
			fields: &Crypto{
				inPlay:  true,
				encoded: "1111",
			},
			args: args{
				req: []string{"=test", "init"},
				msg: &dg.MessageCreate{},
			},
			want:    fmt.Sprintf("Mining in progress...\nCurrent encoding:\n1111"),
			wantErr: false,
		},
		{
			name: "correct guess",
			fields: &Crypto{
				XP:     xp.NewXpStore(&sync.Mutex{}),
				inPlay: true,
				words:  []byte("this is a test"),
			},
			args: args{
				req: []string{"=test", "this", "is", "a", "test"},
				msg: &dg.MessageCreate{&dg.Message{Author: &dg.User{Username: "User1"}}},
			},
			want:    "User1 won this round! Will you be next?",
			wantErr: false,
		},
		{
			name: "incorrect guess",
			fields: &Crypto{
				inPlay: true,
				words:  []byte("this is a test"),
			},
			args: args{
				req: []string{"=test", "this", "is", "a", "test2"},
				msg: &dg.MessageCreate{&dg.Message{Author: &dg.User{Username: "User1"}}},
			},
			want:    "User1 sorry, that is incorrect :smirk:",
			wantErr: false,
		},
		{
			name: "game cool-down period",
			fields: &Crypto{
				inPlay:   false,
				words:    []byte("this is a test"),
				waitTime: 10,
				roundEnd: time.Now(),
			},
			args: args{
				req: []string{"=test", "init"},
				msg: &dg.MessageCreate{&dg.Message{Author: &dg.User{Username: "User1"}}},
			},
			want:    fmt.Sprintf("Please wait 10 minutes to open a new mine."),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := tt.fields
			got, err := c.Game(tt.args.req, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Crypto.Game() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Crypto.Game() = %v, want %v", got, tt.want)
			}
		})
	}
}
