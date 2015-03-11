package packagist

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

//PackageListResult is the envelope for a package list result
type PackageListResult struct {
	PackageNames []string `json:"packageNames"`
}

//PackageResult is the envelope for a package result
type PackageResult struct {
	Package Package `json:"package"`
}

//Package Entity
type Package struct {
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Time        string              `json:"time"`
	Maintainers []map[string]string `json:"maintainers"`
	Versions    map[string]Version  `json:"versions"`
	Type        string              `json:"type"`
	Repository  string              `json:"repository"`
	Download    map[string]int      `json:"download"`
	Favers      int                 `json:"favers"`
}

// Version of a Package
type Version struct {
	Name              string                       `json:"name"`
	Description       string                       `json:"description"`
	Keywords          []string                     `json:"keywords"`
	Homepage          string                       `json:"homepage"`
	Version           string                       `json:"version"`
	VersionNormalized string                       `json:"version_normalized"`
	License           []string                     `json:"license"`
	Authors           []map[string]string          `json:"authors"`
	Source            map[string]string            `json:"source"`
	Dist              map[string]string            `json:"dist"`
	Type              string                       `json:"type"`
	Time              string                       `json:"time"`
	Autoload          map[string]map[string]string `json:"autoload"`
	Require           map[string]string            `json:"require"`
	RequireDev        map[string]string            `json:"require-dev"`
	Suggest           map[string]string            `json:"suggest"`
	Extra             map[string]map[string]string `json:"extra"`
}

//Client for the Packagist API
type Client struct {
	host       string
	httpClient *http.Client
}

//ListPackages Gets a list of packages optionally filtered by type and vendor
func (p *Client) ListPackages(filters map[string]string) (PackageListResult, error) {

	//list
	var packageListResult PackageListResult

	//setup request
	req, err := http.NewRequest("GET", p.MakeURI("/packages/list.json", filters), nil)

	if err != nil {
		return packageListResult, err
	}

	//get response
	res, err := p.httpClient.Do(req)
	if err != nil {
		return packageListResult, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	decodeError := decoder.Decode(&packageListResult)

	return packageListResult, decodeError
}

//GetPackage gets a single package record by name
func (p *Client) GetPackage(name string) (PackageResult, error) {
	var packageResult PackageResult
	req, err := http.NewRequest("GET", p.MakeURI("/packages/"+name+".json", make(map[string]string)), nil)
	if err != nil {
		return packageResult, err
	}

	//get response
	res, err := p.httpClient.Do(req)
	if err != nil {
		return packageResult, err
	}
	defer res.Body.Close()

	decoder := json.NewDecoder(res.Body)
	decodeError := decoder.Decode(&packageResult)

	return packageResult, decodeError
}

//MakeURI creates a valid packagist API uri for a given path. params are converted to a querystring
func (p *Client) MakeURI(path string, paramMap map[string]string) string {

	uri, err := url.Parse(p.host + path)
	if err != nil {
		log.Fatal(err)
	}

	params := url.Values{}
	for filterName, filterValue := range paramMap {
		params.Set(filterName, filterValue)
	}

	uri.RawQuery = params.Encode()

	return uri.String()
}

//NewAPIClient creates a new Packagist API client
func NewAPIClient() *Client {
	client := Client{host: "https://packagist.org", httpClient: &http.Client{}}
	return &client
}
