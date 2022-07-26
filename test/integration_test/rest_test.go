package integration_test

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/odpf/salt/log"

	"github.com/odpf/shield/core/resource"
	"github.com/odpf/shield/core/rule"
	"github.com/odpf/shield/internal/proxy"
	"github.com/odpf/shield/internal/proxy/hook"
	authz_hook "github.com/odpf/shield/internal/proxy/hook/authz"
	basic_auth "github.com/odpf/shield/internal/proxy/middleware/basic_auth"
	"github.com/odpf/shield/internal/proxy/middleware/prefix"
	"github.com/odpf/shield/internal/proxy/middleware/rulematch"
	"github.com/odpf/shield/internal/store/blob"
	"github.com/stretchr/testify/assert"

	"gocloud.dev/blob/fileblob"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const (
	restBackendPort = 13777
	restProxyPort   = restBackendPort + 1
	httpProtocol    = "http"
	h2cProtocol     = "h2c"
)

// @TODO: add tests for hooks

func TestREST(t *testing.T) {
	baseCtx, baseCancel := context.WithCancel(context.Background())
	defer baseCancel()

	blobFS, err := fileblob.OpenBucket("./fixtures", &fileblob.Options{
		CreateDir: true,
		Metadata:  fileblob.MetadataDontWrite,
	})
	if err != nil {
		t.Fatal(err)
	}

	responseHooks := hookPipeline(log.NewNoop())
	h2cProxy := proxy.NewH2c(proxy.NewH2cRoundTripper(log.NewNoop(), responseHooks), proxy.NewDirector())
	ruleRepo := blob.NewRuleRepository(log.NewNoop(), blobFS)
	if err := ruleRepo.InitCache(baseCtx, time.Minute); err != nil {
		t.Fatal(err)
	}
	defer ruleRepo.Close()
	ruleService := rule.NewService(ruleRepo)
	pipeline := buildPipeline(log.NewNoop(), h2cProxy, ruleService)

	proxyURL := fmt.Sprintf(":%d", restProxyPort)
	mux := http.NewServeMux()
	mux.Handle("/", pipeline)

	//create a tcp listener
	proxyListener, err := net.Listen("tcp", proxyURL)
	if err != nil {
		t.Fatal(err)
	}
	proxySrv := http.Server{
		Addr:    proxyURL,
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}
	defer proxySrv.Close()
	go func() {
		if err := proxySrv.Serve(proxyListener); err != nil && err != http.ErrServerClosed {
			t.Error(err)
		}
	}()

	for _, proto := range []string{httpProtocol, h2cProtocol} {
		func() {
			ts := startTestHTTPServer(restBackendPort, http.StatusOK, "", proto)
			defer ts.Close()

			// wait for proxy to start
			time.Sleep(time.Second * 1)
			t.Run("should handle GET request with 200", func(t *testing.T) {
				backendReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:%d/basic-authn", restProxyPort), nil)
				if err != nil {
					assert.Nil(t, err)
				}
				backendReq.SetBasicAuth("user", "password")
				resp, err := http.DefaultClient.Do(backendReq)
				if err != nil {
					assert.Nil(t, err)
				}
				assert.Equal(t, 200, resp.StatusCode)
				resp.Body.Close()
			})
			t.Run("should handle valid method request with 200", func(t *testing.T) {
				backendReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:%d/basic/", restProxyPort), nil)
				if err != nil {
					assert.Nil(t, err)
				}
				resp, err := http.DefaultClient.Do(backendReq)
				if err != nil {
					assert.Nil(t, err)
				}
				assert.Equal(t, 200, resp.StatusCode)
				resp.Body.Close()
			})
			t.Run("should handle invalid method request with 400", func(t *testing.T) {
				backendReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://127.0.0.1:%d/basic/", restProxyPort), nil)
				if err != nil {
					assert.Nil(t, err)
				}
				resp, err := http.DefaultClient.Do(backendReq)
				if err != nil {
					assert.Nil(t, err)
				}
				assert.Equal(t, 400, resp.StatusCode)
				resp.Body.Close()
			})
			t.Run("should give 401 if authn fails", func(t *testing.T) {
				backendReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:%d/basic-authn/", restProxyPort), nil)
				if err != nil {
					assert.Nil(t, err)
				}
				backendReq.SetBasicAuth("user", "XX")
				resp, err := http.DefaultClient.Do(backendReq)
				if err != nil {
					assert.Nil(t, err)
				}
				assert.Equal(t, 401, resp.StatusCode)
				resp.Body.Close()
			})
			t.Run("should give 401 if authz fails on json payload", func(t *testing.T) {
				buff := bytes.NewReader([]byte(`{"project": "xx"}`))
				backendReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://127.0.0.1:%d/basic-authz/", restProxyPort), buff)
				if err != nil {
					t.Fatal(err)
				}
				backendReq.SetBasicAuth("user", "password")
				resp, err := http.DefaultClient.Do(backendReq)
				if err != nil {
					assert.Nil(t, err)
				}
				assert.Equal(t, 401, resp.StatusCode)
				resp.Body.Close()
			})
			t.Run("should give 200 if authz success on json payload", func(t *testing.T) {
				buff := bytes.NewReader([]byte(`{"project": "foo"}`))
				backendReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://127.0.0.1:%d/basic-authz/", restProxyPort), buff)
				if err != nil {
					t.Fatal(err)
				}
				backendReq.SetBasicAuth("user", "password")
				resp, err := http.DefaultClient.Do(backendReq)
				if err != nil {
					assert.Nil(t, err)
				}
				assert.Equal(t, 200, resp.StatusCode)
				resp.Body.Close()
			})
		}()
	}
}

func BenchmarkProxyOverHttp(b *testing.B) {
	baseCtx, baseCancel := context.WithCancel(context.Background())
	defer baseCancel()

	blobFS, err := fileblob.OpenBucket("./fixtures", &fileblob.Options{
		CreateDir: true,
		Metadata:  fileblob.MetadataDontWrite,
	})
	if err != nil {
		b.Fatal(err)
	}

	h2cProxy := proxy.NewH2c(proxy.NewH2cRoundTripper(log.NewNoop(), hook.New()), proxy.NewDirector())
	ruleRepo := blob.NewRuleRepository(log.NewNoop(), blobFS)
	if err := ruleRepo.InitCache(baseCtx, time.Minute); err != nil {
		b.Fatal(err)
	}
	defer ruleRepo.Close()
	ruleService := rule.NewService(ruleRepo)
	pipeline := buildPipeline(log.NewNoop(), h2cProxy, ruleService)

	proxyURL := fmt.Sprintf(":%d", restProxyPort)
	mux := http.NewServeMux()
	mux.Handle("/", pipeline)

	//create a tcp listener
	proxyListener, err := net.Listen("tcp", proxyURL)
	if err != nil {
		b.Fatal(err)
	}
	proxySrv := http.Server{
		Addr:    proxyURL,
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}
	defer proxySrv.Close()
	go func() {
		if err := proxySrv.Serve(proxyListener); err != nil && err != http.ErrServerClosed {
			b.Error(err)
		}
	}()

	for _, proto := range []string{httpProtocol, h2cProtocol} {
		func() {
			ts := startTestHTTPServer(restBackendPort, http.StatusOK, "", proto)
			defer ts.Close()

			// wait for proxy to start
			time.Sleep(time.Second * 1)
			b.Run("200 status code GET on http1.1", func(b *testing.B) {
				backendReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:%d/basic/", restProxyPort), nil)
				if err != nil {
					b.Fatal(err)
				}
				for i := 0; i < b.N; i++ {
					resp, err := http.DefaultClient.Do(backendReq)
					if err != nil {
						panic(err)
					}
					if 200 != resp.StatusCode {
						b.Fatal("response code non 200")
					}
					resp.Body.Close()
				}
			})
			b.Run("200 status code with basic md5 authn on http1.1", func(b *testing.B) {
				backendReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:%d/basic-authn/", restProxyPort), nil)
				if err != nil {
					b.Fatal(err)
				}
				backendReq.SetBasicAuth("user", "password")
				for i := 0; i < b.N; i++ {
					resp, err := http.DefaultClient.Do(backendReq)
					if err != nil {
						panic(err)
					}
					if 200 != resp.StatusCode {
						b.Fatal("response code non 200")
					}
					resp.Body.Close()
				}
			})
			b.Run("200 status code with basic bcrypt authn on http1.1", func(b *testing.B) {
				backendReq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://127.0.0.1:%d/basic-authn-bcrypt/", restProxyPort), nil)
				if err != nil {
					b.Fatal(err)
				}
				backendReq.SetBasicAuth("user", "password")
				for i := 0; i < b.N; i++ {
					resp, err := http.DefaultClient.Do(backendReq)
					if err != nil {
						panic(err)
					}
					if 200 != resp.StatusCode {
						b.Fatal("response code non 200")
					}
					resp.Body.Close()
				}
			})
			b.Run("200 status code with basic authz on http1.1", func(b *testing.B) {
				buff := bytes.NewReader([]byte(`{"project": "foo"}`))
				backendReq, err := http.NewRequest(http.MethodPost, fmt.Sprintf("http://127.0.0.1:%d/basic-authz/", restProxyPort), buff)
				if err != nil {
					b.Fatal(err)
				}
				backendReq.SetBasicAuth("user", "password")
				for i := 0; i < b.N; i++ {
					resp, err := http.DefaultClient.Do(backendReq)
					if err != nil {
						panic(err)
					}
					if 200 != resp.StatusCode {
						b.Fatal("response code non 200")
					}
					resp.Body.Close()
				}
			})
		}()
	}
}

// buildPipeline builds middleware sequence
func buildPipeline(logger log.Logger, proxy http.Handler, ruleService *rule.Service) http.Handler {
	// Note: execution order is bottom up
	prefixWare := prefix.New(logger, proxy)
	//casbinAuthz := authz.New(logger, "", server.Deps{}, prefixWare)
	basicAuthn := basic_auth.New(logger, prefixWare)
	matchWare := rulematch.New(logger, basicAuthn, rulematch.NewRouteMatcher(ruleService))
	return matchWare
}

func hookPipeline(log log.Logger) hook.Service {
	rootHook := hook.New()
	return authz_hook.New(log, rootHook, rootHook, "", &resource.Service{})
}

func startTestHTTPServer(port, statusCode int, content, proto string) (ts *httptest.Server) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		if content != "" {
			_, err := w.Write([]byte(content))
			if err != nil {
				panic(err)
			}
		}
	})
	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
	if err != nil {
		panic(err)
	}

	var testHandler http.Handler = handler
	if proto == h2cProtocol {
		testHandler = h2c.NewHandler(handler, &http2.Server{})
	}

	ts = &httptest.Server{
		Listener: listener,
		Config: &http.Server{
			Handler:      testHandler,
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
			IdleTimeout:  time.Second,
		},
		EnableHTTP2: true,
	}
	ts.Start()
	return ts
}
