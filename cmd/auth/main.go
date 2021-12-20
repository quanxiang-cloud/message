package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	authi "github.com/quanxiang-cloud/message/pkg/auth"
	"github.com/quanxiang-cloud/message/pkg/auth/lowcode"
	"go.uber.org/zap"
)

var (
	log logr.Logger
)

func main() {
	var port string
	var entrypoint string
	var timeout time.Duration
	var keepAlive time.Duration
	var maxIdleConns int
	var idleConnTimeout time.Duration
	var tlsHandshakeTimeout time.Duration
	var expectContinueTimeout time.Duration

	flag.StringVar(&port, "port", ":40001", "service port default: :40001")
	flag.StringVar(&entrypoint, "entrypoint", "http://localhost:80", "service entrypoint default: http://localhost:80")
	flag.DurationVar(&timeout, "timeout", 20*time.Second, "Timeout is the maximum amount of time a dial will wait for a connect to complete. If Deadline is also set, it may fail earlier")
	flag.DurationVar(&keepAlive, "keep-alive", 20*time.Second, "KeepAlive specifies the interval between keep-alive probes for an active network connection.")
	flag.IntVar(&maxIdleConns, "max-idle-conns", 10, "MaxIdleConns controls the maximum number of idle (keep-alive) connections across all hosts. Zero means no limit.")
	flag.DurationVar(&idleConnTimeout, "idle-conn-timeout", 20*time.Second, "IdleConnTimeout is the maximum amount of time an idle (keep-alive) connection will remain idle before closing itself.")
	flag.DurationVar(&tlsHandshakeTimeout, "tls-handshake-timeout", 10*time.Second, "TLSHandshakeTimeout specifies the maximum amount of time waiting to wait for a TLS handshake. Zero means no timeout.")
	flag.DurationVar(&expectContinueTimeout, "expect-continue-timeout", 1*time.Second, "")
	flag.Parse()

	zapLog, err := zap.NewDevelopment()
	if err != nil {
		panic(fmt.Sprintf("who watches the watchmen (%v)?", err))
	}
	log = zapr.NewLogger(zapLog)

	uri, err := url.ParseRequestURI(entrypoint)
	if err != nil {
		panic(err)
	}

	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   timeout * time.Second,
			KeepAlive: keepAlive * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          maxIdleConns,
		IdleConnTimeout:       idleConnTimeout,
		TLSHandshakeTimeout:   tlsHandshakeTimeout,
		ExpectContinueTimeout: expectContinueTimeout * time.Second,
	}

	e := gin.New()
	e.Use(gin.Logger(), gin.Recovery())

	group := e.Group("/", auth(lowcode.NewLowcodeAuth(log)))
	group.Any("*path", proxy(uri, transport))

	log.Info("start...")
	e.Run(port)
}

func proxy(url *url.URL, transport *http.Transport) func(c *gin.Context) {
	return func(c *gin.Context) {
		proxy := httputil.NewSingleHostReverseProxy(url)
		proxy.Transport = transport
		r := c.Request
		r.Host = url.Host
		proxy.ServeHTTP(c.Writer, r)
	}
}

func auth(i authi.Interface) func(c *gin.Context) {
	return func(c *gin.Context) {
		if !i.Auth(c.Writer, c.Request) {
			c.Abort()
			return
		}
		c.Next()
	}
}
