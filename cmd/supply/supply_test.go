package supply_test

import (
  "bytes"
	"io/ioutil"
	"os"
  "github.com/andy-paine/haproxy-buildpack/cmd/supply"
  "github.com/cloudfoundry/libbuildpack"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

    supplier = &supply.Supplier{
      Manifest: mockManifest,
      Installer: mockInstaller,
      Stager:   mockStager,
      Command:  mockCommand,
      Log:      logger,
    }
	})

  AfterEach(func() {
		mockCtrl.Finish()
		os.RemoveAll(depDir)
	})

	It("should download HAProxy tarball", func() {
		Expect(supplier.Run()).To(Succeed())
	})
})
