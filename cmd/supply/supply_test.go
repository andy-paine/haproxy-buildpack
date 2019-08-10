package supply_test

import (
	"bytes"
	"fmt"
	"github.com/andy-paine/haproxy-buildpack/cmd/supply"
	"github.com/cloudfoundry/libbuildpack"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io/ioutil"
	"os"
)

//go:generate mockgen -source=supply.go --destination=mocks_test.go --package=supply_test

var _ = Describe("Supply", func() {
	var (
		depDir        string
		supplier      *supply.Supplier
		logger        *libbuildpack.Logger
		mockCtrl      *gomock.Controller
		mockStager    *MockStager
		mockManifest  *MockManifest
		mockInstaller *MockInstaller
		mockCommand   *MockCommand
		buffer        *bytes.Buffer
	)

	BeforeEach(func() {
		var err error
		buffer = new(bytes.Buffer)
		logger = libbuildpack.NewLogger(buffer)

		mockCtrl = gomock.NewController(GinkgoT())
		mockStager = NewMockStager(mockCtrl)
		mockManifest = NewMockManifest(mockCtrl)
		mockInstaller = NewMockInstaller(mockCtrl)
		mockCommand = NewMockCommand(mockCtrl)

		Expect(err).ToNot(HaveOccurred())
		depDir, err = ioutil.TempDir("", "haproxy.depdir")
		mockStager.EXPECT().DepDir().AnyTimes().Return(depDir)
		mockManifest.EXPECT().AllDependencyVersions("haproxy").AnyTimes().Return([]string{"0.0.0", "0.0.1"})
		mockManifest.EXPECT().GetEntry(gomock.Any()).AnyTimes().Return(&libbuildpack.ManifestEntry{
      Dependency: libbuildpack.Dependency{ Version: "0.0.1" },
      URI: "https://downloaded.from/v0.0.1",
    }, nil)

		supplier = &supply.Supplier{
			Manifest:  mockManifest,
			Installer: mockInstaller,
			Stager:    mockStager,
			CompileCommand:   mockCommand,
			Log:       logger,
		}
	})

	AfterEach(func() {
		mockCtrl.Finish()
		os.RemoveAll(depDir)
	})

  It("should use the latest version available", func() {
    version, err := supplier.VersionToInstall()
		Expect(err).To(BeNil())
    Expect(version.Dependency.Version).To(Equal("0.0.1"))
  })

	It("should install HAProxy tarball", func() {
    dep := libbuildpack.Dependency{Name: "haproxy", Version: "0.0.1"}
		mockInstaller.EXPECT().InstallDependency(gomock.Eq(dep), fmt.Sprintf("%s/haproxy", depDir))
		Expect(supplier.InstallArchive(dep)).To(Succeed())
	})

  It("should compile HAProxy with correct flags", func() {
		mockCommand.EXPECT().Run(fmt.Sprintf("%s/haproxy", depDir))
    Expect(supplier.Compile()).To(Succeed())
  })
})
