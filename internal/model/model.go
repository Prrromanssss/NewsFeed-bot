package model

import "time"

type Item struct {
	Title      string
	Categories []string
	Link       string
	Date       time.Time
	Summary    string
	SourceName string
}

type Source struct {
	SourceID  int64
	Name      string
	FeedUrl   string
	CreatedAt time.Time
}

type Article struct {
	ArticleID   int64
	SourceID    int64
	Title       string
	Link        string
	Summary     string
	PublishedAt time.Time
	CreatedAt   time.Time
	PostedAt    time.Time
}
