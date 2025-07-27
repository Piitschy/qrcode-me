package proxies

import "github.com/pocketbase/pocketbase/core"

var _ core.RecordProxy = (*Link)(nil)

type Link struct {
	// Id     string `db:"id"     json:"id"`
	// Slug   string `db:"slug"   json:"slug"`
	// Url    string `db:"url"    json:"url"`
	// Status string `db:"status" json:"status"`
	// Count  int    `db:"count"  json:"count"`
	core.BaseRecordProxy
}

func (l *Link) Id() string {
	return l.GetString("id")
}

func (l *Link) Slug() string {
	return l.GetString("slug")
}

func (l *Link) Url() string {
	return l.GetString("url")
}

func (l *Link) Status() string {
	return l.GetString("status")
}

func (l *Link) Count() int {
	return l.GetInt("count")
}

func (l *Link) SetSlug(slug string) {
	l.Set("slug", slug)
}

func (l *Link) SetUrl(url string) {
	l.Set("url", url)
}

func (l *Link) SetStatus(status string) {
	l.Set("status", status)
}

func (l *Link) SetCount(count int) {
	l.Set("count", count)
}

func (l *Link) IncrementCount() {
	l.SetCount(l.Count() + 1)
}
