package domain

import "time"

type InstagramDailyReport struct {
	ID                      uint64
	InstagramConnectionID   int64
	ReportDate              time.Time
	FollowerCount           int
	Impressions             int
	Reach                   int
	ProfileViews            int
	WebsiteClicks           int
	PostEngagement          []PostEngagement
	AudienceDemographics    *AudienceDemographics
	StoriesInsights         []StoriesInsight
	FetchedAt               time.Time
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type PostEngagement struct {
	MediaID    string `json:"media_id"`
	MediaType  string `json:"media_type"`
	Timestamp  string `json:"timestamp"`
	LikeCount  int    `json:"like_count"`
	Comments   int    `json:"comments"`
	Reach      int    `json:"reach"`
	Impressions int   `json:"impressions"`
	Saved      int    `json:"saved"`
	Engagement int    `json:"engagement"`
}

type AudienceDemographics struct {
	GenderAge []GenderAgeBreakdown `json:"gender_age"`
	Cities    []CityBreakdown      `json:"cities"`
	Countries []CountryBreakdown   `json:"countries"`
}

type GenderAgeBreakdown struct {
	Key   string `json:"key"`
	Value int    `json:"value"`
}

type CityBreakdown struct {
	City  string `json:"city"`
	Value int    `json:"value"`
}

type CountryBreakdown struct {
	Country string `json:"country"`
	Value   int    `json:"value"`
}

type StoriesInsight struct {
	MediaID     string `json:"media_id"`
	Timestamp   string `json:"timestamp"`
	Impressions int    `json:"impressions"`
	Reach       int    `json:"reach"`
	Replies     int    `json:"replies"`
	Exits       int    `json:"exits"`
	TapsForward int    `json:"taps_forward"`
	TapsBack    int    `json:"taps_back"`
}