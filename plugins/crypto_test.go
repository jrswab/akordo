package plugins

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	dg "github.com/bwmarrin/discordgo"
)

func TestNewCrypto(t *testing.T) {
	tests := []struct {
		name string
		want *Crypto
	}{
		{
			name: "Create blank struct",
			want: &Crypto{},
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
	type fields struct {
		words        []byte
		encoded      string
		lastEncoding string
		inPlay       bool
		roundStart   time.Time
		roundEnd     time.Time
		wasDecoded   bool
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
			name:   "catch length less than 2",
			fields: fields{},
			args: args{
				req: []string{"~test"},
				msg: &dg.MessageCreate{},
			},
			want:    "Usage: `<prefix>crypto init` to start a game",
			wantErr: false,
		},
		{
			name: "catch init during gameplay",
			fields: fields{
				inPlay:  true,
				encoded: "1111",
			},
			args: args{
				req: []string{"~test", "init"},
				msg: &dg.MessageCreate{},
			},
			want:    fmt.Sprintf("Game in progress. Current encoding:\n1111"),
			wantErr: false,
		},
		{
			name: "check user input concatenation",
			fields: fields{
				inPlay: true,
				words:  []byte("this is a test"),
			},
			args: args{
				req: []string{"~test", "this", "is", "a", "test"},
				msg: &dg.MessageCreate{&dg.Message{Author: &dg.User{Username: "User1"}}},
			},
			want:    "User1 won this round! Will you be next?",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Crypto{
				words:        tt.fields.words,
				encoded:      tt.fields.encoded,
				lastEncoding: tt.fields.lastEncoding,
				inPlay:       tt.fields.inPlay,
				roundStart:   tt.fields.roundStart,
				roundEnd:     tt.fields.roundEnd,
				wasDecoded:   tt.fields.wasDecoded,
			}
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
