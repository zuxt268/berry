package adapter

import (
	"context"
	"fmt"
	"time"

	"github.com/zuxt268/berry/internal/config"
	"github.com/zuxt268/berry/internal/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	searchconsole "google.golang.org/api/searchconsole/v1"
)

// GSCDataAdapter Search Console APIからレポートデータを取得するアダプター
//
//go:generate mockgen -source=$GOFILE -destination=./mock/mock_$GOFILE -package mock
type GSCDataAdapter interface {
	FetchDailyReport(ctx context.Context, refreshToken string, siteURL string, date time.Time) (*GSCReportData, error)
}

// GSCReportData Search Consoleから取得したレポートデータ
type GSCReportData struct {
	Impressions      int
	Clicks           int
	CTR              float64
	AveragePosition  float64
	KeywordBreakdown []domain.KeywordBreakdown
	PageBreakdown    []domain.GSCPageBreakdown
}

type gscDataAdapter struct {
	oauthConfig *oauth2.Config
}

func NewGSCDataAdapter() GSCDataAdapter {
	oauthConfig := &oauth2.Config{
		ClientID:     config.Env.GoogleClientID,
		ClientSecret: config.Env.GoogleClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/webmasters.readonly",
		},
		Endpoint: google.Endpoint,
	}
	return &gscDataAdapter{oauthConfig: oauthConfig}
}

func (a *gscDataAdapter) FetchDailyReport(ctx context.Context, refreshToken string, siteURL string, date time.Time) (*GSCReportData, error) {
	service, err := a.buildService(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrGSCTokenRefresh, err)
	}

	dateStr := date.Format("2006-01-02")

	data := &GSCReportData{}

	// リクエスト1: サマリー指標
	if err := a.fetchSummary(ctx, service, siteURL, dateStr, data); err != nil {
		return nil, err
	}

	// リクエスト2: キーワード別
	if err := a.fetchKeywords(ctx, service, siteURL, dateStr, data); err != nil {
		return nil, err
	}

	// リクエスト3: ページ別
	if err := a.fetchPages(ctx, service, siteURL, dateStr, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (a *gscDataAdapter) buildService(ctx context.Context, refreshToken string) (*searchconsole.Service, error) {
	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}
	tokenSource := a.oauthConfig.TokenSource(ctx, token)
	service, err := searchconsole.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, err
	}
	return service, nil
}

// fetchSummary サマリー指標を取得（impressions, clicks, ctr, position）
func (a *gscDataAdapter) fetchSummary(ctx context.Context, service *searchconsole.Service, siteURL, dateStr string, data *GSCReportData) error {
	req := &searchconsole.SearchAnalyticsQueryRequest{
		StartDate: dateStr,
		EndDate:   dateStr,
	}

	resp, err := service.Searchanalytics.Query(siteURL, req).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("%w: summary: %v", domain.ErrGSCAPICall, err)
	}

	if len(resp.Rows) > 0 {
		row := resp.Rows[0]
		data.Impressions = int(row.Impressions)
		data.Clicks = int(row.Clicks)
		data.CTR = row.Ctr
		data.AveragePosition = row.Position
	}

	return nil
}

// fetchKeywords キーワード別データを取得（上位100件）
func (a *gscDataAdapter) fetchKeywords(ctx context.Context, service *searchconsole.Service, siteURL, dateStr string, data *GSCReportData) error {
	req := &searchconsole.SearchAnalyticsQueryRequest{
		StartDate:  dateStr,
		EndDate:    dateStr,
		Dimensions: []string{"query"},
		RowLimit:   100,
	}

	resp, err := service.Searchanalytics.Query(siteURL, req).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("%w: keywords: %v", domain.ErrGSCAPICall, err)
	}

	for _, row := range resp.Rows {
		data.KeywordBreakdown = append(data.KeywordBreakdown, domain.KeywordBreakdown{
			Query:           row.Keys[0],
			Impressions:     int(row.Impressions),
			Clicks:          int(row.Clicks),
			CTR:             row.Ctr,
			AveragePosition: row.Position,
		})
	}

	return nil
}

// fetchPages ページ別データを取得（上位100件）
func (a *gscDataAdapter) fetchPages(ctx context.Context, service *searchconsole.Service, siteURL, dateStr string, data *GSCReportData) error {
	req := &searchconsole.SearchAnalyticsQueryRequest{
		StartDate:  dateStr,
		EndDate:    dateStr,
		Dimensions: []string{"page"},
		RowLimit:   100,
	}

	resp, err := service.Searchanalytics.Query(siteURL, req).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("%w: pages: %v", domain.ErrGSCAPICall, err)
	}

	for _, row := range resp.Rows {
		data.PageBreakdown = append(data.PageBreakdown, domain.GSCPageBreakdown{
			Page:            row.Keys[0],
			Impressions:     int(row.Impressions),
			Clicks:          int(row.Clicks),
			CTR:             row.Ctr,
			AveragePosition: row.Position,
		})
	}

	return nil
}
