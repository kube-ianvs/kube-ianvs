package envtest

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"html/template"
	"io/fs"
	"io/ioutil"
	"time"

	"k8s.io/client-go/rest"
	"k8s.io/utils/env"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
)

var (
	IanvsAssets = env.GetString("IANVS_ASSETS", "/data/projects/kube-ianvs/testbin/bin")
)

// Environment comment lint rebel
type Environment struct {
	envtest.Environment
	Dex                   Dex
	APIServerToken        string
	KubeConfigTestingPath string
	RawIDToken            string
	Config                *rest.Config
}

// Start comment lint rebel
func (e *Environment) Start() (err error) {
	err = e.beforeStart()
	if err != nil {
		return
	}

	e.Config, err = e.Environment.Start()
	if err != nil {
		return
	}
	err = e.Dex.Start()
	if err != nil {
		return
	}
	err = e.afterStart()
	return
}

// SetDefault comment lint rebel
func (e *Environment) beforeStart() (err error) {
	e.KubeAPIServerFlags = DefaultKubeAPIServerFlags()
	e.Environment.BinaryAssetsDirectory = IanvsAssets

	md5Ctx := md5.New()
	md5Ctx.Write([]byte(time.Now().String()))
	e.APIServerToken = fmt.Sprintf("%x", md5Ctx.Sum(nil))

	err = e.InitTokenCsv()
	return
}

func (e *Environment) afterStart() (err error) {

	err = e.InitKubeConfigFile()
	if err != nil {
		return
	}

	e.RawIDToken, err = e.Dex.GetIDToken()
	if err != nil {
		return
	}

	return
}

// InitTokenCsv comment lint rebel
func (e *Environment) InitTokenCsv() (err error) {
	var templateData []byte
	templateData, err = ioutil.ReadFile(IanvsAssets + "/token.csv.template")
	if err != nil {
		return
	}

	templateText := string(templateData)
	var t *template.Template
	t, err = template.New(templateText).Parse(templateText)
	if err != nil {
		return
	}

	buf := &bytes.Buffer{}
	err = t.Execute(buf, e)
	if err != nil {
		return
	}

	e.KubeConfigTestingPath = fmt.Sprintf("%s/token.csv", IanvsAssets)
	err = ioutil.WriteFile(e.KubeConfigTestingPath, buf.Bytes(), fs.ModePerm)
	return
}

// InitKubeConfigFile comment lint rebel
func (e *Environment) InitKubeConfigFile() (err error) {
	var configTemplate []byte
	configTemplate, err = ioutil.ReadFile(IanvsAssets + "/config.template")
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
	err = t.Execute(buf, e)
	if err != nil {
		return
	}

	e.KubeConfigTestingPath = fmt.Sprintf("%s/config.testing", IanvsAssets)
	err = ioutil.WriteFile(e.KubeConfigTestingPath, buf.Bytes(), fs.ModePerm)
	return
}

// DefaultKubeAPIServerFlags comment lint rebel
func DefaultKubeAPIServerFlags() []string {
	return []string{
		//"--advertise-address=127.0.0.1",
		"--etcd-servers={{ if .EtcdURL }}{{ .EtcdURL.String }}{{ end }}",
		"--cert-dir={{ .CertDir }}",
		"--insecure-port={{ if .URL }}{{ .URL.Port }}{{ end }}",
		"--insecure-bind-address={{ if .URL }}{{ .URL.Hostname }}{{ end }}",
		"--secure-port={{ if .SecurePort }}{{ .SecurePort }}{{ end }}",
		// we're keeping this disabled because if enabled, default SA is missing which would force all tests to create one
		// in normal apiserver operation this SA is created by controller, but that is not run in integration environment
		"--disable-admission-plugins=ServiceAccount",
		//"--service-cluster-ip-range=10.0.0.0/24",
		"--allow-privileged=true",
		fmt.Sprintf("--token-auth-file=%s/token.csv", IanvsAssets),
	}
}
