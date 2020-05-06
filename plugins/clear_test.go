package plugins

import (
	"reflect"
	"testing"

	dg "github.com/bwmarrin/discordgo"
	"gitlab.com/technonauts/akordo/plugins/mocks"
)

func TestNewEraser(t *testing.T) {
	mockSession := new(mocks.AkSession)
	type args struct {
		s AkSession
	}
	tests := []struct {
		name string
		args args
		want Eraser
	}{
		{
			name: "Create new struct",
			args: args{s: mockSession},
			want: &clear{
				dgs: mockSession,
				Authd: &authorizedRoles{
					roleID: make(map[string]bool),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewEraser(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewEraser() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_clear_ClearHandler(t *testing.T) {
	type fields struct {
		dgs   AkSession
		Authd *authorizedRoles
	}
	type args struct {
		request []string
		msg     *dg.MessageCreate
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "Clear bot messages",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &clear{
				dgs:   tt.fields.dgs,
				Authd: tt.fields.Authd,
			}
			if err := c.ClearHandler(tt.args.request, tt.args.msg); (err != nil) != tt.wantErr {
				t.Errorf("clear.ClearHandler() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_clear_LoadAuthList(t *testing.T) {
	type fields struct {
		dgs   AkSession
		Authd *authorizedRoles
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
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &clear{
				dgs:   tt.fields.dgs,
				Authd: tt.fields.Authd,
			}
			if err := c.LoadAuthList(tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("clear.LoadAuthList() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
