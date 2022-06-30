// Code generated by ifacemaker; DO NOT EDIT.

package jackett

import (
	"context"
)

// Client interface describes wrapped Jackett client.
type Client interface {
	// Fetch make a query to Jackett and returns found results.
	// categoryID must match an existent category ID from Jackett. You can check them on your Jackett instance's indexers configuration.
	// You can specify multiples categories.
	// Custom indexer can be specified. Default will be "all".
	Fetch(ctx context.Context, options ...Option) (*Response, error)
}