package utils

import (
	"errors"
	"slices"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Session struct {
	User      primitive.ObjectID `json:"user_id"`
	Id        string             `json:"session"`
	LastUse   int64              `json:"last_use"`
	Ip        string             `json:"ip"`
	Timestamp int64              `json:"timestamp"`
	UserAgent string             `json:"user_agent"`
}

type Sessions struct {
	SArray []Session `json:"sessions"`
}

func (s *Session) UpdateUsingTime() {
	s.LastUse = time.Now().UnixNano()
}

func (S *Sessions) Clear(t int64) (c int) {
	now := time.Now().UnixNano()
	S.SArray = slices.DeleteFunc(S.SArray, func(s Session) bool {
		if s.LastUse+t < now {
			c++
			return true
		}
		return false
	})
	return c
}

func (S *Sessions) Find(id string) (*Session, error) {
	i := slices.IndexFunc(S.SArray, func(s Session) bool {
		return s.Id == id
	})
	if i == -1 {
		return &Session{}, errors.New("session not found")
	}
	return &S.SArray[i], nil
}

func (S *Sessions) List(userID primitive.ObjectID) []Session {
	var userSessions []Session
	for _, session := range S.SArray {
		if session.User == userID {
			userSessions = append(userSessions, session)
		}
	}
	return userSessions
}

func (S *Sessions) Kill(sessionID string) error {
	currentSession, err := S.Find(sessionID)
	if err != nil {
		return errors.New("session not found")
	}

	oneWeekAgo := time.Now().Add(-7 * 24 * time.Hour).UnixNano()
	if currentSession.Timestamp > oneWeekAgo {
		return errors.New("session must be at least one week old")
	}

	S.SArray = slices.DeleteFunc(S.SArray, func(s Session) bool {
		return s.User == currentSession.User && s.Id != sessionID
	})

	return nil
}
