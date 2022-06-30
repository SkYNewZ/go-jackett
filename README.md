# go-jackett

Go [Jackett](https://github.com/Jackett/Jackett) client.

## Get started

```go
package main

import (
	"context"

	"github.com/SkYNewZ/go-jackett"
)

func main() {
	// Read JACKETT_API_URL and JACKETT_API_KEY by default
	client := jackett.New(nil)
	results, err := client.Fetch(context.Background(), jackett.WithQuery("foo"))
	if err != nil {
		panic(err)
	}

	_ = results // use results
}
```

## Trackers

Available trackers are statically types using [`cmd/jackett-sync`](cmd/jackett-sync) as code generator
from https://github.com/Jackett/Jackett/tree/master/src/Jackett.Common/Definitions.
You can see the output at [`trackers/trackers.go`](trackers/trackers.go).