package test

import (
	"path/filepath"
	"testing"

	devenv "github.com/kube-ianvs/kube-ianvs/pkg/env"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"k8s.io/utils/env"
	"sigs.k8s.io/controller-runtime/pkg/envtest/printer"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

var testEnv *devenv.Environment

// TestAPIs comment lint rebel
func TestAPIs(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecsWithDefaultAndCustomReporters(t,
		"Ianvs Gateway Suite",
		[]Reporter{printer.NewlineReporter{}})
}

var _ = BeforeSuite(func() {
	var err error
	logf.SetLogger(zap.New(zap.WriteTo(GinkgoWriter), zap.UseDevMode(true)))

	By("bootstrapping test environment")
	IanvsAssets := env.GetString("IANVS_ASSETS", "/data/projects/kube-ianvs/testbin/bin")
	// oidcIssuer := "http://127.0.0.1:15556/dex"
	testEnv = &devenv.Environment{
		Dex: devenv.Dex{
			ClientID: "ianvs",
			//	OidcIssuer:  oidcIssuer,
			RedirectURI: "http://localhost:9989",
		},
	}
	testEnv.CRDDirectoryPaths = []string{filepath.Join(IanvsAssets, "crd")}
	testEnv.ErrorIfCRDPathMissing = true

	err = testEnv.Start()

	Expect(err).NotTo(HaveOccurred())

}, 60)

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
})
