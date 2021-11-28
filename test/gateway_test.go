package test

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/kube-ianvs/kube-ianvs/cmd"
	"github.com/kube-ianvs/kube-ianvs/pkg/addr"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var k8sClient client.Client
var serverURL = "https://localhost:18085"

var _ = Describe("ianvs gateway api", func() {
	var httpClient *http.Client

	BeforeEach(func() {
		gatewayCmd := cmd.NewGatewayCmd()
		port, host, err := addr.Suggest("", 10)
		if err != nil {
			Expect(err).NotTo(HaveOccurred())
			return
		}

		serverURL = fmt.Sprintf("https://%s:%d", host, port)
		gatewayCmd.SetArgs([]string{
			fmt.Sprintf("--addr=%s:%d", host, port),
			fmt.Sprintf("--certFile=%s/apiserver.crt", testEnv.ControlPlane.APIServer.CertDir),
			fmt.Sprintf("--keyFile=%s/apiserver.key", testEnv.ControlPlane.APIServer.CertDir),
		})
		go func() {
			gatewayCmd.Execute()
		}()

		time.Sleep(5 * time.Second)

		httpClient = oauth2.NewClient(context.TODO(), oauth2.StaticTokenSource(&oauth2.Token{
			AccessToken: testEnv.RawIDToken,
			TokenType:   "Bearer",
		}))
		httpClient.
			Transport.(*oauth2.Transport).
			Base = http.DefaultTransport
		httpClient.
			Transport.(*oauth2.Transport).
			Base.(*http.Transport).
			TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	})

	It("check healthz should response status code 200", func() {
		resp, err := httpClient.Get(serverURL + "/healthz")
		Expect(err).NotTo(HaveOccurred())
		Expect(resp.StatusCode).To(BeEquivalentTo(http.StatusOK))
	})

})
