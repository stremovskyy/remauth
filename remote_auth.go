package remauth

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/karmadon/remauth/cache"
)

type remoteAuth struct {
	options *Options
	client  *http.Client
	url     *url.URL
	cache   cache.Cache
}

func (r *remoteAuth) Valid(token string) bool {
	defer timeTrack(time.Now(), "valid method")
	cacheKey := fmt.Sprintf("%x", md5.Sum([]byte(token)))

	if valid, exists := r.cache.Get(cacheKey); exists {
		return valid
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    r.url,
		Header: http.Header{
			"User-Agent":    {"RemoteAuthClient/" + version},
			"Authorization": {token},
		},
	}
	if r.options.Debug {
		fmt.Printf("[RemoteAuthClient][URL]: %s\n", r.url.String())
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return false
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
		r.cache.Set(cacheKey, true, 5*time.Minute)

		return true
	}

	r.cache.Set(cacheKey, false, 5*time.Minute)
	return false
}

func timeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", name, elapsed)
}

func New(options *Options) (*remoteAuth, error) {
	cache := cache.New(5*time.Minute, 10*time.Minute)

	transport := &http.Transport{
		MaxIdleConnsPerHost: -1,
		TLSHandshakeTimeout: 5 * time.Second,
		DisableKeepAlives:   false,
	}

	c := &http.Client{
		Transport: transport,
		Timeout:   time.Duration(options.Timeout) * time.Second,
	}

	parsedUrl, err := url.Parse(options.CheckUrl)
	if err != nil {
		return nil, err
	}

	return &remoteAuth{options: options, client: c, url: parsedUrl, cache: cache}, nil
}
