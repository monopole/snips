package myjira

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"github.com/monopole/snips/internal/myhttp"
)

func (jb *jiraBoss) doJiraRequest(loc *url.URL, req issueSearchRequest) (resp *issueSearchResponse, err error) {
	var (
		ans  io.ReadCloser
		body []byte
	)
	body, err = json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("trouble marshaling data from request; %w", err)
	}
	ans, err = jb.sendPost(loc, bytes.NewBuffer(body), jb.args.Token)
	if err != nil {
		return nil, err
	}
	body, err = io.ReadAll(ans)
	if err != nil {
		return nil, fmt.Errorf("ReadAll failure: %w", err)
	}
	if debug := false; debug {
		var pretty bytes.Buffer
		_ = json.Indent(&pretty, body, "  ", "  ")
		fmt.Fprintln(os.Stderr, pretty.String())
	}
	resp = &issueSearchResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil, fmt.Errorf("trouble unmarshaling data from response; %w", err)
	}
	return
}

func (jb *jiraBoss) sendPost(loc *url.URL, body io.Reader, token string) (ans io.ReadCloser, err error) {
	const debug = false
	var (
		req  *http.Request
		resp *http.Response
	)
	req, err = http.NewRequest(http.MethodPost, loc.String(), body)
	if err != nil {
		return
	}
	req.Header.Set(myhttp.HeaderAccept, myhttp.ContentTypeJson)
	req.Header.Set(myhttp.HeaderContentType, myhttp.ContentTypeJson)
	req.Header.Set(myhttp.HeaderAAuthorization, fmt.Sprintf("Bearer: %s", token))
	resp, err = jb.htCl.Do(req)
	if err != nil {
		return
	}
	if debug {
		myhttp.PrintResponse(resp, myhttp.PrArgs{Headers: true, Body: true})
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("status code %d", resp.StatusCode)
		return
	}
	return resp.Body, nil
}
