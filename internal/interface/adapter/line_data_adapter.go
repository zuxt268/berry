package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/usecase/port"
)

const lineAPIBase = "https://api.line.me/v2/bot"

type lineDataAdapter struct{}

func NewLineDataAdapter() port.LineDataAdapter {
	return &lineDataAdapter{}
}

func (a *lineDataAdapter) FetchDailyReport(ctx context.Context, channelAccessToken string, date time.Time) (*domain.LineDailyReport, error) {
	report := &domain.LineDailyReport{}
	dateStr := date.Format("20060102")

	// 1. 友だち数
	if err := a.fetchFollowers(ctx, channelAccessToken, dateStr, report); err != nil {
		return nil, err
	}

	// 2. メッセージ配信統計
	if err := a.fetchMessageDelivery(ctx, channelAccessToken, dateStr, report); err != nil {
		// メッセージ配信がない日はスキップ
		report.MessageDelivery = nil
	}

	// 3. 友だち属性
	if err := a.fetchDemographic(ctx, channelAccessToken, report); err != nil {
		// 属性データが取得できない場合はスキップ
		report.Demographic = nil
	}

	return report, nil
}

// fetchFollowers 友だち数を取得
func (a *lineDataAdapter) fetchFollowers(ctx context.Context, token string, dateStr string, report *domain.LineDailyReport) error {
	url := fmt.Sprintf("%s/insight/followers?date=%s", lineAPIBase, dateStr)

	body, err := a.doGetWithAuth(ctx, url, token)
	if err != nil {
		return fmt.Errorf("%w: followers: %v", domain.ErrLineAPICall, err)
	}

	var result struct {
		Status          string `json:"status"`
		Followers       int    `json:"followers"`
		TargetedReaches int    `json:"targetedReaches"`
		Blocks          int    `json:"blocks"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("%w: followers decode: %v", domain.ErrLineAPICall, err)
	}

	// statusが"ready"でない場合はデータ未集計
	if result.Status != "ready" {
		return fmt.Errorf("%w: followers data not ready (status: %s)", domain.ErrLineAPICall, result.Status)
	}

	report.Followers = result.Followers
	report.TargetedReaches = result.TargetedReaches
	report.Blocks = result.Blocks

	return nil
}

// fetchMessageDelivery メッセージ配信統計を取得
func (a *lineDataAdapter) fetchMessageDelivery(ctx context.Context, token string, dateStr string, report *domain.LineDailyReport) error {
	url := fmt.Sprintf("%s/insight/message/delivery?date=%s", lineAPIBase, dateStr)

	body, err := a.doGetWithAuth(ctx, url, token)
	if err != nil {
		return fmt.Errorf("%w: message delivery: %v", domain.ErrLineAPICall, err)
	}

	var result struct {
		Status      string `json:"status"`
		Success     int    `json:"success"`
		UniqueClick int    `json:"uniqueClick"`
		UniqueOpen  int    `json:"uniqueOpen"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("%w: message delivery decode: %v", domain.ErrLineAPICall, err)
	}

	if result.Status != "ready" {
		return fmt.Errorf("message delivery data not ready (status: %s)", result.Status)
	}

	report.MessageDelivery = &domain.LineMessageDelivery{
		Status:      result.Status,
		Success:     result.Success,
		UniqueClick: result.UniqueClick,
		UniqueOpen:  result.UniqueOpen,
	}

	return nil
}

// fetchDemographic 友だち属性を取得
func (a *lineDataAdapter) fetchDemographic(ctx context.Context, token string, report *domain.LineDailyReport) error {
	url := fmt.Sprintf("%s/insight/demographic", lineAPIBase)

	body, err := a.doGetWithAuth(ctx, url, token)
	if err != nil {
		return fmt.Errorf("%w: demographic: %v", domain.ErrLineAPICall, err)
	}

	var result struct {
		Available bool `json:"available"`
		Genders   []struct {
			Gender     string  `json:"gender"`
			Percentage float64 `json:"percentage"`
		} `json:"genders"`
		Ages []struct {
			Age        string  `json:"age"`
			Percentage float64 `json:"percentage"`
		} `json:"ages"`
		Areas []struct {
			Area       string  `json:"area"`
			Percentage float64 `json:"percentage"`
		} `json:"areas"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("%w: demographic decode: %v", domain.ErrLineAPICall, err)
	}

	if !result.Available {
		return fmt.Errorf("demographic data not available")
	}

	demographic := &domain.LineDemographic{
		Available: result.Available,
	}

	for _, g := range result.Genders {
		demographic.Genders = append(demographic.Genders, domain.LineDemographicItem{
			Key:        g.Gender,
			Percentage: g.Percentage,
		})
	}
	for _, a := range result.Ages {
		demographic.Ages = append(demographic.Ages, domain.LineDemographicItem{
			Key:        a.Age,
			Percentage: a.Percentage,
		})
	}
	for _, a := range result.Areas {
		demographic.Areas = append(demographic.Areas, domain.LineDemographicItem{
			Key:        a.Area,
			Percentage: a.Percentage,
		})
	}

	report.Demographic = demographic
	return nil
}

// doGetWithAuth Bearer Token付きGETリクエストを実行
func (a *lineDataAdapter) doGetWithAuth(ctx context.Context, url string, token string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status %d: %s", resp.StatusCode, string(body))
	}

	return body, nil
}