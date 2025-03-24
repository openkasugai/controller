/*
Copyright 2025 NTT Corporation , FUJITSU LIMITED

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pkgruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	examplecomv1 "FPGAFunction/api/v1"
	controllertestcpu "FPGAFunction/internal/controller/test/type/CPU"
	controllertestethernet "FPGAFunction/internal/controller/test/type/Ethernet"
	controllertestgpu "FPGAFunction/internal/controller/test/type/GPU"
	controllertestpcie "FPGAFunction/internal/controller/test/type/PCIe"

	corev1 "k8s.io/api/core/v1"
	//+kubebuilder:scaffold:imports
)

// These tests use Ginkgo (BDD-style Go testing framework). Refer to
// http://onsi.github.io/ginkgo/ to learn more about Ginkgo.

var cfg *rest.Config
var k8sClient client.Client
var testEnv *envtest.Environment
var testScheme *pkgruntime.Scheme

var cfg2 *rest.Config
var k8sClient2 client.Client
var testEnv2 *envtest.Environment
var testScheme2 *pkgruntime.Scheme

func TestControllers(t *testing.T) {
	RegisterFailHandler(Fail)

	RunSpecs(t, "Controller Suite")
}

const (
	TESTNAMESPACE string = "default"
)

var _ = BeforeSuite(func() {

	By("bootstrapping test environment")
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "config", "crd", "bases"),
			filepath.Join("test", "crd"),
		},
		ErrorIfCRDPathMissing: true,

		// The BinaryAssetsDirectory is only required if you want to run the tests directly
		// without call the makefile target test. If not informed it will look for the
		// default path defined in controller-runtime which is /usr/local/kubebuilder/.
		// Note that you must have the required binaries setup under the bin directory to perform
		// the tests directly. When we run make test it will be setup and used automatically.
		BinaryAssetsDirectory: filepath.Join("..", "..", "bin", "k8s",
			fmt.Sprintf("1.28.0-%s-%s", runtime.GOOS, runtime.GOARCH)),
	}

	testEnv2 = &envtest.Environment{
		CRDDirectoryPaths: []string{
			filepath.Join("..", "..", "config", "crd", "bases"),
			filepath.Join("test", "crd"),
		},
		ErrorIfCRDPathMissing: true,

		// The BinaryAssetsDirectory is only required if you want to run the tests directly
		// without call the makefile target test. If not informed it will look for the
		// default path defined in controller-runtime which is /usr/local/kubebuilder/.
		// Note that you must have the required binaries setup under the bin directory to perform
		// the tests directly. When we run make test it will be setup and used automatically.
		BinaryAssetsDirectory: filepath.Join("..", "..", "bin", "k8s",
			fmt.Sprintf("1.28.0-%s-%s", runtime.GOOS, runtime.GOARCH)),
	}

	var err error
	// cfg is defined in this file globally.
	cfg, err = testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	testScheme = pkgruntime.NewScheme()

	err = examplecomv1.AddToScheme(testScheme)
	Expect(err).NotTo(HaveOccurred())

	err = controllertestcpu.AddToScheme(testScheme)
	Expect(err).NotTo(HaveOccurred())

	err = controllertestgpu.AddToScheme(testScheme)
	Expect(err).NotTo(HaveOccurred())

	err = controllertestethernet.AddToScheme(testScheme)
	Expect(err).NotTo(HaveOccurred())

	err = controllertestpcie.AddToScheme(testScheme)
	Expect(err).NotTo(HaveOccurred())

	err = corev1.AddToScheme(testScheme)
	Expect(err).NotTo(HaveOccurred())

	err = examplecomv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme
	k8sClient, err = client.New(cfg, client.Options{
		Scheme: testScheme,
		Cache:  &client.CacheOptions{},
	})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	// cfg is defined in this file globally.
	cfg2, err = testEnv2.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg2).NotTo(BeNil())

	testScheme2 = pkgruntime.NewScheme()

	err = examplecomv1.AddToScheme(testScheme2)
	Expect(err).NotTo(HaveOccurred())

	err = controllertestcpu.AddToScheme(testScheme2)
	Expect(err).NotTo(HaveOccurred())

	err = controllertestgpu.AddToScheme(testScheme2)
	Expect(err).NotTo(HaveOccurred())

	err = controllertestethernet.AddToScheme(testScheme2)
	Expect(err).NotTo(HaveOccurred())

	err = controllertestpcie.AddToScheme(testScheme2)
	Expect(err).NotTo(HaveOccurred())

	err = corev1.AddToScheme(testScheme2)
	Expect(err).NotTo(HaveOccurred())

	err = examplecomv1.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	//+kubebuilder:scaffold:scheme
	k8sClient2, err = client.New(cfg2, client.Options{
		Scheme: testScheme2,
		Cache:  &client.CacheOptions{},
	})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient2).NotTo(BeNil())
})

var _ = AfterSuite(func() {
	By("tearing down the test environment")
	err := testEnv.Stop()
	Expect(err).NotTo(HaveOccurred())
	err = testEnv2.Stop()
	Expect(err).NotTo(HaveOccurred())
})
