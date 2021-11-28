package envtest

import (
	"bytes"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"net/http"
	"net/url"
	"os/exec"
	"path"
	"sync"
	"time"

	"github.com/kube-ianvs/kube-ianvs/pkg/addr"
	"github.com/kube-ianvs/kube-ianvs/pkg/oidc"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

// Dex comment lint rebel
type Dex struct {
	cmd           *exec.Cmd
	ClientID      string
	OidcIssuer    string
	RedirectURI   string
	ConfigPath    string
	Host          string
	Port          int
	TelemetryPort int
}

// InitKubeConfigFile comment lint rebel
func (d *Dex) InitConfigFile() (err error) {
	var configTemplate []byte
	configTemplate, err = ioutil.ReadFile(IanvsAssets + "/dex-config.yaml.template")
	if err != nil {
		return
	}

	configTemplateString := string(configTemplate)
	var t *template.Template
	t, err = template.New(configTemplateString).Parse(configTemplateString)
	if err != nil {
		return
	}

	buf := &bytes.Buffer{}
	err = t.Execute(buf, d)
	if err != nil {
		fmt.Println(err)
		return
	}

	configPath := fmt.Sprintf("%s/dex-config.yaml.testing", IanvsAssets)
	err = ioutil.WriteFile(configPath, buf.Bytes(), fs.ModePerm)
	return
}

// Start comment lint rebel
func (d *Dex) Start() (err error) {
	d.Port, d.Host, err = addr.Suggest(d.Host, 10)
	if err != nil {
		return
	}
	d.TelemetryPort, _, err = addr.Suggest(d.Host, 10)
	if err != nil {
		return
	}
	d.InitConfigFile()
	d.OidcIssuer = fmt.Sprintf("http://%s:%d/dex", d.Host, d.Port)

	d.cmd = exec.Command(IanvsAssets+"/dex", "serve", IanvsAssets+"/dex-config.yaml.testing")
	outBuffer := bytes.NewBufferString("")
	errBuffer := bytes.NewBufferString("")
	d.cmd.Stdout = outBuffer
	d.cmd.Stderr = errBuffer
	go func() {
		_ = d.cmd.Run()
	}()
	ready := make(chan bool)
	timedOut := time.After(30 * time.Second)
	healthCheckURL, _ := url.Parse(d.OidcIssuer + "/healthz")
	HealthCheckPollInterval := 100 * time.Millisecond
	pollerStopCh := make(StopChannel)
	go PollURLUntilOK(nil, *healthCheckURL, HealthCheckPollInterval, ready, pollerStopCh)

	select {
	case <-ready:
		return nil
	case <-timedOut:
		if pollerStopCh != nil {
			close(pollerStopCh)
		}
		d.Stop()
		return fmt.Errorf(
			"timeout waiting for process %s to start \nerrBuffer:%s \noutBuffer:%s",
			path.Base(d.cmd.Path),
			errBuffer.String(),
			outBuffer.String(),
		)
	}
}

// Stop comment lint rebel
func (d *Dex) Stop() {
	_ = d.cmd.Process.Kill()
}

// NewHTTPClientFetcher comment lint rebel
func NewHTTPClientFetcher() (c oidc.CodeFetcher) {
	return &HTTPClientFetcher{}
}

// HTTPClientFetcher comment lint rebel
type HTTPClientFetcher struct {
	oidc.CodeFetcher
}

// CallbackServer comment lint rebel
type CallbackServer struct {
	svr  http.Server
	Code chan string
}

var callbackServer *CallbackServer
var callbackServerOnce sync.Once

// getCallbackServer comment lint rebel
func getCallbackServer() *CallbackServer {
	callbackServerOnce.Do(func() {
		var err error
		callbackServer = &CallbackServer{
			svr: http.Server{
				Addr: ":9989",
			},
			Code: make(chan string),
		}

		go func() {
			if err = callbackServer.svr.ListenAndServe(); err != http.ErrServerClosed {
				logrus.Fatalf("ListenAndServe(): %v", err)
			}
		}()

		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			code := r.URL.Query().Get("code")
			callbackServer.Code <- code
		})
	})
	return callbackServer
}

// Fetch comment lint rebel
func (f *HTTPClientFetcher) Fetch(a *oidc.App) (code string) {
	cs := getCallbackServer()
	go func() {
		res, _ := http.DefaultClient.Get(a.GetLoginURL())
		defer func() {
			_ = res.Body.Close()
		}()
	}()

	code = <-cs.Code
	return
}

// GetIDToken comment lint rebel
func (d *Dex) GetIDToken() (rawIDToken string, err error) {
	var oidcApp oidc.App
	oidcApp, err = oidc.Setup(
		"ianvs",
		d.OidcIssuer,
		"http://localhost:9989",
	)
	if err != nil {
		return
	}
	codeFetcher := NewHTTPClientFetcher()
	code := codeFetcher.Fetch(&oidcApp)
	var token *oauth2.Token
	token, err = oidcApp.Exchange(code)
	if err != nil {
		return
	}
	rawIDToken, _ = token.Extra("id_token").(string)
	return
}

// StopChannel comment lint rebel
type StopChannel chan struct{}

// PollURLUntilOK comment lint rebel
func PollURLUntilOK(httpClient *http.Client, url url.URL, interval time.Duration, ready chan bool, stopCh StopChannel) {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	if interval <= 0 {
		interval = 100 * time.Millisecond
	}
	for {
		res, err := httpClient.Get(url.String())
		if err == nil {
			_ = res.Body.Close()
			if res.StatusCode == http.StatusOK {
				ready <- true
				return
			}
		}

		select {
		case <-stopCh:
			return
		default:
			time.Sleep(interval)
		}
	}
}
