package jackett

import (
	"net/url"
	"strconv"

	"github.com/SkYNewZ/go-jackett/trackers"
)

// Option describes a query option.
type Option func(values *url.Values)

// WithQuery specify the query.
func WithQuery(query string) Option {
	return func(q *url.Values) {
		q.Set("Query", query)
	}
}

// WithCategory specify category to filter results with.
func WithCategory(category ...uint) Option {
	return func(q *url.Values) {
		for _, c := range category {
			q.Add("Category[]", strconv.Itoa(int(c)))
		}
	}
}

// WithTracker specify tracker to filter results with.
func WithTracker(tracker ...trackers.Tracker) Option {
	return func(q *url.Values) {
		for _, t := range tracker {
			q.Add("Tracker[]", t.String())
		}
	}
}
