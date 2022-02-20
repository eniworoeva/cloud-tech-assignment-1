package model

type UniCache struct {
	Name         string   `json:"name"`
	Country      string   `json:"country"`
	AlphaTwoCode string   `json:"alpha_two_code"`
	WebPages     []string `json:"web_pages"`
}

type CountryCache struct {
	CCA2      string            `json:"cca2"`
	Languages map[string]string `json:"languages"`
	Maps      map[string]string `json:"maps"`
}

type UniInfoResponse struct {
	Name      string            `json:"name"`
	Country   string            `json:"country"`
	IsoCode   string            `json:"isocode"`
	WebPages  []string          `json:"webpages"`
	Languages map[string]string `json:"languages"`
	Map       string            `json:"map"`
}

type BordersCache struct {
	CCA3    string   `json:"cca3"`
	Borders []string `json:"borders"`
}
