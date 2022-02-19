package handler

type UniCache struct {
	Name         string   `json:"name"`
	Country      string   `json:"country"`
	AlphaTwoCode string   `json:"alpha_two_code"`
	WebPages     []string `json:"web_pages"`
}

type CountryCache struct {
}

type Response struct {
	Name     string   `json:"name"`
	Country  string   `json:"country"`
	IsoCode  string   `json:"isocode"`
	WebPages []string `json:"webpages"`
}
