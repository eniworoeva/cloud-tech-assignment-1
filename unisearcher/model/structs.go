package model

// UniCache is a struct that holds relevant information retrieved from the university API
type UniCache struct {
	Name         string   `json:"name"`
	Country      string   `json:"country"`
	AlphaTwoCode string   `json:"alpha_two_code"`
	WebPages     []string `json:"web_pages"`
}

// CountryCache is a struct that holds relevant information retrieved from the RESTCountries API
type CountryCache struct {
	CCA2      string            `json:"cca2"`
	Languages map[string]string `json:"languages"`
	Maps      map[string]string `json:"maps"`
}

// UniInfoResponse is a combined struct of UniCache and CountryCache and is used as the standard response
type UniInfoResponse struct {
	Name      string            `json:"name"`
	Country   string            `json:"country"`
	IsoCode   string            `json:"isocode"`
	WebPages  []string          `json:"webpages"`
	Languages map[string]string `json:"languages"`
	Map       string            `json:"map"`
}

// BordersCache is a struct that holds a countries cca3 code and it's bordering countries' cca3 codes
type BordersCache struct {
	CCA3    string   `json:"cca3"`
	Borders []string `json:"borders"`
}

// Diag Default struct for diagnostics
type Diag struct {
	UniversityAPI string `json:"universityAPI"`
	CountryAPI    string `json:"countryAPI"`
	Version       string `json:"version"`
	Uptime        string `json:"uptime"`
}
