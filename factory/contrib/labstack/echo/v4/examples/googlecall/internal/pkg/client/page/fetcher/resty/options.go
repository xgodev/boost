package resty

// Options configures the page to fetch; inject generates its config key
// (internal.pkg.client.page.fetcher.resty.url).
// @Config
type Options struct {
	URL string `default:"http://www.google.com" desc:"page the fetcher asks for"`
}
