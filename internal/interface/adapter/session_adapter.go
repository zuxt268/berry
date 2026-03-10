package adapter

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/zuxt268/berry/internal/config"
)

type sessionStore struct {
	store       *sessions.CookieStore
	sessionName string
}

// SessionAdapter セッションクッキー操作のインターフェース
//
//go:generate mockgen -source=$GOFILE -destination=./mock/mock_$GOFILE -package mock
type SessionAdapter interface {
	SaveSessionToken(r *http.Request, w http.ResponseWriter, token string) error
	GetSessionToken(r *http.Request) (string, bool, error)
	DeleteSessionToken(r *http.Request, w http.ResponseWriter) error
}

func NewSessionStore(sessionName string) SessionAdapter {
	store := sessions.NewCookieStore([]byte(config.Env.SessionSecret))
	store.Options = &sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7,
		HttpOnly: true,
		Secure:   config.Env.AppEnv != "local",
		SameSite: http.SameSiteLaxMode,
	}
	return &sessionStore{store: store, sessionName: sessionName}
}

// SaveSessionToken セッショントークンをクッキーに保存
func (s *sessionStore) SaveSessionToken(r *http.Request, w http.ResponseWriter, token string) error {
	session, err := s.store.Get(r, s.sessionName)
	if err != nil {
		slog.Error("failed to get session for save", "sessionName", s.sessionName, "error", err)
		return err
	}

	session.Values["session_token"] = token
	return session.Save(r, w)
}

// GetSessionToken クッキーからセッショントークンを取得
func (s *sessionStore) GetSessionToken(r *http.Request) (string, bool, error) {
	session, err := s.store.Get(r, s.sessionName)
	if err != nil {
		slog.Error("failed to get session for read", "sessionName", s.sessionName, "error", err)
		return "", false, err
	}

	token, ok := session.Values["session_token"].(string)
	if !ok || token == "" {
		return "", false, nil
	}

	return token, true, nil
}

// DeleteSessionToken クッキーからセッショントークンを削除
func (s *sessionStore) DeleteSessionToken(r *http.Request, w http.ResponseWriter) error {
	session, err := s.store.Get(r, s.sessionName)
	if err != nil {
		slog.Error("failed to get session for delete", "sessionName", s.sessionName, "error", err)
		return err
	}

	delete(session.Values, "session_token")
	session.Options.MaxAge = -1

	return session.Save(r, w)
}
