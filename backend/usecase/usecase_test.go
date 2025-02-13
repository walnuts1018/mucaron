package usecase

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/walnuts1018/mucaron/backend/config"
	mock_usecase "github.com/walnuts1018/mucaron/backend/usecase/mock_usecase"
	"go.uber.org/mock/gomock"
)

type mockRepostitories struct {
	EntityRepository *mock_usecase.MockEntityRepository
	MetadataReader   *mock_usecase.MockMetadataReader
	ObjectStorage    *mock_usecase.MockObjectStorage
}

func NewMockUsecase() (*Usecase, mockRepostitories) {
	cfg := config.Config{}
	entityRepository := mock_usecase.NewMockEntityRepository(ctrl)
	metadataReader := mock_usecase.NewMockMetadataReader(ctrl)
	objectStorage := mock_usecase.NewMockObjectStorage(ctrl)

	return NewUsecase(cfg, entityRepository, metadataReader, objectStorage),
		mockRepostitories{
			EntityRepository: entityRepository,
			MetadataReader:   metadataReader,
			ObjectStorage:    objectStorage,
		}
}

var ctrl *gomock.Controller

func TestUsecase(t *testing.T) {
	RegisterFailHandler(Fail)

	ctrl = gomock.NewController(t)

	RunSpecs(t, "Postgres Suite")
}
