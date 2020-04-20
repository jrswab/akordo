package plugins

import (
	"testing"

	man "git.sr.ht/~jrswab/akordo/manuals"
	dg "github.com/bwmarrin/discordgo"
)

func TestManual(t *testing.T) {
	type args struct {
		req []string
		s   *dg.Session
		msg *dg.MessageCreate
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Command missing",
			args: args{
				req: []string{"~testCMD"},
				s:   &dg.Session{},
				msg: &dg.MessageCreate{},
			},
			want: "Usage: `<prefix>man command`",
		},
		{
			name: "Return undefined command",
			args: args{
				req: []string{"~testCMD", "satesat"},
				s:   &dg.Session{},
				msg: &dg.MessageCreate{},
			},
			want: "",
		},
		{
			name: "Return gif",
			args: args{
				req: []string{"~testCMD", "gif"},
				s:   &dg.Session{},
				msg: &dg.MessageCreate{},
			},
			want: man.Gif,
		},
		{
			name: "Return man",
			args: args{
				req: []string{"~testCMD", "man"},
				s:   &dg.Session{},
				msg: &dg.MessageCreate{},
			},
			want: man.Man,
		},
		{
			name: "Return meme",
			args: args{
				req: []string{"~testCMD", "meme"},
				s:   &dg.Session{},
				msg: &dg.MessageCreate{},
			},
			want: man.Meme,
		},
		{
			name: "Return ping",
			args: args{
				req: []string{"~testCMD", "ping"},
				s:   &dg.Session{},
				msg: &dg.MessageCreate{},
			},
			want: man.Ping,
		},
		{
			name: "Return rule34",
			args: args{
				req: []string{"~testCMD", "rule34"},
				s:   &dg.Session{},
				msg: &dg.MessageCreate{},
			},
			want: man.Rule34,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Manual(tt.args.req, tt.args.s, tt.args.msg)
			if got != tt.want {
				t.Errorf("Manual() = %v, want %v", got, tt.want)
			}
		})
	}
}
