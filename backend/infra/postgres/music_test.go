package postgres

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/walnuts1018/mucaron/backend/domain"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"github.com/walnuts1018/mucaron/backend/domain/entity/gormmodel"
)

var _ = Describe("Music", Ordered, func() {
	user1 := entity.User{
		UUIDModel: gormmodel.UUIDModel{
			ID: uuid.New(),
		},
		UserName: "user1",
	}

	artist1 := entity.Artist{
		UUIDModel: gormmodel.UUIDModel{
			ID: uuid.New(),
		},
		OwnerID: user1.ID,
		Name:    "artist1",
	}

	artist2 := entity.Artist{
		UUIDModel: gormmodel.UUIDModel{
			ID: uuid.New(),
		},
		OwnerID: user1.ID,
		Name:    "artist2",
	}

	music1 := entity.Music{
		UUIDModel: gormmodel.UUIDModel{
			ID: uuid.New(),
		},
		OwnerID: user1.ID,
		Name:    "music1",
		Artists: []entity.Artist{
			artist1,
		},
		FileHash: "filehash1",
	}

	music2 := entity.Music{
		UUIDModel: gormmodel.UUIDModel{
			ID: uuid.New(),
		},
		OwnerID: user1.ID,
		Name:    "music2",
		Artists: []entity.Artist{
			artist1,
			artist2,
		},
		FileHash: "filehash2",
	}

	ctx := context.Background()

	BeforeAll(func() {
		By("CreateUser")
		Expect(p.CreateUser(ctx, user1)).To(Succeed())
		Expect(p.CreateArtist(ctx, artist1)).To(Succeed())
	})

	It("test musics CRUD", func() {
		By("create")
		Expect(p.CreateMusic(ctx, music1)).To(Succeed())
		Expect(p.CreateMusic(ctx, music2)).To(Succeed())

		By("get")
		m1, err := p.GetMusicByID(ctx, music1.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(m1.ID).To(Equal(music1.ID))
		Expect(m1.Name).To(Equal(music1.Name))
		Expect(m1.Artists).To(HaveLen(1))
		Expect(m1.Artists[0].ID).To(Equal(artist1.ID))
		Expect(m1.Artists[0].Name).To(Equal(artist1.Name))
		Expect(m1.FileHash).To(Equal(music1.FileHash))

		m2, err := p.GetMusicByID(ctx, music2.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(m2.ID).To(Equal(music2.ID))
		Expect(m2.Name).To(Equal(music2.Name))
		Expect(m2.Artists).To(HaveLen(2))
		Expect(m2.Artists[0].ID).To(Equal(artist1.ID))
		Expect(m2.Artists[0].Name).To(Equal(artist1.Name))
		Expect(m2.Artists[1].ID).To(Equal(artist2.ID))
		Expect(m2.Artists[1].Name).To(Equal(artist2.Name))
		Expect(m2.FileHash).To(Equal(music2.FileHash))

		By("update")
		m1.Name = "music1_updated"
		Expect(p.UpdateMusic(ctx, m1)).To(Succeed())
		m1Updated, err := p.GetMusicByID(ctx, m1.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(m1Updated.ID).To(Equal(m1.ID))
		Expect(m1Updated.Name).To(Equal("music1_updated"))

		By("delete")
		Expect(p.DeleteMusics(ctx, []uuid.UUID{music1.ID})).To(Succeed())
		_, err = p.GetMusicByID(ctx, music1.ID)
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(domain.ErrNotFound))

		By("create music1 again")
		err = p.CreateMusic(ctx, music1)
		Expect(err).To(HaveOccurred())

		By("hard delete")
		Expect(p.HardDeleteMusic(ctx, music1)).To(Succeed())
		_, err = p.GetMusicByID(ctx, music1.ID)
		Expect(err).To(HaveOccurred())
		Expect(err).To(MatchError(domain.ErrNotFound))

		By("create music1 again")
		Expect(p.CreateMusic(ctx, music1)).To(Succeed())
		newM1, err := p.GetMusicByID(ctx, music1.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(newM1.DeletedAt.Valid).To(BeFalse())
	})
})
