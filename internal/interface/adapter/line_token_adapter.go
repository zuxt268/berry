package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/usecase/port"
)

type lineTokenAdapter struct{}

func NewLineTokenAdapter() port.LineTokenAdapter {
	return &lineTokenAdapter{}
}

// ValidateToken トークンの有効性を検証し、Bot情報を返す
func (a *lineTokenAdapter) ValidateToken(ctx context.Context, channelAccessToken string) (*port.LineBotInfo, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.line.me/v2/bot/info", nil)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrLineInvalidToken, err)
	}
	req.Header.Set("Authorization", "Bearer "+channelAccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrLineInvalidToken, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("%w: status %d: %s", domain.ErrLineInvalidToken, resp.StatusCode, string(body))
	}

	var info port.LineBotInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, fmt.Errorf("%w: decode: %v", domain.ErrLineInvalidToken, err)
	}

	return &info, nil
}
