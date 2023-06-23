package oauth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/monopole/snips/internal/myhttp"
)

const (
	scopes = "repo read:org user"

	defaultWaitingInterval = 8 * time.Second
	maxAttempts            = 10
	WarningPrefix          = " ***** "
)

type Params struct {
	// GhDomain is the GitHub domain, likely "github.com", or "github.company.com"
	GhDomain string
	// ClientId for this program as registered at that domain.
	// We don't appear to need a client secret for RO access.
	ClientId string
	// HttpCl is used to communicate with the GhDomain.
	HttpCl *http.Client
	// Verbose being true yields more print statements.
	Verbose bool
}

type devCodeData struct {
	DeviceCode         string `json:"device_code"`
	UserCode           string `json:"user_code"`
	VerifyUri          string `json:"verification_uri"`
	ExpiresIn          int    `json:"expires_in"`
	MinIntervalSeconds int    `json:"interval"`
}

// GetAccessToken implements the instructions at
// https://docs.github.com/en/apps/oauth-apps/building-oauth-apps/authorizing-oauth-apps#device-flow
func GetAccessToken(params *Params) (accessCode string, err error) {
	loc, err := url.Parse(myhttp.Scheme + params.GhDomain + "/login/device/code")
	if err != nil {
		return
	}

	postData := url.Values{}
	postData.Add("client_id", params.ClientId)
	postData.Add("scope", scopes)
	if params.Verbose {
		fmt.Printf("posting to %s\n", loc.String())
	}
	r, err := sendPost(params.HttpCl, loc, postData, params.Verbose)
	if err != nil {
		return
	}

	var codes devCodeData
	p := json.NewDecoder(r)
	if err = p.Decode(&codes); err != nil {
		err = fmt.Errorf("failed to parse device code from %+v\n; %w", codes, err)
		return
	}
	if params.Verbose {
		fmt.Printf("CODES %+v\n", codes)
	}
	waitingInterval := computeWaitingInterval(codes.MinIntervalSeconds)
	warnUser(waitingInterval, &codes)
	return pollForTheUsersApproval(params.HttpCl, params, waitingInterval, codes.DeviceCode)
}

func warnUser(waitingInterval time.Duration, codes *devCodeData) {
	fmt.Fprintf(os.Stderr, `
%s You have %s to visit  %s  and enter this code:  %s
`, WarningPrefix, time.Duration(maxAttempts)*waitingInterval, codes.VerifyUri, codes.UserCode)
}

func computeWaitingInterval(minIntervalSeconds int) time.Duration {
	minInterval := time.Duration(minIntervalSeconds) + time.Second
	if defaultWaitingInterval > minInterval {
		return defaultWaitingInterval
	}
	return minInterval
}

func pollForTheUsersApproval(
	cl *http.Client, params *Params,
	pollInterval time.Duration, devCode string) (accessCode string, err error) {
	loc, err := url.Parse(myhttp.Scheme + params.GhDomain + "/login/oauth/access_token")
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
		r, err = sendPost(cl, loc, postData, params.Verbose)
		if err != nil {
			if params.Verbose {
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
			if params.Verbose {
				fmt.Printf("attempt %2d: response parse error %s\n", attempts, err.Error())
			}
			continue
		}
		if params.Verbose {
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

func sendPost(cl *http.Client, loc *url.URL, postData url.Values, debug bool) (ans io.ReadCloser, err error) {
	var (
		req  *http.Request
		resp *http.Response
	)
	req, err = http.NewRequest(http.MethodPost, loc.String(), bytes.NewBufferString(postData.Encode()))
	if err != nil {
		return
	}
	req.Header.Set(myhttp.HeaderAccept, myhttp.ContentTypeJson)
	req.Header.Set(myhttp.HeaderContentType, myhttp.ContentTypeFormURLEncoded)
	resp, err = cl.Do(req)
	if err != nil {
		return
	}
	if debug {
		myhttp.PrintResponse(resp, myhttp.PrArgs{Headers: true})
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code %d", resp.StatusCode)
		return
	}
	return resp.Body, nil
}
