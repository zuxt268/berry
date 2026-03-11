package handlers

import (
	"net/http"

	"github.com/zuxt268/berry/internal/config"
)

// setOAuthStateCookie OAuthのstateをクッキーに保存
func setOAuthStateCookie(w http.ResponseWriter, name, state string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    state,
		MaxAge:   300,
		HttpOnly: true,
		Secure:   config.Env.AppEnv != "local",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
}

// getOAuthStateCookie OAuthのstateをクッキーから取得
func getOAuthStateCookie(r *http.Request, name string) (string, error) {
	cookie, err := r.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// setValueCookie 値をクッキーに保存（OAuth連携の一時データ用）
func setValueCookie(w http.ResponseWriter, name, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		MaxAge:   300,
		HttpOnly: true,
		Secure:   config.Env.AppEnv != "local",
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
	})
}

// clearCookie クッキーを削除
func clearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    "",
		MaxAge:   -1,
		HttpOnly: true,
		Path:     "/",
	})
}

// verifyOAuthState OAuthのstateを検証
func verifyOAuthState(r *http.Request, cookieName, state string) error {
	cookieState, err := getOAuthStateCookie(r, cookieName)
	if err != nil {
		return err
	}
	if state != cookieState {
		return http.ErrNoCookie
	}
	return nil
}