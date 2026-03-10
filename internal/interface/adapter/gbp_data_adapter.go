package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/zuxt268/berry/internal/config"
	"github.com/zuxt268/berry/internal/domain"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	businessprofileperformance "google.golang.org/api/businessprofileperformance/v1"
	"google.golang.org/api/option"
)

// GBPDataAdapter Google Business Profile APIからレポートデータを取得するアダプター
//
//go:generate mockgen -source=$GOFILE -destination=./mock/mock_$GOFILE -package mock
type GBPDataAdapter interface {
	FetchDailyReport(ctx context.Context, refreshToken string, accountID string, locationID string, date time.Time) (*GBPReportData, error)
}

// GBPReportData GBPから取得したレポートデータ
type GBPReportData struct {
	ProfileViews         int
	PhoneCalls           int
	DirectionRequests    int
	PhotoViews           int
	ReviewCount          int
	AverageRating        float64
	SearchQueryBreakdown []domain.SearchQueryBreakdown
}

type gbpDataAdapter struct {
	oauthConfig *oauth2.Config
}

func NewGBPDataAdapter() GBPDataAdapter {
	oauthConfig := &oauth2.Config{
		ClientID:     config.Env.GoogleClientID,
		ClientSecret: config.Env.GoogleClientSecret,
		Scopes: []string{
			"https://www.googleapis.com/auth/business.manage",
		},
		Endpoint: google.Endpoint,
	}
	return &gbpDataAdapter{oauthConfig: oauthConfig}
}

func (a *gbpDataAdapter) FetchDailyReport(ctx context.Context, refreshToken string, accountID string, locationID string, date time.Time) (*GBPReportData, error) {
	service, err := a.buildService(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrGBPTokenRefresh, err)
	}

	// ロケーションIDのフォーマット確認
	if !strings.HasPrefix(locationID, "locations/") {
		locationID = "locations/" + locationID
	}

	data := &GBPReportData{}

	// リクエスト1: パフォーマンス指標（閲覧数・電話・ルート検索など）
	if err := a.fetchPerformanceMetrics(ctx, service, locationID, date, data); err != nil {
		return nil, err
	}

	// リクエスト2: 検索キーワード（月次データ）
	if err := a.fetchSearchKeywords(ctx, service, locationID, date, data); err != nil {
		return nil, err
	}

	// リクエスト3: クチコミ数・平均評価（REST直接呼び出し）
	if err := a.fetchReviews(ctx, refreshToken, accountID, locationID, data); err != nil {
		return nil, err
	}

	return data, nil
}

func (a *gbpDataAdapter) buildService(ctx context.Context, refreshToken string) (*businessprofileperformance.Service, error) {
	token := &oauth2.Token{
		RefreshToken: refreshToken,
	}
	tokenSource := a.oauthConfig.TokenSource(ctx, token)
	service, err := businessprofileperformance.NewService(ctx, option.WithTokenSource(tokenSource))
	if err != nil {
		return nil, err
	}
	return service, nil
}

// fetchPerformanceMetrics パフォーマンス指標を一括取得
func (a *gbpDataAdapter) fetchPerformanceMetrics(ctx context.Context, service *businessprofileperformance.Service, locationID string, date time.Time, data *GBPReportData) error {
	year := int64(date.Year())
	month := int64(date.Month())
	day := int64(date.Day())

	// 終了日は翌日（APIは排他的終了日）
	endDate := date.AddDate(0, 0, 1)
	endYear := int64(endDate.Year())
	endMonth := int64(endDate.Month())
	endDay := int64(endDate.Day())

	resp, err := service.Locations.FetchMultiDailyMetricsTimeSeries(locationID).
		DailyMetrics(
			"BUSINESS_IMPRESSIONS_DESKTOP_MAPS",
			"BUSINESS_IMPRESSIONS_DESKTOP_SEARCH",
			"BUSINESS_IMPRESSIONS_MOBILE_MAPS",
			"BUSINESS_IMPRESSIONS_MOBILE_SEARCH",
			"BUSINESS_DIRECTION_REQUESTS",
			"CALL_CLICKS",
			"WEBSITE_CLICKS",
		).
		DailyRangeStartDateYear(year).
		DailyRangeStartDateMonth(month).
		DailyRangeStartDateDay(day).
		DailyRangeEndDateYear(endYear).
		DailyRangeEndDateMonth(endMonth).
		DailyRangeEndDateDay(endDay).
		Context(ctx).
		Do()
	if err != nil {
		return fmt.Errorf("%w: performance metrics: %v", domain.ErrGBPAPICall, err)
	}

	for _, multi := range resp.MultiDailyMetricTimeSeries {
		for _, ts := range multi.DailyMetricTimeSeries {
			value := extractDailyValue(ts, date)
			switch ts.DailyMetric {
			case "BUSINESS_IMPRESSIONS_DESKTOP_MAPS",
				"BUSINESS_IMPRESSIONS_DESKTOP_SEARCH",
				"BUSINESS_IMPRESSIONS_MOBILE_MAPS",
				"BUSINESS_IMPRESSIONS_MOBILE_SEARCH":
				data.ProfileViews += value
			case "BUSINESS_DIRECTION_REQUESTS":
				data.DirectionRequests = value
			case "CALL_CLICKS":
				data.PhoneCalls = value
			case "WEBSITE_CLICKS":
				// WebsiteClicksはPhotoViewsの代わりに取得
				// GBP APIに写真閲覧数の直接的な指標がないため、ウェブサイトクリック数で代用
				data.PhotoViews = value
			}
		}
	}

	return nil
}

// fetchSearchKeywords 検索キーワード別インプレッションを取得（月次データ）
func (a *gbpDataAdapter) fetchSearchKeywords(ctx context.Context, service *businessprofileperformance.Service, locationID string, date time.Time, data *GBPReportData) error {
	year := int64(date.Year())
	month := int64(date.Month())

	resp, err := service.Locations.Searchkeywords.Impressions.Monthly.List(locationID).
		MonthlyRangeStartMonthYear(year).
		MonthlyRangeStartMonthMonth(month).
		MonthlyRangeStartMonthDay(1).
		MonthlyRangeEndMonthYear(year).
		MonthlyRangeEndMonthMonth(month).
		MonthlyRangeEndMonthDay(1).
		PageSize(100).
		Context(ctx).
		Do()
	if err != nil {
		// 検索キーワードが取得できない場合はスキップ（データ不足の可能性）
		return nil
	}

	for _, kw := range resp.SearchKeywordsCounts {
		impressions := 0
		if kw.InsightsValue != nil {
			impressions = int(kw.InsightsValue.Value)
		}
		data.SearchQueryBreakdown = append(data.SearchQueryBreakdown, domain.SearchQueryBreakdown{
			Query:       kw.SearchKeyword,
			Impressions: impressions,
		})
	}

	return nil
}

// gbpReviewsResponse GBP Reviews APIのレスポンス構造
type gbpReviewsResponse struct {
	Reviews          []gbpReview `json:"reviews"`
	AverageRating    float64     `json:"averageRating"`
	TotalReviewCount int         `json:"totalReviewCount"`
	NextPageToken    string      `json:"nextPageToken"`
}

type gbpReview struct {
	ReviewID   string `json:"reviewId"`
	StarRating string `json:"starRating"`
}

// fetchReviews クチコミ数・平均評価をREST直接呼び出しで取得
func (a *gbpDataAdapter) fetchReviews(ctx context.Context, refreshToken string, accountID, locationID string, data *GBPReportData) error {
	// アクセストークンを取得
	token := &oauth2.Token{RefreshToken: refreshToken}
	tokenSource := a.oauthConfig.TokenSource(ctx, token)
	accessToken, err := tokenSource.Token()
	if err != nil {
		return fmt.Errorf("%w: reviews token refresh: %v", domain.ErrGBPTokenRefresh, err)
	}

	// accountID / locationID のフォーマット調整
	if !strings.HasPrefix(accountID, "accounts/") {
		accountID = "accounts/" + accountID
	}
	// locationIDは "locations/xxx" 形式で来るので、数値部分だけ取得
	locID := strings.TrimPrefix(locationID, "locations/")

	url := fmt.Sprintf(
		"https://mybusiness.googleapis.com/v4/%s/locations/%s/reviews",
		accountID, locID,
	)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return fmt.Errorf("%w: reviews request: %v", domain.ErrGBPAPICall, err)
	}
	req.Header.Set("Authorization", "Bearer "+accessToken.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("%w: reviews: %v", domain.ErrGBPAPICall, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("%w: reviews: status %d: %s", domain.ErrGBPAPICall, resp.StatusCode, string(body))
	}

	var result gbpReviewsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("%w: reviews decode: %v", domain.ErrGBPAPICall, err)
	}

	data.ReviewCount = result.TotalReviewCount
	data.AverageRating = result.AverageRating

	return nil
}

// extractDailyValue 指定日のデータポイント値を抽出
func extractDailyValue(ts *businessprofileperformance.DailyMetricTimeSeries, date time.Time) int {
	if ts.TimeSeries == nil {
		return 0
	}
	for _, dv := range ts.TimeSeries.DatedValues {
		if dv.Date != nil &&
			dv.Date.Year == int64(date.Year()) &&
			dv.Date.Month == int64(date.Month()) &&
			dv.Date.Day == int64(date.Day()) {
			return int(dv.Value)
		}
	}
	return 0
}
