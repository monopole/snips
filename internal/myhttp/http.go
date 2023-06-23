package myhttp

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"
)

const (
	HeaderAccept              = "Accept"
	HeaderContentType         = "Content-Type"
	HeaderAAuthorization      = "Authorization"
	ContentTypeFormURLEncoded = "application/x-www-form-urlencoded"
	ContentTypeJson           = "application/json"
	Scheme                    = "https://"
)

// MakeHttpClient returns client ready to make HTTP requests.
// It's primed with certs loaded from the given caPath.
// If no caPath provided, TLS will be unauthenticated.
// The certs are used to establish that the servers are who they say they are.
func MakeHttpClient(caPath string) (*http.Client, error) {
	pool, err := loadCertPoolFromFile(caPath)
	if err != nil {
		return nil, err
	}
	jar, err := cookiejar.New(nil /* no options */)
	if err != nil {
		return nil, err
	}
	return &http.Client{
		Transport: makeTransport(makeTlsConfig(pool)),
		Timeout:   8 * time.Second,
		// Don't automatically follow redirects; we want debug mode to expose redirect hops.
		// CheckRedirect: func(req *http.Request, via []*http.Request) error {
		//	return http.ErrUseLastResponse
		// },
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

func makeTransport(tlsConfig *tls.Config) *http.Transport {
	return &http.Transport{
		TLSClientConfig:     tlsConfig,
		Proxy:               http.ProxyFromEnvironment,
		TLSHandshakeTimeout: 10 * time.Second,
	}
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

type PrArgs struct {
	Headers bool
	Body    bool
}

// PrintRequest prints a request.
func PrintRequest(r *http.Request, args PrArgs) error {
	fmt.Printf("==== Request to %s\n", r.URL)
	if args.Headers {
		for k, v := range r.Header {
			fmt.Printf("%50s : %s\n", k, v)
		}
	}
	if args.Body {
		const bodyDelim = "---------------------------------"
		// Doing this makes the body unavailable elsewhere.
		defer r.Body.Close()
		b, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		fmt.Println(bodyDelim)
		if body := strings.TrimSpace(string(b)); body != "" {
			fmt.Println(body)
		}
		fmt.Println(bodyDelim)
	}
	return nil
}

// PrintResponse prints a response.
func PrintResponse(r *http.Response, args PrArgs) error {
	fmt.Printf("==== Response Code %d ==================\n", r.StatusCode)
	if args.Headers {
		printResponseHeaders(r)
	}
	if args.Body {
		const bodyDelim = "---------------------------------"
		// Doing this makes the body unavailable elsewhere.
		defer r.Body.Close()
		b, err := io.ReadAll(r.Body)
		if err != nil {
			return err
		}
		fmt.Println(bodyDelim)
		if body := strings.TrimSpace(string(b)); body != "" {
			fmt.Println(body)
		}
		fmt.Println(bodyDelim)
	}
	return nil
}

// printResponseHeaders prints response headers.
func printResponseHeaders(r *http.Response) {
	for k, v := range r.Header {
		fmt.Printf("%50s : %s\n", k, v)
	}
}
