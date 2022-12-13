package remauth

import (
	"compress/gzip"
	"crypto/md5"
	"errors"
	"fmt"
	"io"
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
	if r.options.Debug {
		defer timeTrack(time.Now(), "valid method")
	}

	cacheKey := fmt.Sprintf("%x", md5.Sum([]byte(token)))

	if valid, exists := r.cache.Get(cacheKey); exists {
		return valid.(bool)
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

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

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

func New(options *Options) (RemoteAuth, error) {
	ch := cache.New(options.CacheTime, 10*time.Minute)

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

	return &remoteAuth{options: options, client: c, url: parsedUrl, cache: ch}, nil
}

func (r *remoteAuth) Check(token string, result func(response interface{}, err error)) {
	if r.options.Debug {
		defer timeTrack(time.Now(), "check method")
	}

	cacheKey := fmt.Sprintf("obj-%x", md5.Sum([]byte(token)))

	if obj, exists := r.cache.Get(cacheKey); exists {
		result(obj, nil)
		return
	}

	req := &http.Request{
		Method: http.MethodGet,
		URL:    r.url,
		Header: http.Header{
			"User-Agent":      {"RemoteAuthClient/" + version},
			"Authorization":   {token},
			"Accept-Encoding": {"gzip"},
		},
	}

	if r.options.Debug {
		fmt.Printf("[RemoteAuthClient][URL]: %s\n", r.url.String())
	}

	resp, err := r.client.Do(req)
	if err != nil {
		result(nil, err)
		return
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	if resp.Header.Get("Content-Encoding") == "gzip" {
		resp.Body, err = gzip.NewReader(resp.Body)
		if err != nil {
			result(nil, err)
			return
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result(nil, err)
		return
	}

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusNoContent {
		r.cache.Set(cacheKey, body, 5*time.Minute)

		result(body, nil)
		return
	}

	r.cache.Set(cacheKey, nil, 5*time.Minute)
	result(nil, errors.New("invalid token"))
}
