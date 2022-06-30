package jackett

import "time"

// Response is a Jackett search response.
type Response struct {
	Results  []Result  `json:"Results"`
	Indexers []Indexer `json:"Indexers"`
}

// Result is a search result.
type Result struct {
	FirstSeen            string      `json:"FirstSeen"`
	Tracker              string      `json:"Tracker"`
	TrackerID            string      `json:"TrackerId"`
	TrackerType          string      `json:"TrackerType"`
	CategoryDesc         string      `json:"CategoryDesc"`
	BlackholeLink        string      `json:"BlackholeLink"`
	Title                string      `json:"Title"`
	GUID                 string      `json:"Guid"`
	Link                 string      `json:"Link"`
	Details              string      `json:"Details"`
	PublishDate          time.Time   `json:"PublishDate"`
	Category             []int       `json:"Category"`
	Size                 int64       `json:"Size"`
	Files                interface{} `json:"Files"`
	Grabs                int         `json:"Grabs"`
	Description          string      `json:"Description"`
	RageID               interface{} `json:"RageID"`
	TVDBID               interface{} `json:"TVDBId"`
	Imdb                 interface{} `json:"Imdb"`
	TMDb                 interface{} `json:"TMDb"`
	DoubanID             interface{} `json:"DoubanId"`
	Author               string      `json:"Author"`
	BookTitle            string      `json:"BookTitle"`
	Seeders              int         `json:"Seeders"`
	Peers                int         `json:"Peers"`
	Poster               interface{} `json:"Poster"`
	InfoHash             interface{} `json:"InfoHash"`
	MagnetURI            interface{} `json:"MagnetUri"`
	MinimumRatio         interface{} `json:"MinimumRatio"`
	MinimumSeedTime      interface{} `json:"MinimumSeedTime"`
	DownloadVolumeFactor float64     `json:"DownloadVolumeFactor"`
	UploadVolumeFactor   float64     `json:"UploadVolumeFactor"`
	Gain                 float64     `json:"Gain"`
}

// Indexer configured on Jackett.
type Indexer struct {
	ID      string `json:"ID"`
	Name    string `json:"Name"`
	Status  int    `json:"Status"`
	Results int    `json:"Results"`
	Error   string `json:"Error"`
}
