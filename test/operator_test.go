package test

import (
	"context"
	"fmt"

	v1 "github.com/kube-ianvs/ianvs-operator/apis/cluster/v1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var mockCredentialName = "mock-credential"

var _ = Describe("Ianvs operator", func() {
	var err error

	It("should added test credential", func() {

		err = v1.AddToScheme(scheme.Scheme)

		Expect(err).NotTo(HaveOccurred())

		k8sClient, err = client.New(testEnv.Config, client.Options{Scheme: scheme.Scheme})
		Expect(err).NotTo(HaveOccurred())

		var testc v1.Credential
		testc.Name = mockCredentialName
		testc.Spec.Cluster.InsecureSkipTLSVerify = true
		testc.Spec.Cluster.Server = fmt.Sprintf("https://127.0.0.1:%d", testEnv.ControlPlane.APIServer.SecurePort)
		testc.Spec.AuthInfo.Token = testEnv.APIServerToken

		err = k8sClient.Create(context.TODO(), &testc)

		Expect(err).NotTo(HaveOccurred())

		Expect(k8sClient).NotTo(BeNil())

	})
})
