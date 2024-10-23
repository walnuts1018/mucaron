package usecase

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/walnuts1018/mucaron/backend/domain"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"github.com/walnuts1018/mucaron/backend/domain/entity/gormmodel"
	"go.uber.org/mock/gomock"
)

var _ = Describe("music.go", func() {
	user1 := entity.User{
		UUIDModel: gormmodel.UUIDModel{
			ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		},
	}

	user2 := entity.User{
		UUIDModel: gormmodel.UUIDModel{
			ID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		},
	}

	music1 := entity.Music{
		OwnerID: user1.ID,
		UUIDModel: gormmodel.UUIDModel{
			ID: uuid.MustParse("00000000-0000-0000-0001-000000000001"),
		},
	}

	Context("GetMusics", func() {
		It("Normal", func() {
			usecase, repos := NewMockUsecase()

			By("expect GetMusicsByUserID")
			repos.EntityRepository.EXPECT().GetMusicsByUserID(gomock.Any(), user1.ID).Return([]entity.Music{music1}, nil)

			By("get musics by user1")
			got, err := usecase.GetMusics(context.Background(), user1)
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal([]entity.Music{music1}))
		})
	})

	Context("GetMusicIDs", func() {
		It("Normal", func() {
			usecase, repos := NewMockUsecase()

			By("expect GetMusicIDsByUserID")
			repos.EntityRepository.EXPECT().GetMusicIDsByUserID(gomock.Any(), user1.ID).Return([]uuid.UUID{music1.ID}, nil)

			By("get music ids by user1")
			got, err := usecase.GetMusicIDs(context.Background(), user1)
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal([]uuid.UUID{music1.ID}))
		})
	})

	Context("DeleteMusics", func() {
		It("Normal", func() {
			usecase, repos := NewMockUsecase()

			By("expect GetMusicByIDs and DeleteMusics")
			repos.EntityRepository.EXPECT().GetMusicByIDs(gomock.Any(), []uuid.UUID{music1.ID}).Return([]entity.Music{music1}, nil)
			repos.EntityRepository.EXPECT().DeleteMusics(gomock.Any(), []uuid.UUID{music1.ID}).Return(nil)

			By("delete musics by user1")
			err := usecase.DeleteMusics(context.Background(), user1, []uuid.UUID{music1.ID})
			Expect(err).NotTo(HaveOccurred())
		})

		It("access denied", func() {
			usecase, repos := NewMockUsecase()

			By("expect GetMusicByIDs")
			repos.EntityRepository.EXPECT().GetMusicByIDs(gomock.Any(), []uuid.UUID{music1.ID}).Return([]entity.Music{music1}, nil)

			By("delete musics by user2")
			err := usecase.DeleteMusics(context.Background(), user2, []uuid.UUID{music1.ID})
			Expect(err).To(Equal(domain.ErrAccessDenied))
		})
	})
})
