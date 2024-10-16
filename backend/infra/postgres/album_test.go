package postgres

import (
	"context"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/walnuts1018/mucaron/backend/domain/entity"
	"github.com/walnuts1018/mucaron/backend/domain/entity/gormmodel"
)

var _ = Describe("Album", Ordered, func() {
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

	album1 := entity.Album{
		UUIDModel: gormmodel.UUIDModel{
			ID: uuid.New(),
		},
		OwnerID: user1.ID,
		Name:    "album1",
		Musics: []entity.Music{
			music1,
		},
	}

	album2 := entity.Album{
		UUIDModel: gormmodel.UUIDModel{
			ID: uuid.New(),
		},
		OwnerID: user1.ID,
		Name:    "album2",
		Musics: []entity.Music{
			music1,
			music2,
		},
	}

	album2_dup := entity.Album{
		UUIDModel: gormmodel.UUIDModel{
			ID: uuid.New(),
		},
		OwnerID: user1.ID,
		Name:    "album2",
		Musics: []entity.Music{
			music1,
		},
	}

	ctx := context.Background()

	BeforeAll(func() {
		By("CreateUser")
		err := p.CreateUser(ctx, user1)
		Expect(err).NotTo(HaveOccurred())

		By("CreateMusic")
		err = p.CreateMusic(ctx, music1)
		Expect(err).NotTo(HaveOccurred())
		err = p.CreateMusic(ctx, music2)
		Expect(err).NotTo(HaveOccurred())
	})

	It("Album Normal CRUD", func() {
		By("CreateAlbum")
		err := p.CreateAlbum(ctx, album1)
		Expect(err).NotTo(HaveOccurred())
		Expect(album1.ID).NotTo(Equal(uuid.Nil))

		By("GetAlbum")
		a, err := p.GetAlbumByID(ctx, album1.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(a.ID).To(Equal(album1.ID))
		Expect(a.Name).To(Equal(album1.Name))
		Expect(a.Musics).To(HaveLen(1))

		By("UpdateAlbum")
		album1.Name = "album1_updated"
		err = p.UpdateAlbum(ctx, album1)
		Expect(err).NotTo(HaveOccurred())

		By("GetAlbum")
		a, err = p.GetAlbumByID(ctx, album1.ID)
		Expect(err).NotTo(HaveOccurred())
		Expect(a.ID).To(Equal(album1.ID))
		Expect(a.Name).To(Equal("album1_updated"))
		Expect(a.Musics).To(HaveLen(len(album1.Musics)))

		By("CreateAlbum 2")
		err = p.CreateAlbum(ctx, album2)
		Expect(err).NotTo(HaveOccurred())

		By("GetAlbumByIDs")
		albums, err := p.GetAlbumByIDs(ctx, []uuid.UUID{album1.ID, album2.ID})
		Expect(err).NotTo(HaveOccurred())
		Expect(albums).To(HaveLen(2))
		Expect(albums[0].ID).To(Equal(album1.ID))
		Expect(albums[1].ID).To(Equal(album2.ID))
		Expect(albums[1].Name).To(Equal(album2.Name))
		Expect(albums[1].Musics).To(HaveLen(len(album2.Musics)))

		By("DeleteAlbum")
		err = p.DeleteAlbums(ctx, []entity.Album{album1, album2})
		Expect(err).NotTo(HaveOccurred())

		By("GetAlbumByIDs")
		albums, err = p.GetAlbumByIDs(ctx, []uuid.UUID{album1.ID, album2.ID})
		Expect(err).NotTo(HaveOccurred())
		Expect(albums).To(BeEmpty())

		By("Music should not be deleted")
		musics, err := p.GetMusicByIDs(ctx, []uuid.UUID{music1.ID, music2.ID})
		Expect(err).NotTo(HaveOccurred())
		Expect(musics).To(HaveLen(2))
		Expect(musics[0].ID).To(Equal(music1.ID))
		Expect(musics[1].ID).To(Equal(music2.ID))
		Expect(musics[1].Name).To(Equal(music2.Name))
	})

	It("GetAlbumsByNameAndArtist", func() {
		By("Delete all albums")
		err := p.DB(ctx).Unscoped().Where(("deleted_at IS NOT NULL")).Delete(&entity.Album{}).Error
		Expect(err).NotTo(HaveOccurred())

		By("CreateAlbum")
		err = p.CreateAlbum(ctx, album1)
		Expect(err).NotTo(HaveOccurred())
		err = p.CreateAlbum(ctx, album2)
		Expect(err).NotTo(HaveOccurred())
		err = p.CreateAlbum(ctx, album2_dup)
		Expect(err).NotTo(HaveOccurred())

		By("Create Check")
		albums, err := p.GetAlbumByIDs(ctx, []uuid.UUID{album1.ID, album2.ID, album2_dup.ID})
		Expect(err).NotTo(HaveOccurred())
		Expect(albums).To(HaveLen(3))

		By("GetAlbumsByNameAndArtist")
		albums, err = p.GetAlbumsByNameAndArtist(ctx, user1.ID, album2.Name, artist2)
		Expect(err).NotTo(HaveOccurred())
		Expect(albums).To(HaveLen(1))
		Expect(albums[0].ID).To(Equal(album2.ID))
		Expect(albums[0].Name).To(Equal(album2.Name))
		Expect(albums[0].Musics).To(HaveLen(len(album2.Musics)))
	})
})
