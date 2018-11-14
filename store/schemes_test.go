package store

import (
	"encoding/json"
	"testing"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/im-kulikov/helium"
	"github.com/im-kulikov/simplinic-task/models"
	. "github.com/onsi/gomega"
	"github.com/pkg/errors"
)

var (
	rollback         = errors.New("rollback")
	testSchemeSuites = []struct {
		name    string
		handler func(g *GomegaWithT, tx orm.DB) error
	}{
		// Create:
		{
			name: "create",
			handler: func(g *GomegaWithT, tx orm.DB) error {

				fixture := Scheme{
					Version: 1,
					Tags:    []string{"a", "b", "c"},
					Data:    json.RawMessage(`{"hello": "world"}`),
				}

				s := &schemes{db: tx}

				err := s.Create(&fixture)
				g.Expect(err).NotTo(HaveOccurred())

				var items []*models.SchemeVersion

				err = tx.Model(&items).Where("scheme_id = ?",
					fixture.ID).Select()

				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(items).To(HaveLen(1))
				g.Expect(items[0].Version).To(Equal(fixture.Version))
				g.Expect(items[0].Tags).To(Equal(fixture.Tags))
				g.Expect(items[0].Data).To(Equal(fixture.Data))

				return rollback
			},
		},

		// Read:
		{
			name: "read",
			handler: func(g *GomegaWithT, tx orm.DB) error {

				fixture := Scheme{
					Version: 1,
					Tags:    []string{"a", "b", "c"},
					Data:    json.RawMessage(`{"hello": "world"}`),
				}

				s := &schemes{db: tx}

				err := s.Create(&fixture)
				g.Expect(err).NotTo(HaveOccurred())

				item, err := s.Read(fixture.ID)

				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(*item).To(Equal(fixture))

				return rollback
			},
		},

		// Update:
		{
			name: "update",
			handler: func(g *GomegaWithT, tx orm.DB) error {

				fixture := Scheme{
					Version: 1,
					Tags:    []string{"a", "b", "c"},
					Data:    json.RawMessage(`{"hello": "world"}`),
				}

				s := &schemes{db: tx}

				err := s.Create(&fixture)
				g.Expect(err).NotTo(HaveOccurred())

				err = s.Update(&fixture)
				g.Expect(err).NotTo(HaveOccurred())

				var items []*Scheme

				err = tx.Model(&items).Where("scheme_id = ?", fixture.ID).Select()
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(items).To(HaveLen(2))

				for i, item := range items {
					g.Expect(item.Version).To(BeEquivalentTo(i + 1))
					g.Expect(item.Tags).To(BeEquivalentTo(fixture.Tags))
					g.Expect(item.Data).To(BeEquivalentTo(fixture.Data))
				}

				return rollback
			},
		},

		// Delete:
		{
			name: "delete",
			handler: func(g *GomegaWithT, tx orm.DB) error {

				fixture := Scheme{
					Version: 1,
					Tags:    []string{"a", "b", "c"},
					Data:    json.RawMessage(`{"hello": "world"}`),
				}

				s := &schemes{db: tx}

				err := s.Create(&fixture)
				g.Expect(err).NotTo(HaveOccurred())

				err = s.Delete(fixture.ID)
				g.Expect(err).NotTo(HaveOccurred())

				var items []*Scheme

				err = tx.Model(&items).Where("scheme_id = ?", fixture.ID).Select()
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(items).To(HaveLen(1))

				for i, item := range items {
					g.Expect(item.Version).To(BeEquivalentTo(i + 1))
					g.Expect(item.Tags).To(BeEquivalentTo(fixture.Tags))
					g.Expect(item.Data).To(BeEquivalentTo(fixture.Data))
				}

				return rollback
			},
		},

		// Search:
		{
			name: "search",
			handler: func(g *GomegaWithT, tx orm.DB) error {
				fixtures := []Scheme{
					{
						Tags: []string{"a", "b", "c"},
					},

					{
						Tags: []string{"b", "c", "d"},
					},

					{
						Tags: []string{"c", "d", "e"},
					},
				}
				s := &schemes{db: tx}

				var ids []int64

				for _, item := range fixtures {
					item.Version = 1
					err := s.Create(&item)
					g.Expect(err).NotTo(HaveOccurred())

					ids = append(ids, item.ID)
				}

				items, err := s.Search(SearchRequest{
					Tags: []string{"c"},
				})

				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(items).To(HaveLen(len(fixtures)))

				// Must ignore deleted schemes:
				for i, id := range ids {
					err = s.Delete(id)
					g.Expect(err).NotTo(HaveOccurred())

					items, err := s.Search(SearchRequest{
						Tags: []string{"c"},
					})

					g.Expect(err).NotTo(HaveOccurred())
					g.Expect(items).To(HaveLen(len(fixtures) - 1 - i))
				}

				return rollback
			},
		},
	}
)

func TestSchemes_Create(t *testing.T) {
	g := NewGomegaWithT(t)

	h, err := helium.New(&helium.Settings{
		File:   "../config.yml",
		Prefix: "TEST",
	}, testModule)

	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(h.Invoke(func(db *pg.DB) {

		err := db.RunInTransaction(func(tx *pg.Tx) error {
			s := &schemes{db: tx}
			err := s.Create(&Scheme{
				Version: 0,
				Tags:    []string{"a", "b", "c"},
				Data:    json.RawMessage(`{"hello": "world"}`),
			})

			g.Expect(err).NotTo(HaveOccurred())

			var model Scheme

			err = tx.Model(&model).Limit(1).Select()
			g.Expect(err).NotTo(HaveOccurred())

			return rollback
		})

		g.Expect(err).To(BeEquivalentTo(rollback))

	})).NotTo(HaveOccurred())
}

func TestSchemes(t *testing.T) {
	g := NewGomegaWithT(t)

	h, err := helium.New(&helium.Settings{
		File:   "../config.yml",
		Prefix: "TEST",
	}, testModule)

	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(h.Invoke(func(db *pg.DB) {

		for _, suite := range testSchemeSuites {
			t.Run(suite.name, func(t *testing.T) {
				g := NewGomegaWithT(t)

				err := db.RunInTransaction(func(tx *pg.Tx) error {
					return suite.handler(g, tx)
				})

				g.Expect(err).To(BeEquivalentTo(rollback))
			})
		}

	})).NotTo(HaveOccurred())
}
