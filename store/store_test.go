package store

import (
	"encoding/json"

	"github.com/go-pg/pg"
	"github.com/im-kulikov/helium"
	"github.com/im-kulikov/helium/grace"
	"github.com/im-kulikov/helium/logger"
	"github.com/im-kulikov/helium/module"
	"github.com/im-kulikov/helium/orm"
	"github.com/im-kulikov/helium/redis"
	"github.com/im-kulikov/helium/settings"
	"github.com/im-kulikov/simplinic-task/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var testModule = module.Module{}.Append(
	grace.Module,
	settings.Module,
	logger.Module,
	redis.Module,
	orm.Module,
)

var _ = Describe("Store Suite", func() {
	var db *pg.DB

	BeforeSuite(func() {
		var err error

		//err = os.Setenv("TEST_POSTGRES_DEBUG", "false")
		//Expect(err).NotTo(HaveOccurred())
		//
		//err = os.Setenv("TEST_LOGGER_LEVEL", "info")
		//Expect(err).NotTo(HaveOccurred())

		h, err := helium.New(&helium.Settings{
			File:   "../config.yml",
			Prefix: "TEST",
		}, testModule)
		Expect(err).NotTo(HaveOccurred())

		Expect(h.Invoke(func(pdb *pg.DB) {
			db = pdb

			// Cleanup all tables...
			_, err := db.Exec("TRUNCATE schemes RESTART IDENTITY CASCADE;")
			Expect(err).NotTo(HaveOccurred())
		})).NotTo(HaveOccurred())
	})

	AfterSuite(func() {

		Expect(db.Close()).NotTo(HaveOccurred())
	})

	Context("try CRUD+S of schemes", func() {
		var (
			s       Schemes
			fixture Scheme
		)

		BeforeEach(func() {
			fixture = Scheme{
				Version: 1,
				Tags:    []string{"a", "b", "c"},
				Data:    json.RawMessage(`{"hello": "world"}`),
			}

			s = NewSchemeStore(db)
		})

		It("should create scheme without errors", func() {
			err := s.Create(&fixture)
			Expect(err).NotTo(HaveOccurred())

			var items []*models.SchemeVersion

			err = db.Model(&items).Where("scheme_id = ?",
				fixture.ID).Select()

			Expect(err).NotTo(HaveOccurred())
			Expect(items).To(HaveLen(1))
			Expect(items[0].Version).To(Equal(fixture.Version))
			Expect(items[0].Tags).To(Equal(fixture.Tags))
			Expect(items[0].Data).To(Equal(fixture.Data))
		})

		It("should read created scheme without errors", func() {
			err := s.Create(&fixture)
			Expect(err).NotTo(HaveOccurred())

			item, err := s.Read(fixture.ID)

			Expect(err).NotTo(HaveOccurred())
			Expect(*item).To(Equal(fixture))
		})

		It("should update created scheme without errors", func() {
			err := s.Create(&fixture)
			Expect(err).NotTo(HaveOccurred())

			err = s.Update(&fixture)
			Expect(err).NotTo(HaveOccurred())

			var items []*Scheme

			err = db.Model(&items).Where("scheme_id = ?", fixture.ID).Select()
			Expect(err).NotTo(HaveOccurred())
			Expect(items).To(HaveLen(2))

			for i, item := range items {
				Expect(item.Version).To(BeEquivalentTo(i + 1))
				Expect(item.Tags).To(BeEquivalentTo(fixture.Tags))
				Expect(item.Data).To(BeEquivalentTo(fixture.Data))
			}
		})

		It("should delete created scheme without errors", func() {
			err := s.Create(&fixture)
			Expect(err).NotTo(HaveOccurred())

			err = s.Delete(fixture.ID)
			Expect(err).NotTo(HaveOccurred())

			item, err := s.Read(fixture.ID)
			Expect(err).To(HaveOccurred()) // not found
			Expect(item).To(BeNil())
		})

		It("should search created schemes without errors", func() {
			fixtures := []Scheme{
				{
					Tags: []string{"a1", "b1", "c1"},
				},

				{
					Tags: []string{"b1", "c1", "d1"},
				},

				{
					Tags: []string{"c1", "d1", "e1"},
				},
			}

			var ids []int64

			for _, item := range fixtures {
				item.Version = 1
				err := s.Create(&item)
				Expect(err).NotTo(HaveOccurred())

				ids = append(ids, item.ID)
			}

			items, err := s.Search(SearchRequest{
				Tags: []string{"c1"},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(items).To(HaveLen(len(fixtures)))

			// Must ignore deleted schemes:
			for i, id := range ids {
				err = s.Delete(id)
				Expect(err).NotTo(HaveOccurred())

				items, err := s.Search(SearchRequest{
					Tags: []string{"c1"},
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(items).To(HaveLen(len(fixtures) - 1 - i))
			}
		})
	})

	Context("try CRUD+S of configs", func() {
		var (
			s       Configs
			fixture Config
			scheme  Scheme
		)

		BeforeEach(func() {

			scheme = Scheme{
				Version: 1,
				Tags:    []string{"a", "b", "c"},
				Data:    json.RawMessage(`{"hello": "world"}`),
			}

			err := NewSchemeStore(db).Create(&scheme)
			Expect(err).NotTo(HaveOccurred())

			fixture = Config{
				SchemeID: scheme.ID,
				Version:  1,
				Tags:     []string{"a", "b", "c"},
				Data:     json.RawMessage(`{"hello": "world"}`),
			}

			s = NewConfigStore(db)
		})

		It("should create config without errors", func() {
			err := s.Create(&fixture)
			Expect(err).NotTo(HaveOccurred())

			var items []*Config

			err = db.Model(&items).Where("scheme_id = ?", scheme.ID).Select()

			Expect(err).NotTo(HaveOccurred())
			Expect(items).To(HaveLen(1))
			Expect(items[0].Version).To(BeEquivalentTo(fixture.Version))
			Expect(items[0].Tags).To(BeEquivalentTo(fixture.Tags))
			Expect(items[0].Data).To(BeEquivalentTo(fixture.Data))
		})

		It("should read created config without errors", func() {
			err := s.Create(&fixture)
			Expect(err).NotTo(HaveOccurred())

			item, err := s.Read(fixture.ID)
			Expect(err).NotTo(HaveOccurred())

			Expect(item.Version).To(BeEquivalentTo(fixture.Version))
			Expect(item.Tags).To(BeEquivalentTo(fixture.Tags))
			Expect(item.Data).To(BeEquivalentTo(fixture.Data))
		})

		It("should update created config without errors", func() {
			err := s.Create(&fixture)
			Expect(err).NotTo(HaveOccurred())

			err = s.Update(&fixture)
			Expect(err).NotTo(HaveOccurred())

			var items []*Config
			err = db.Model(&items).
				Where("scheme_id = ? AND config_id = ?", scheme.ID, fixture.ID).
				Order("version DESC").
				Select()

			Expect(err).NotTo(HaveOccurred())
			Expect(items).To(HaveLen(2))

			for i := 0; i < 2; i++ {
				Expect(items[i].Version).To(BeEquivalentTo(fixture.Version-int64(i)), "version %d", i)
				Expect(items[i].Tags).To(BeEquivalentTo(fixture.Tags), "tags %d", i)
				Expect(items[i].Data).To(BeEquivalentTo(fixture.Data), "data %d", i)
			}
		})

		It("should delete created config without errors", func() {
			err := s.Create(&fixture)
			Expect(err).NotTo(HaveOccurred())

			err = s.Delete(fixture.ID)
			Expect(err).NotTo(HaveOccurred())

			item, err := s.Read(fixture.ID)
			Expect(err).To(HaveOccurred()) // not found
			Expect(item).To(BeNil())
		})

		It("should search created configs without errors", func() {
			fixtures := []*Config{
				{SchemeID: scheme.ID, Tags: []string{"a2", "b2", "c2"}},
				{SchemeID: scheme.ID, Tags: []string{"b2", "c2", "d2"}},
				{SchemeID: scheme.ID, Tags: []string{"c2", "d2", "e2"}},
			}

			var ids []int64

			for _, item := range fixtures {
				err := s.Create(item)
				Expect(err).NotTo(HaveOccurred())

				ids = append(ids, item.ID)
			}

			items, err := s.Search(SearchRequest{
				Tags: []string{"c2"},
			})

			Expect(err).NotTo(HaveOccurred())
			Expect(items).To(HaveLen(len(fixtures)))

			// Must ignore deleted schemes:
			for i, id := range ids {
				err = s.Delete(id)
				Expect(err).NotTo(HaveOccurred())

				items, err := s.Search(SearchRequest{
					Tags: []string{"c2"},
				})

				Expect(err).NotTo(HaveOccurred())
				Expect(items).To(HaveLen(len(fixtures)-1-i), "step %d / %d", i+1, len(ids))
			}
		})

	})
})
