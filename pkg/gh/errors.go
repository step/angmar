package gh

import "fmt"

// ClientFetchError is typically returned when an http.Get or an http.Post
// returns an error(other methods as well). Either due to a malformed URL
// or a connection refused etc. It accepts the url, the method used,
// the actual error returned by http.Get/Do etc and a location. The location
// is simply a string that helps you identify where the error occcurred.
// Usually a function name.
type ClientFetchError struct {
	url       string
	method    string
	actualErr error
	location  string
}

// Error returns a string that reports the method and the url that failed
// along with the location it failed and the actual error returned by http.Get/Do
func (c ClientFetchError) Error() string {
	return fmt.Sprintf("Unable to %s from/to %s at %s\n%s", c.method, c.url, c.location, c.actualErr)
}

// StatusCodeError is returned when a non 2xx or 3xx status code
// is returned on an http.Get/Post/Put etc. It accepts a url, the method used
// and the location where the error occurred.
type StatusCodeError struct {
	statusCode int
	url        string
	method     string
	location   string
}

// Error returns a string that reports the method, the url the status code
// along with the location it failed at.
func (s StatusCodeError) Error() string {
	return fmt.Sprintf("Got %d from %s %s at %s\n", s.statusCode, s.method, s.url, s.location)
}

// FetchUntarError is returned when an untar fails as a part of a fetch
// of some sort. Ideally this returned from FetchTarball in gh, but
// any fetch that untars and fails is a valid scenario for this error
// It accepts the url from where the untar fetch failed. The location
// is simply a string that helps you identify where the error occcurred.
// Usually a function name.
type FetchUntarError struct {
	url       string
	actualErr error
	location  string
}

// Error returns a string that reports the url from where the untar fetch failed
// It also reports the actual error along with the location it failed at.
func (f FetchUntarError) Error() string {
	return fmt.Sprintf("Unable to untar while fetching from %s at %s\n%s", f.url, f.location, f.actualErr)
}
