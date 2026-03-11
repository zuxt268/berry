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

const instagramGraphAPIBase = "https://graph.facebook.com/v21.0"

type instagramDataAdapter struct{}

func NewInstagramDataAdapter() port.InstagramDataAdapter {
	return &instagramDataAdapter{}
}

func (a *instagramDataAdapter) FetchDailyReport(ctx context.Context, accessToken string, igAccountID string, date time.Time) (*domain.InstagramDailyReport, error) {
	report := &domain.InstagramDailyReport{}

	// 1. アカウントレベル日次インサイト
	if err := a.fetchAccountInsights(ctx, accessToken, igAccountID, date, report); err != nil {
		return nil, err
	}

	// 2. 投稿別エンゲージメント
	if err := a.fetchPostEngagements(ctx, accessToken, igAccountID, date, report); err != nil {
		return nil, err
	}

	// 3. フォロワー属性（lifetime）
	if err := a.fetchAudienceDemographics(ctx, accessToken, igAccountID, report); err != nil {
		// フォロワー属性はフォロワー100人未満だと取得不可の場合があるのでスキップ
		report.AudienceDemographics = nil
	}

	// 4. ストーリーズインサイト
	if err := a.fetchStoriesInsights(ctx, accessToken, igAccountID, report); err != nil {
		// ストーリーズがない場合はスキップ
		report.StoriesInsights = nil
	}

	return report, nil
}

// fetchAccountInsights アカウントレベルの日次インサイトを取得
func (a *instagramDataAdapter) fetchAccountInsights(ctx context.Context, accessToken string, igAccountID string, date time.Time, report *domain.InstagramDailyReport) error {
	since := date.Unix()
	until := date.AddDate(0, 0, 1).Unix()

	url := fmt.Sprintf(
		"%s/%s/insights?metric=follower_count,impressions,reach,profile_views,website_clicks&period=day&since=%d&until=%d&access_token=%s",
		instagramGraphAPIBase, igAccountID, since, until, accessToken,
	)

	body, err := a.doGet(ctx, url)
	if err != nil {
		return fmt.Errorf("%w: account insights: %v", domain.ErrInstagramAPICall, err)
	}

	var result struct {
		Data []struct {
			Name   string `json:"name"`
			Values []struct {
				Value int `json:"value"`
			} `json:"values"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("%w: account insights decode: %v", domain.ErrInstagramAPICall, err)
	}

	for _, metric := range result.Data {
		if len(metric.Values) == 0 {
			continue
		}
		value := metric.Values[0].Value
		switch metric.Name {
		case "follower_count":
			report.FollowerCount = value
		case "impressions":
			report.Impressions = value
		case "reach":
			report.Reach = value
		case "profile_views":
			report.ProfileViews = value
		case "website_clicks":
			report.WebsiteClicks = value
		}
	}

	return nil
}

// fetchPostEngagements 投稿別エンゲージメントを取得
func (a *instagramDataAdapter) fetchPostEngagements(ctx context.Context, accessToken string, igAccountID string, date time.Time, report *domain.InstagramDailyReport) error {
	since := date.Unix()
	until := date.AddDate(0, 0, 1).Unix()

	// 対象日に作成されたメディアを取得
	url := fmt.Sprintf(
		"%s/%s/media?fields=id,media_type,timestamp,like_count,comments_count&since=%d&until=%d&access_token=%s",
		instagramGraphAPIBase, igAccountID, since, until, accessToken,
	)

	body, err := a.doGet(ctx, url)
	if err != nil {
		return fmt.Errorf("%w: media list: %v", domain.ErrInstagramAPICall, err)
	}

	var mediaResult struct {
		Data []struct {
			ID            string `json:"id"`
			MediaType     string `json:"media_type"`
			Timestamp     string `json:"timestamp"`
			LikeCount     int    `json:"like_count"`
			CommentsCount int    `json:"comments_count"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &mediaResult); err != nil {
		return fmt.Errorf("%w: media decode: %v", domain.ErrInstagramAPICall, err)
	}

	for _, media := range mediaResult.Data {
		engagement := domain.PostEngagement{
			MediaID:   media.ID,
			MediaType: media.MediaType,
			Timestamp: media.Timestamp,
			LikeCount: media.LikeCount,
			Comments:  media.CommentsCount,
		}

		// 各メディアのインサイトを取得
		insightsURL := fmt.Sprintf(
			"%s/%s/insights?metric=impressions,reach,saved,engagement&access_token=%s",
			instagramGraphAPIBase, media.ID, accessToken,
		)

		insightsBody, err := a.doGet(ctx, insightsURL)
		if err == nil {
			var insightsResult struct {
				Data []struct {
					Name   string `json:"name"`
					Values []struct {
						Value int `json:"value"`
					} `json:"values"`
				} `json:"data"`
			}
			if json.Unmarshal(insightsBody, &insightsResult) == nil {
				for _, insight := range insightsResult.Data {
					if len(insight.Values) == 0 {
						continue
					}
					switch insight.Name {
					case "impressions":
						engagement.Impressions = insight.Values[0].Value
					case "reach":
						engagement.Reach = insight.Values[0].Value
					case "saved":
						engagement.Saved = insight.Values[0].Value
					case "engagement":
						engagement.Engagement = insight.Values[0].Value
					}
				}
			}
		}

		report.PostEngagement = append(report.PostEngagement, engagement)
	}

	return nil
}

// fetchAudienceDemographics フォロワー属性を取得（lifetime）
func (a *instagramDataAdapter) fetchAudienceDemographics(ctx context.Context, accessToken string, igAccountID string, report *domain.InstagramDailyReport) error {
	url := fmt.Sprintf(
		"%s/%s/insights?metric=audience_city,audience_country,audience_gender_age&period=lifetime&access_token=%s",
		instagramGraphAPIBase, igAccountID, accessToken,
	)

	body, err := a.doGet(ctx, url)
	if err != nil {
		return fmt.Errorf("%w: audience: %v", domain.ErrInstagramAPICall, err)
	}

	var result struct {
		Data []struct {
			Name   string `json:"name"`
			Values []struct {
				Value map[string]float64 `json:"value"`
			} `json:"values"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("%w: audience decode: %v", domain.ErrInstagramAPICall, err)
	}

	demographics := &domain.AudienceDemographics{}

	for _, metric := range result.Data {
		if len(metric.Values) == 0 {
			continue
		}
		values := metric.Values[0].Value
		switch metric.Name {
		case "audience_gender_age":
			for key, val := range values {
				demographics.GenderAge = append(demographics.GenderAge, domain.GenderAgeBreakdown{
					Key:   key,
					Value: int(val),
				})
			}
		case "audience_city":
			for city, val := range values {
				demographics.Cities = append(demographics.Cities, domain.CityBreakdown{
					City:  city,
					Value: int(val),
				})
			}
		case "audience_country":
			for country, val := range values {
				demographics.Countries = append(demographics.Countries, domain.CountryBreakdown{
					Country: country,
					Value:   int(val),
				})
			}
		}
	}

	report.AudienceDemographics = demographics
	return nil
}

// fetchStoriesInsights ストーリーズインサイトを取得
func (a *instagramDataAdapter) fetchStoriesInsights(ctx context.Context, accessToken string, igAccountID string, report *domain.InstagramDailyReport) error {
	// アクティブなストーリーズを取得
	url := fmt.Sprintf(
		"%s/%s/stories?fields=id,timestamp&access_token=%s",
		instagramGraphAPIBase, igAccountID, accessToken,
	)

	body, err := a.doGet(ctx, url)
	if err != nil {
		return fmt.Errorf("%w: stories list: %v", domain.ErrInstagramAPICall, err)
	}

	var storiesResult struct {
		Data []struct {
			ID        string `json:"id"`
			Timestamp string `json:"timestamp"`
		} `json:"data"`
	}
	if err := json.Unmarshal(body, &storiesResult); err != nil {
		return fmt.Errorf("%w: stories decode: %v", domain.ErrInstagramAPICall, err)
	}

	for _, story := range storiesResult.Data {
		insight := domain.StoriesInsight{
			MediaID:   story.ID,
			Timestamp: story.Timestamp,
		}

		insightsURL := fmt.Sprintf(
			"%s/%s/insights?metric=impressions,reach,replies,exits,taps_forward,taps_back&access_token=%s",
			instagramGraphAPIBase, story.ID, accessToken,
		)

		insightsBody, err := a.doGet(ctx, insightsURL)
		if err == nil {
			var insightsResult struct {
				Data []struct {
					Name   string `json:"name"`
					Values []struct {
						Value int `json:"value"`
					} `json:"values"`
				} `json:"data"`
			}
			if json.Unmarshal(insightsBody, &insightsResult) == nil {
				for _, m := range insightsResult.Data {
					if len(m.Values) == 0 {
						continue
					}
					switch m.Name {
					case "impressions":
						insight.Impressions = m.Values[0].Value
					case "reach":
						insight.Reach = m.Values[0].Value
					case "replies":
						insight.Replies = m.Values[0].Value
					case "exits":
						insight.Exits = m.Values[0].Value
					case "taps_forward":
						insight.TapsForward = m.Values[0].Value
					case "taps_back":
						insight.TapsBack = m.Values[0].Value
					}
				}
			}
		}

		report.StoriesInsights = append(report.StoriesInsights, insight)
	}

	return nil
}

// doGet GETリクエストを実行してレスポンスボディを返す
func (a *instagramDataAdapter) doGet(ctx context.Context, url string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

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