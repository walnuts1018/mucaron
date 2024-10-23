package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/walnuts1018/mucaron/backend/domain/entity/gormmodel"
)

var _ = Describe("Transaction", func() {
	ctx := context.Background()

	to1 := testObject{
		UUIDModel: gormmodel.UUIDModel{
			ID: uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		},
		Name: "test1",
	}
	to2 := testObject{
		UUIDModel: gormmodel.UUIDModel{
			ID: uuid.MustParse("00000000-0000-0000-0000-000000000002"),
		},
		Name: "test2",
	}

	It("Test Transaction", func() {
		By("auto migrate")
		err := p.DB(ctx).AutoMigrate(&testObject{})
		Expect(err).ShouldNot(HaveOccurred())

		By("should commit")
		err = p.Transaction(ctx, func(ctx context.Context) error {
			By("create test object 1")
			err := p.createTestObject(ctx, to1)
			Expect(err).ShouldNot(HaveOccurred())

			By("get test object 1")
			to1copy, err := p.getTestObjectByID(ctx, to1.ID)
			Expect(err).ShouldNot(HaveOccurred())
			Expect(to1copy.Name).To(Equal(to1.Name))

			By("update test object 1")
			to1copy.Name = "test1copy"
			err = p.updateTestObject(ctx, to1copy)
			Expect(err).ShouldNot(HaveOccurred())

			By("get test object 1 copy")
			to1copy, err = p.getTestObjectByID(ctx, to1.ID)
			Expect(err).ShouldNot(HaveOccurred())
			return nil
		})
		Expect(err).ShouldNot(HaveOccurred())

		By("get test object 1 copy")
		to1copy, err := p.getTestObjectByID(ctx, to1.ID)
		Expect(err).ShouldNot(HaveOccurred())
		Expect(to1copy.Name).To(Equal("test1copy"))

		By("should rollback")
		err = p.Transaction(ctx, func(ctx context.Context) error {
			By("create test object 2")
			err := p.createTestObject(ctx, to2)
			Expect(err).ShouldNot(HaveOccurred())

			By("return error")
			return errors.New("should rollback")
		})
		Expect(err).Should(HaveOccurred())

		By("test object 2 should not exist")
		to2copy, err := p.getTestObjectByID(ctx, to2.ID)
		Expect(err).Should(HaveOccurred())
		Expect(to2copy).To(Equal(testObject{}))
	})
})
