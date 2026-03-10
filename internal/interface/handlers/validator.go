package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
}

// FieldError バリデーションエラーの1フィールド分
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ValidationError バリデーションエラー全体
type ValidationError struct {
	Fields []FieldError
}

func (e *ValidationError) Error() string {
	msgs := make([]string, len(e.Fields))
	for i, f := range e.Fields {
		msgs[i] = fmt.Sprintf("%s: %s", f.Field, f.Message)
	}
	return strings.Join(msgs, "; ")
}

// Validate 構造体のバリデーションを実行する
func Validate(v any) error {
	err := validate.Struct(v)
	if err == nil {
		return nil
	}

	var fields []FieldError
	for _, fe := range err.(validator.ValidationErrors) {
		fields = append(fields, FieldError{
			Field:   toSnakeCase(fe.Field()),
			Message: fieldErrorMessage(fe),
		})
	}
	return &ValidationError{Fields: fields}
}

// fieldErrorMessage タグに応じた日本語メッセージを返す
func fieldErrorMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "必須項目です"
	case "email":
		return "有効なメールアドレスを入力してください"
	case "min":
		return fmt.Sprintf("%s 以上の値を指定してください", fe.Param())
	case "max":
		return fmt.Sprintf("%s 以下の値を指定してください", fe.Param())
	case "oneof":
		return fmt.Sprintf("%s のいずれかを指定してください", fe.Param())
	default:
		return "不正な値です"
	}
}

// toSnakeCase PascalCase/camelCase を snake_case に変換する
func toSnakeCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result.WriteByte('_')
			}
			result.WriteRune(unicode.ToLower(r))
		} else {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// ValidationErrorResponse バリデーションエラーのレスポンス
type ValidationErrorResponse struct {
	Error   string       `json:"error"`
	Message string       `json:"message"`
	Details []FieldError `json:"details"`
}

// HandleValidationError バリデーションエラーを400レスポンスとして返す
func HandleValidationError(w http.ResponseWriter, ve *ValidationError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	_ = json.NewEncoder(w).Encode(ValidationErrorResponse{
		Error:   "validation_error",
		Message: "入力内容に誤りがあります",
		Details: ve.Fields,
	})
}
