// Code generated by mockery v1.0.0. DO NOT EDIT.

package mocks

import (
	discordgo "github.com/bwmarrin/discordgo"
	mock "github.com/stretchr/testify/mock"
)

// AkSession is an autogenerated mock type for the AkSession type
type AkSession struct {
	mock.Mock
}

// GuildMember provides a mock function with given fields: guildID, userID
func (_m *AkSession) GuildMember(guildID string, userID string) (*discordgo.Member, error) {
	ret := _m.Called(guildID, userID)

	var r0 *discordgo.Member
	if rf, ok := ret.Get(0).(func(string, string) *discordgo.Member); ok {
		r0 = rf(guildID, userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*discordgo.Member)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(guildID, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GuildMemberRoleAdd provides a mock function with given fields: guildID, userID, roleID
func (_m *AkSession) GuildMemberRoleAdd(guildID string, userID string, roleID string) error {
	ret := _m.Called(guildID, userID, roleID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, string) error); ok {
		r0 = rf(guildID, userID, roleID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GuildMembers provides a mock function with given fields: guildID, after, limit
func (_m *AkSession) GuildMembers(guildID string, after string, limit int) ([]*discordgo.Member, error) {
	ret := _m.Called(guildID, after, limit)

	var r0 []*discordgo.Member
	if rf, ok := ret.Get(0).(func(string, string, int) []*discordgo.Member); ok {
		r0 = rf(guildID, after, limit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*discordgo.Member)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, int) error); ok {
		r1 = rf(guildID, after, limit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GuildRoles provides a mock function with given fields: guildID
func (_m *AkSession) GuildRoles(guildID string) ([]*discordgo.Role, error) {
	ret := _m.Called(guildID)

	var r0 []*discordgo.Role
	if rf, ok := ret.Get(0).(func(string) []*discordgo.Role); ok {
		r0 = rf(guildID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*discordgo.Role)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(guildID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}