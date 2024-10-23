package usecase

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/walnuts1018/mucaron/backend/domain"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"github.com/walnuts1018/mucaron/backend/domain/entity/gormmodel"
	newuuid "github.com/walnuts1018/mucaron/backend/util/new_uuid"
	"github.com/walnuts1018/mucaron/backend/util/random"
	"go.uber.org/mock/gomock"
)

var _ = Describe("user.go", func() {
	salt := "salt"

	user1Password := entity.RawPassword("password1")
	user1 := entity.User{
		UUIDModel: gormmodel.UUIDModel{
			ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		},
	}

	user2Password := entity.RawPassword("password2")
	user2 := entity.User{
		UUIDModel: gormmodel.UUIDModel{
			ID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		},
	}

	BeforeEach(func() {
		By("seed用乱数固定")
		random.SetTestValue([]byte(salt))

		user1LoginInfo, err := entity.NewLoginInfo(user1Password)
		Expect(err).NotTo(HaveOccurred())
		user1.LoginInfo = user1LoginInfo

		user2LoginInfo, err := entity.NewLoginInfo(user2Password)
		Expect(err).NotTo(HaveOccurred())
		user2.LoginInfo = user2LoginInfo
	})

	Context("Login", func() {
		It("Normal", func() {
			usecase, repos := NewMockUsecase()

			By("expect GetUserByName")
			repos.EntityRepository.EXPECT().GetUserByName(gomock.Any(), user1.UserName).Return(user1, nil)

			By("login user1")
			got, err := usecase.Login(context.Background(), user1.UserName, user1Password)
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(user1))
		})

		It("UserNotFound", func() {
			usecase, repos := NewMockUsecase()

			By("expect GetUserByName")
			repos.EntityRepository.EXPECT().GetUserByName(gomock.Any(), user1.UserName).Return(entity.User{}, domain.ErrNotFound)

			By("login user1")
			_, err := usecase.Login(context.Background(), user2.UserName, user2Password)
			Expect(err).To(Equal(ErrUserNotFound))
		})

		It("IncorrectPW", func() {
			usecase, repos := NewMockUsecase()

			By("expect GetUserByName")
			repos.EntityRepository.EXPECT().GetUserByName(gomock.Any(), user1.UserName).Return(user1, nil)

			By("login user1 with incorrect password")
			_, err := usecase.Login(context.Background(), user1.UserName, user2Password)
			Expect(err).To(Equal(ErrIncorrectPW))
		})
	})

	Context("CreateUser", func() {
		It("Normal", func() {
			usecase, repos := NewMockUsecase()

			By("expect same user not found")
			repos.EntityRepository.EXPECT().GetUserByName(gomock.Any(), user1.UserName).Return(entity.User{}, domain.ErrNotFound)

			By("expect CreateUser")
			repos.EntityRepository.EXPECT().CreateUser(gomock.Any(), user1).Return(nil)

			By("UUID 固定")
			newuuid.SetUUIDValue(user1.ID)

			By("seed用乱数固定")
			random.SetTestValue([]byte(salt))

			By("create user1")
			got, err := usecase.CreateUser(context.Background(), user1.UserName, user1Password)
			Expect(err).NotTo(HaveOccurred())
			Expect(got).To(Equal(user1))
		})

		It("UserExists", func() {
			usecase, repos := NewMockUsecase()

			By("expect same user found")
			repos.EntityRepository.EXPECT().GetUserByName(gomock.Any(), user1.UserName).Return(user1, nil)

			By("create user1")
			_, err := usecase.CreateUser(context.Background(), user1.UserName, user1Password)
			Expect(err).To(Equal(ErrUserExists))
		})
	})

})
