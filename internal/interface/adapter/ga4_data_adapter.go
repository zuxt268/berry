package adapter

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/zuxt268/berry/internal/config"
	"github.com/zuxt268/berry/internal/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	analyticsdata "google.golang.org/api/analyticsdata/v1beta"
	"google.golang.org/api/option"
)

// GA4DataAdapter GA4 Data APIからレポートデータを取得するアダプター
//
//go:generate mockgen -source=$GOFILE -destination=./mock/mock_$GOFILE -package mock
type GA4DataAdapter interface {
	FetchDailyReport(ctx context.Context, refreshToken string, propertyID string, date time.Time) (*GA4ReportData, error)
}

// GA4ReportData GA4から取得したレポートデータ
type GA4ReportData struct {
	Sessions           int
	TotalUsers         int
	BounceRate         float64
	AvgSessionDuration float64
	Conversions        int
	ChannelBreakdown   []domain.ChannelBreakdown
	DeviceBreakdown    []domain.DeviceBreakdown
	PageBreakdown      []domain.PageBreakdown
}

type ga4DataAdapter struct {
	oauthConfig *oauth2.Config
}

func NewGA4DataAdapter() GA4DataAdapter {
	oauthConfig := &oauth2.Config{
		ClientID:     config.Env.GoogleClientID,
		ClientSecret: config.Env.GoogleClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/analytics.readonly",
		},
		Endpoint: google.Endpoint,
	}
	return &ga4DataAdapter{oauthConfig: oauthConfig}
}

func (a *ga4DataAdapter) FetchDailyReport(ctx context.Context, refreshToken string, propertyID string, date time.Time) (*GA4ReportData, error) {
	service, err := a.buildService(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrGA4TokenRefresh, err)
	}

	dateStr := date.Format("2006-01-02")

	// プロパティIDのフォーマット確認
	if !strings.HasPrefix(propertyID, "properties/") {
		propertyID = "properties/" + propertyID
	}

	data := &GA4ReportData{}

	// リクエスト1: サマリー指標
	if err := a.fetchSummary(ctx, service, propertyID, dateStr, data); err != nil {
		return nil, err
	}

	// リクエスト2: 流入経路 + デバイス別
	if err := a.fetchChannelDevice(ctx, service, propertyID, dateStr, data); err != nil {
		return nil, err
	}

	// リクエスト3: ページ別PV
	if err := a.fetchPages(ctx, service, propertyID, dateStr, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (a *ga4DataAdapter) buildService(ctx context.Context, refreshToken string) (*analyticsdata.Service, error) {
	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}
	tokenSource := a.oauthConfig.TokenSource(ctx, token)
	service, err := analyticsdata.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, err
	}
	return service, nil
}

// fetchSummary サマリー指標を取得（sessions, totalUsers, bounceRate, averageSessionDuration, conversions）
func (a *ga4DataAdapter) fetchSummary(ctx context.Context, service *analyticsdata.Service, propertyID, dateStr string, data *GA4ReportData) error {
	req := &analyticsdata.RunReportRequest{
		DateRanges: []*analyticsdata.DateRange{
			{StartDate: dateStr, EndDate: dateStr},
		},
		Metrics: []*analyticsdata.Metric{
			{Name: "sessions"},
			{Name: "totalUsers"},
			{Name: "bounceRate"},
			{Name: "averageSessionDuration"},
			{Name: "conversions"},
		},
	}

	resp, err := service.Properties.RunReport(propertyID, req).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("%w: summary: %v", domain.ErrGA4APICall, err)
	}

	if len(resp.Rows) > 0 {
		row := resp.Rows[0]
		data.Sessions = parseIntMetric(row.MetricValues, 0)
		data.TotalUsers = parseIntMetric(row.MetricValues, 1)
		data.BounceRate = parseFloatMetric(row.MetricValues, 2)
		data.AvgSessionDuration = parseFloatMetric(row.MetricValues, 3)
		data.Conversions = parseIntMetric(row.MetricValues, 4)
	}

	return nil
}

// fetchChannelDevice 流入経路別 + デバイス別を取得
func (a *ga4DataAdapter) fetchChannelDevice(ctx context.Context, service *analyticsdata.Service, propertyID, dateStr string, data *GA4ReportData) error {
	req := &analyticsdata.RunReportRequest{
		DateRanges: []*analyticsdata.DateRange{
			{StartDate: dateStr, EndDate: dateStr},
		},
		Dimensions: []*analyticsdata.Dimension{
			{Name: "sessionDefaultChannelGroup"},
			{Name: "deviceCategory"},
		},
		Metrics: []*analyticsdata.Metric{
			{Name: "sessions"},
			{Name: "totalUsers"},
		},
	}

	resp, err := service.Properties.RunReport(propertyID, req).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("%w: channel_device: %v", domain.ErrGA4APICall, err)
	}

	// チャンネル別とデバイス別に集約
	channelMap := make(map[string]*domain.ChannelBreakdown)
	deviceMap := make(map[string]*domain.DeviceBreakdown)

	for _, row := range resp.Rows {
		channel := row.DimensionValues[0].Value
		device := row.DimensionValues[1].Value
		sessions := parseIntMetric(row.MetricValues, 0)
		users := parseIntMetric(row.MetricValues, 1)

		if _, ok := channelMap[channel]; !ok {
			channelMap[channel] = &domain.ChannelBreakdown{Channel: channel}
		}
		channelMap[channel].Sessions += sessions
		channelMap[channel].Users += users

		if _, ok := deviceMap[device]; !ok {
			deviceMap[device] = &domain.DeviceBreakdown{DeviceCategory: device}
		}
		deviceMap[device].Sessions += sessions
		deviceMap[device].Users += users
	}

	for _, v := range channelMap {
		data.ChannelBreakdown = append(data.ChannelBreakdown, *v)
	}
	for _, v := range deviceMap {
		data.DeviceBreakdown = append(data.DeviceBreakdown, *v)
	}

	return nil
}

// fetchPages ページ別PVを取得（上位50件）
func (a *ga4DataAdapter) fetchPages(ctx context.Context, service *analyticsdata.Service, propertyID, dateStr string, data *GA4ReportData) error {
	req := &analyticsdata.RunReportRequest{
		DateRanges: []*analyticsdata.DateRange{
			{StartDate: dateStr, EndDate: dateStr},
		},
		Dimensions: []*analyticsdata.Dimension{
			{Name: "pagePath"},
		},
		Metrics: []*analyticsdata.Metric{
			{Name: "screenPageViews"},
			{Name: "averageSessionDuration"},
		},
		OrderBys: []*analyticsdata.OrderBy{
			{
				Metric: &analyticsdata.MetricOrderBy{MetricName: "screenPageViews"},
				Desc:   true,
			},
		},
		Limit: 50,
	}

	resp, err := service.Properties.RunReport(propertyID, req).Context(ctx).Do()
	if err != nil {
		return fmt.Errorf("%w: pages: %v", domain.ErrGA4APICall, err)
	}

	for _, row := range resp.Rows {
		data.PageBreakdown = append(data.PageBreakdown, domain.PageBreakdown{
			PagePath:      row.DimensionValues[0].Value,
			PageViews:     parseIntMetric(row.MetricValues, 0),
			AvgTimeOnPage: parseFloatMetric(row.MetricValues, 1),
		})
	}

	return nil
}

func parseIntMetric(values []*analyticsdata.MetricValue, idx int) int {
	if idx >= len(values) {
		return 0
	}
	v, _ := strconv.Atoi(values[idx].Value)
	return v
}

func parseFloatMetric(values []*analyticsdata.MetricValue, idx int) float64 {
	if idx >= len(values) {
		return 0
	}
	v, _ := strconv.ParseFloat(values[idx].Value, 64)
	return v
}
