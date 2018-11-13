package store

import (
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

				fixture := models.Scheme{
					Version: 1,
					Tags:    []string{"a", "b", "c"},
				}

				s := &schemes{db: tx}

				err := s.Create(&fixture)
				g.Expect(err).NotTo(HaveOccurred())

				items, err := s.Search(&SearchRequest{
					Version: fixture.Version,
					Tags:    fixture.Tags,
				})

				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(items).To(HaveLen(1))

				return rollback
			},
		},

		// Read:
		{
			name: "read",
			handler: func(g *GomegaWithT, tx orm.DB) error {

				fixture := models.Scheme{
					ID:      1,
					Version: 1,
					Tags:    []string{"a", "b", "c"},
				}

				s := &schemes{db: tx}

				err := s.Create(&fixture)
				g.Expect(err).NotTo(HaveOccurred())

				item, err := s.Read(fixture.ID)

				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(item).NotTo(BeNil())

				return rollback
			},
		},

		// Update:
		{
			name: "update",
			handler: func(g *GomegaWithT, tx orm.DB) error {

				fixture := models.Scheme{
					Version: 1,
					Tags:    []string{"a", "b", "c"},
				}

				s := &schemes{db: tx}

				err := s.Create(&fixture)
				g.Expect(err).NotTo(HaveOccurred())

				err = s.Update(&fixture)
				g.Expect(err).NotTo(HaveOccurred())

				items, err := s.Search(&SearchRequest{
					Tags: fixture.Tags,
				})

				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(items).To(HaveLen(2))

				g.Expect(items[0].Version).To(BeEquivalentTo(1))
				g.Expect(items[1].Version).To(BeEquivalentTo(2))

				return rollback
			},
		},

		// Delete:
		{
			name: "delete",
			handler: func(g *GomegaWithT, tx orm.DB) error {

				fixture := models.Scheme{
					ID:      1,
					Version: 1,
					Tags:    []string{"a", "b", "c"},
				}

				s := &schemes{db: tx}

				err := s.Create(&fixture)
				g.Expect(err).NotTo(HaveOccurred())

				err = s.Delete(&fixture)
				g.Expect(err).NotTo(HaveOccurred())

				items, err := s.Search(&SearchRequest{
					Tags: fixture.Tags,
				})

				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(items).To(HaveLen(0))

				return rollback
			},
		},

		// Search:
		{
			name: "search",
			handler: func(g *GomegaWithT, tx orm.DB) error {
				fixtures := []*models.Scheme{
					{
						Version: 1,
						Tags:    []string{"a", "b", "c"},
					},

					{
						Version: 2,
						Tags:    []string{"b", "c", "d"},
					},

					{
						Version: 3,
						Tags:    []string{"c", "d", "e"},
					},
				}

				_, err := tx.Model(&fixtures).Insert()
				g.Expect(err).NotTo(HaveOccurred())

				s := &schemes{db: tx}

				items, err := s.Search(&SearchRequest{
					Tags: []string{"c"},
				})

				g.Expect(err).NotTo(HaveOccurred())

				g.Expect(items).To(HaveLen(len(fixtures)))

				return rollback
			},
		},
	}
)

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
