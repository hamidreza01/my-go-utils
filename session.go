package utils

import (
	"errors"
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
	SMap map[string]Session `json:"sessions"`
}

func (s *Session) UpdateUsingTime() {
	s.LastUse = time.Now().UnixNano()
}

func (S *Sessions) Clear(t int64) (c int) {
	now := time.Now().UnixNano()
	for id, session := range S.SMap {
		if session.LastUse+t < now {
			delete(S.SMap, id)
			c++
		}
	}
	return c
}

func (S *Sessions) Find(id string) (*Session, error) {
	session, ok := S.SMap[id]
	if !ok {
		return &Session{}, errors.New("session not found")
	}
	return &session, nil
}

func (S *Sessions) List(userID primitive.ObjectID) []Session {
	var userSessions []Session
	for _, session := range S.SMap {
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

	delete(S.SMap, sessionID)

	return nil
}
