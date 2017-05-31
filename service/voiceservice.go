package service

import (
	"sync"

	"github.com/ynori7/go-irc/client"
)

const (
	MUTE_PREFIX   = "mute "
	UNMUTE_PREFIX = "unmute "
)

type VoiceService struct {
	connection client.Client

	mutedUsers      map[string]map[string]bool //channel => map[nick]: true
	mutedUsersMutex *sync.Mutex
}

func NewVoiceService(conn client.Client) VoiceService {
	return VoiceService{
		connection:      conn,
		mutedUsers:      make(map[string]map[string]bool),
		mutedUsersMutex: &sync.Mutex{},
	}
}

func (s *VoiceService) IsMuted(nick, location string) bool {
	s.mutedUsersMutex.Lock()
	defer s.mutedUsersMutex.Unlock()

	mutedUsers, ok := s.mutedUsers[location]
	if ok {
		_, ok := mutedUsers[nick]
		return ok
	}

	return false
}

func (s VoiceService) GiveVoice(nick, location string) {
	s.connection.SetMode(location, "+v", nick)
}

func (s *VoiceService) UnmuteUser(nick, location string) {
	s.mutedUsersMutex.Lock()
	defer s.mutedUsersMutex.Unlock()

	s.GiveVoice(nick, location)
	delete(s.mutedUsers[location], nick)
}

func (s *VoiceService) MuteUser(nick, location string) {
	s.mutedUsersMutex.Lock()
	defer s.mutedUsersMutex.Unlock()

	s.connection.SetMode(location, "-v", nick)

	_, ok := s.mutedUsers[location]
	if !ok {
		s.mutedUsers[location] = make(map[string]bool)
	}
	s.mutedUsers[location][nick] = true
}
