package adapter

import (
	"context"
	"fmt"
	"time"

	"github.com/zuxt268/berry/internal/config"
	"github.com/zuxt268/berry/internal/domain"
	"github.com/zuxt268/berry/internal/usecase/port"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	searchconsole "google.golang.org/api/searchconsole/v1"
)

type gscDataAdapter struct {
	oauthConfig *oauth2.Config
}

func NewGSCDataAdapter() port.GSCDataAdapter {
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

func (a *gscDataAdapter) FetchDailyReport(ctx context.Context, refreshToken string, siteURL string, date time.Time) (*domain.GSCDailyReport, error) {
	service, err := a.buildService(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrGSCTokenRefresh, err)
	}

	dateStr := date.Format("2006-01-02")

	report := &domain.GSCDailyReport{}

	// リクエスト1: サマリー指標
	if err := a.fetchSummary(ctx, service, siteURL, dateStr, report); err != nil {
		return nil, err
	}

	// リクエスト2: キーワード別
	if err := a.fetchKeywords(ctx, service, siteURL, dateStr, report); err != nil {
		return nil, err
	}

	// リクエスト3: ページ別
	if err := a.fetchPages(ctx, service, siteURL, dateStr, report); err != nil {
		return nil, err
	}

	return report, nil
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
func (a *gscDataAdapter) fetchSummary(ctx context.Context, service *searchconsole.Service, siteURL, dateStr string, report *domain.GSCDailyReport) error {
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
		report.Impressions = int(row.Impressions)
		report.Clicks = int(row.Clicks)
		report.CTR = row.Ctr
		report.AveragePosition = row.Position
	}

	return nil
}

// fetchKeywords キーワード別データを取得（上位100件）
func (a *gscDataAdapter) fetchKeywords(ctx context.Context, service *searchconsole.Service, siteURL, dateStr string, report *domain.GSCDailyReport) error {
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
		report.KeywordBreakdown = append(report.KeywordBreakdown, domain.KeywordBreakdown{
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
func (a *gscDataAdapter) fetchPages(ctx context.Context, service *searchconsole.Service, siteURL, dateStr string, report *domain.GSCDailyReport) error {
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
		report.PageBreakdown = append(report.PageBreakdown, domain.GSCPageBreakdown{
			Page:            row.Keys[0],
			Impressions:     int(row.Impressions),
			Clicks:          int(row.Clicks),
			CTR:             row.Ctr,
			AveragePosition: row.Position,
		})
	}

	return nil
}