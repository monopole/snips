package internal

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	headerAccept              = "Accept"
	headerContentType         = "Content-Type"
	contentTypeFormURLEncoded = "application/x-www-form-urlencoded"
	contentTypeJson           = "application/json"
	scheme                    = "https://"

	defaultWaitingInterval = 8 * time.Second
	maxAttempts            = 5
)

const verbose = false

type OAuthParams struct {
	Domain   string
	ClientId string
	// Don't appear to need a client secret for RO access.
}

type devCodeData struct {
	DeviceCode         string `json:"device_code"`
	UserCode           string `json:"user_code"`
	VerifyUri          string `json:"verification_uri"`
	ExpiresIn          int    `json:"expires_in"`
	MinIntervalSeconds int    `json:"interval"`
}

// GetOAuthAccessToken implements the instructions at
// https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps#device-flow
func GetOAuthAccessToken(params *OAuthParams) (accessCode string, err error) {
	cl, err := makeHttpClient("")
	if err != nil {
		return
	}
	loc, err := url.Parse(scheme + params.Domain + "/login/device/code")
	if err != nil {
		return
	}

	postData := url.Values{}
	postData.Add("client_id", params.ClientId)
	postData.Add("scope", "repo read:org user")
	r, err := sendPost(cl, loc, postData)
	if err != nil {
		return
	}

	var codes devCodeData
	p := json.NewDecoder(r)
	if err = p.Decode(&codes); err != nil {
		err = fmt.Errorf("failed to parse device code from %+v\n; %w", codes, err)
		return
	}
	if verbose {
		fmt.Printf("CODES %+v\n", codes)
	}
	waitingInterval := computeWaitingInterval(codes.MinIntervalSeconds)

	fmt.Fprintf(os.Stderr, `
  ***** You have %s to visit  %s  and enter this code:  %s
`, time.Duration(maxAttempts)*waitingInterval, codes.VerifyUri, codes.UserCode)

	return pollForTheUsersApproval(cl, params, waitingInterval, codes.DeviceCode)
}

func computeWaitingInterval(minIntervalSeconds int) time.Duration {
	minInterval := time.Duration(minIntervalSeconds) + time.Second
	if defaultWaitingInterval > minInterval {
		return defaultWaitingInterval
	}
	return minInterval
}

func pollForTheUsersApproval(
	cl *http.Client, params *OAuthParams,
	pollInterval time.Duration, devCode string) (accessCode string, err error) {
	loc, err := url.Parse(scheme + params.Domain + "/login/oauth/access_token")
	if err != nil {
		return
	}

	postData := url.Values{}
	postData.Add("client_id", params.ClientId)
	postData.Add("device_code", devCode)
	postData.Add("grant_type", "urn:ietf:params:oauth:grant-type:device_code")

	attempts := 0
	for attempts < maxAttempts {
		time.Sleep(pollInterval)
		var r io.ReadCloser
		r, err = sendPost(cl, loc, postData)
		if err != nil {
			if verbose {
				fmt.Printf("attempt %2d: post error %s\n", attempts, err.Error())
			}
			continue
		}
		var token struct {
			Value string `json:"access_token"`
			Type  string `json:"token_type"`
			Scope string `json:"scope"`
		}
		p := json.NewDecoder(r)
		if err = p.Decode(&token); err != nil {
			if verbose {
				fmt.Printf("attempt %2d: response parse error %s\n", attempts, err.Error())
			}
			continue
		}
		if verbose {
			fmt.Printf("TOKEN %+v\n", token)
		}
		if token.Value != "" {
			accessCode = token.Value
			return
		}
		attempts++
	}
	err = fmt.Errorf("exhausted %d attempts", maxAttempts)
	return
}

func sendPost(cl *http.Client, loc *url.URL, postData url.Values) (ans io.ReadCloser, err error) {
	var (
		req  *http.Request
		resp *http.Response
	)
	if verbose {
		fmt.Printf("requesting to %s\n", loc.String())
	}
	req, err = http.NewRequest(http.MethodPost, loc.String(), bytes.NewBufferString(postData.Encode()))
	if err != nil {
		return
	}
	req.Header.Set(headerAccept, contentTypeJson)
	req.Header.Set(headerContentType, contentTypeFormURLEncoded)
	resp, err = cl.Do(req)
	if err != nil {
		return
	}
	if verbose {
		printResponse(resp)
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code %d", resp.StatusCode)
		return
	}
	return resp.Body, nil
}

// makeHttpClient returns client ready to make HTTP requests.
// It's primed with certs loaded from the given path.
// If no path provided, TLS will be unauthenticated.
// The certs are used to establish that the servers are who they say they are.
func makeHttpClient(path string) (*http.Client, error) {
	pool, err := loadCertPoolFromFile(path)
	if err != nil {
		return nil, err
	}
	jar, err := cookiejar.New(nil /* no options */)
	if err != nil {
		return nil, err
	}
	return &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: makeTlsConfig(pool),
		},
		// Don't automatically follow redirects; we want to see what DS is doing.
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
		Jar: jar,
	}, nil
}

// loadCertPoolFromFile returns a pool containing the certs read from the given file.
func loadCertPoolFromFile(path string) (*x509.CertPool, error) {
	if path == "" {
		return nil, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("unable to load certs from %q; %w", path, err)
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(data)
	return pool, nil
}

// makeTlsConfig returns a TLS config that uses a cert pool if provided, falling
// back to no cert check if no cert pool provided.
func makeTlsConfig(pool *x509.CertPool) *tls.Config {
	if pool == nil {
		return &tls.Config{
			InsecureSkipVerify: true,
		}
	}
	return &tls.Config{
		RootCAs: pool,
	}
}

// printHeaders prints response headers.
func printHeaders(r *http.Response) {
	for k, v := range r.Header {
		fmt.Printf("%50s : %s\n", k, v)
	}
}

// printResponse prints a response.
func printResponse(r *http.Response) error {
	fmt.Printf("==== %d ==================\n", r.StatusCode)
	if showHeaders := true; showHeaders {
		printHeaders(r)
	}
	if readBody := false; readBody {
		// Doing this makes the body unavailable elsewhere.
		defer r.Body.Close()
		b, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		if body := strings.TrimSpace(string(b)); body != "" {
			fmt.Println(body)
		}
	}
	fmt.Printf("============================\n\n")
	return nil
}
