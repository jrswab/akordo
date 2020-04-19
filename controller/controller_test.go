package controller

import (
	"log"
	"reflect"
	"testing"

	"git.sr.ht/~jrswab/akordo/plugins"
	"github.com/bwmarrin/discordgo"
	dg "github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/mock"
)

func TestNewSessionData(t *testing.T) {
	sess, err := dg.New("Bot fakeToken1111111")
	if err != nil {
		log.Fatalf("Session creation error: %s", err)
	}

	testRecord := plugins.NewRecorder()

	baseController := &SessionData{
		session:    sess,
		gifRecord:  testRecord,
		memeRecord: testRecord,
		pingRecord: testRecord,
		r34Record:  testRecord,
	}

	type args struct {
		s *dg.Session
	}
	tests := []struct {
		name string
		args args
		want *SessionData
	}{
		{
			name: "Returns proper struct data",
			args: args{s: sess},
			want: baseController,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSessionData(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewController() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Mocks for CheckSyntax:
type mockController struct {
	mock.Mock
}

func (_m *mockController) CheckSyntax(s *discordgo.Session, msg *discordgo.MessageCreate) {
	_m.Called(s, msg)
}

func (_m *mockController) ExecuteTask(req []string, s *discordgo.Session, msg *discordgo.MessageCreate) {
	_m.Called(req, s, msg)
}

func TestController_CheckSyntax(t *testing.T) {
	c := new(mockController)
	c.On("ExecuteTask", mock.Anything).Return().Times(2)

	sess, err := dg.New("Bot fakeToken1111111")
	if err != nil {
		log.Fatalf("Session creation error: %s", err)
	}

	correctPrefix := &dg.MessageCreate{
		&dg.Message{
			ID:        "000000",
			ChannelID: "000000",
			GuildID:   "000000",
			Content:   prefix + "blah",
		},
	}

	incorrectPrefix := &dg.MessageCreate{
		&dg.Message{
			ID:        "000000",
			ChannelID: "000000",
			GuildID:   "000000",
			Content:   "\\blah",
		},
	}

	type args struct {
		s   *dg.Session
		msg *dg.MessageCreate
	}
	tests := []struct {
		name string
		//fields fields
		args args
	}{
		{
			name: "Correct Prefix",
			args: args{
				s:   sess,
				msg: correctPrefix,
			},
		},
		{
			name: "Incorrect Prefix",
			args: args{
				s:   sess,
				msg: incorrectPrefix,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sd := NewSessionData(sess)
			sd.CheckSyntax(tt.args.s, tt.args.msg)
		})
	}
}
