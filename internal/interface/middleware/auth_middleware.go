package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/usecase"
)

// UserContextKey コンテキストにユーザー情報を保存するためのキー
type UserContextKey struct{}

// OperatorContextKey コンテキストにオペレーター情報を保存するためのキー
type OperatorContextKey struct{}

// AuthMiddleware 認証ミドルウェア
type AuthMiddleware struct {
	userAuthUseCase     usecase.UserAuthUseCase
	operatorAuthUseCase usecase.OperatorAuthUseCase
}

// NewAuthMiddleware 新しいAuthMiddlewareインスタンスを作成
func NewAuthMiddleware(userAuthUseCase usecase.UserAuthUseCase, operatorAuthUseCase usecase.OperatorAuthUseCase) *AuthMiddleware {
	return &AuthMiddleware{
		userAuthUseCase:     userAuthUseCase,
		operatorAuthUseCase: operatorAuthUseCase,
	}
}

// RequireAuth 認証が必要なエンドポイント用のミドルウェア（customer_users向け）
func (m *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, authenticated, err := m.userAuthUseCase.GetCurrentUser(ctx, r)
		if err != nil {
			slog.Warn("authentication failed", "reason", "internal error", "error", err)
			http.Error(w, "認証エラー", http.StatusInternalServerError)
			return
		}

		if !authenticated || user == nil {
			slog.Warn("authentication failed", "reason", "not authenticated")
			http.Error(w, "認証が必要です", http.StatusUnauthorized)
			return
		}

		if !user.IsActive() {
			slog.Warn("authentication failed", "reason", "account inactive", "uid", user.UID, "email", user.Email)
			http.Error(w, "アカウントが無効です", http.StatusForbidden)
			return
		}

		ctx = context.WithValue(ctx, UserContextKey{}, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequireOperatorAuth 認証が必要なエンドポイント用のミドルウェア（operator向け）
func (m *AuthMiddleware) RequireOperatorAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		operator, authenticated, err := m.operatorAuthUseCase.GetCurrentOperator(ctx, r)
		if err != nil {
			slog.Warn("operator authentication failed", "reason", "internal error", "error", err)
			http.Error(w, "認証エラー", http.StatusInternalServerError)
			return
		}

		if !authenticated || operator == nil {
			slog.Warn("operator authentication failed", "reason", "not authenticated")
			http.Error(w, "認証が必要です", http.StatusUnauthorized)
			return
		}

		if !operator.IsActive {
			slog.Warn("operator authentication failed", "reason", "account inactive", "uid", operator.UID, "email", operator.Email)
			http.Error(w, "アカウントが無効です", http.StatusForbidden)
			return
		}

		ctx = context.WithValue(ctx, OperatorContextKey{}, operator)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// OptionalAuth 認証があれば情報を取得するが、なくてもOKなミドルウェア
func (m *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		user, authenticated, _ := m.userAuthUseCase.GetCurrentUser(ctx, r)

		if authenticated && user != nil {
			ctx = context.WithValue(ctx, UserContextKey{}, user)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}

		next.ServeHTTP(w, r)
	})
}

// GetUserFromContext コンテキストからユーザー情報を取得
func GetUserFromContext(ctx context.Context) (*domain.User, bool) {
	user, ok := ctx.Value(UserContextKey{}).(*domain.User)
	return user, ok
}

// GetOperatorFromContext コンテキストからオペレーター情報を取得
func GetOperatorFromContext(ctx context.Context) (*domain.Operator, bool) {
	operator, ok := ctx.Value(OperatorContextKey{}).(*domain.Operator)
	return operator, ok
}
