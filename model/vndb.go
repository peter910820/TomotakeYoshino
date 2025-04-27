package model

type VndbVnResponse struct {
	More    bool       `json:"more"`
	Results []VnResult `json:"results"`
}

type VnResult struct {
	AltTitle *string `json:"alttitle"`
	ID       string  `json:"id"`
	Image    Image   `json:"image"`
	Rating   float64 `json:"rating"`
	Title    string  `json:"title"`
}

type Image struct {
	URL string `json:"url"`
}
