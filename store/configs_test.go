package store

import (
	"encoding/json"
	"testing"

	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"github.com/im-kulikov/helium"
	. "github.com/onsi/gomega"
)

var (
	testConfigSuites = []struct {
		name    string
		handler func(g *GomegaWithT, tx orm.DB) error
	}{
		// Create:
		{
			name: "create",
			handler: func(g *GomegaWithT, tx orm.DB) error {

				scheme := Scheme{
					Version: 1,
					Tags:    []string{"a", "b", "c"},
					Data:    json.RawMessage(`{"hello": "world"}`),
				}

				err := (&schemes{db: tx}).Create(&scheme)
				g.Expect(err).NotTo(HaveOccurred())

				fixture := Config{
					SchemeID: scheme.ID,
					Version:  1,
					Tags:     []string{"a", "b", "c"},
					Data:     json.RawMessage(`{"hello": "world"}`),
				}

				s := &configs{db: tx}

				err = s.Create(&fixture)
				g.Expect(err).NotTo(HaveOccurred())

				var items []*Config

				err = tx.Model(&items).Where("scheme_id = ?", scheme.ID).Select()

				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(items).To(HaveLen(1))
				g.Expect(items[0].Version).To(BeEquivalentTo(fixture.Version))
				g.Expect(items[0].Tags).To(BeEquivalentTo(fixture.Tags))
				g.Expect(items[0].Data).To(BeEquivalentTo(fixture.Data))

				return rollback
			},
		},

		// Read:
		{
			name: "read",
			handler: func(g *GomegaWithT, tx orm.DB) error {

				scheme := Scheme{
					Version: 1,
					Tags:    []string{"a", "b", "c"},
					Data:    json.RawMessage(`{"hello": "world"}`),
				}

				err := (&schemes{db: tx}).Create(&scheme)
				g.Expect(err).NotTo(HaveOccurred())

				fixture := Config{
					SchemeID: scheme.ID,
					Version:  1,
					Tags:     []string{"a", "b", "c"},
					Data:     json.RawMessage(`{"hello": "world"}`),
				}

				s := &configs{db: tx}

				err = s.Create(&fixture)
				g.Expect(err).NotTo(HaveOccurred())

				item, err := s.Read(fixture.ID)
				g.Expect(err).NotTo(HaveOccurred())

				g.Expect(item.Version).To(BeEquivalentTo(fixture.Version))
				g.Expect(item.Tags).To(BeEquivalentTo(fixture.Tags))
				g.Expect(item.Data).To(BeEquivalentTo(fixture.Data))

				return rollback
			},
		},

		// Update:
		{
			name: "update",
			handler: func(g *GomegaWithT, tx orm.DB) error {

				scheme := Scheme{
					Version: 1,
					Tags:    []string{"a", "b", "c"},
					Data:    json.RawMessage(`{"hello": "world"}`),
				}

				err := (&schemes{db: tx}).Create(&scheme)
				g.Expect(err).NotTo(HaveOccurred())

				fixture := Config{
					SchemeID: scheme.ID,
					Version:  1,
					Tags:     []string{"a", "b", "c"},
					Data:     json.RawMessage(`{"hello": "world"}`),
				}

				s := &configs{db: tx}

				err = s.Create(&fixture)
				g.Expect(err).NotTo(HaveOccurred())

				err = s.Update(&fixture)
				g.Expect(err).NotTo(HaveOccurred())

				var items []*Config
				err = tx.Model(&items).
					Where("scheme_id = ? AND config_id = ?", scheme.ID, fixture.ID).
					Order("version DESC").
					Select()
				g.Expect(err).NotTo(HaveOccurred())
				g.Expect(items).To(HaveLen(2))

				for i := 0; i < 2; i++ {
					g.Expect(items[i].Version).To(BeEquivalentTo(fixture.Version-int64(i)), "version %d", i)
					g.Expect(items[i].Tags).To(BeEquivalentTo(fixture.Tags), "tags %d", i)
					g.Expect(items[i].Data).To(BeEquivalentTo(fixture.Data), "data %d", i)
				}

				return rollback
			},
		},

		// Delete:
		{
			name: "delete",
			handler: func(g *GomegaWithT, tx orm.DB) error {

				scheme := Scheme{
					Version: 1,
					Tags:    []string{"a", "b", "c"},
					Data:    json.RawMessage(`{"hello": "world"}`),
				}

				err := (&schemes{db: tx}).Create(&scheme)
				g.Expect(err).NotTo(HaveOccurred())

				fixture := Config{
					SchemeID: scheme.ID,
					Version:  1,
					Tags:     []string{"a", "b", "c"},
					Data:     json.RawMessage(`{"hello": "world"}`),
				}

				s := &configs{db: tx}

				err = s.Create(&fixture)
				g.Expect(err).NotTo(HaveOccurred())

				err = s.Delete(fixture.ID)
				g.Expect(err).NotTo(HaveOccurred())

				item, err := s.Read(fixture.ID)
				g.Expect(err).To(HaveOccurred()) // not found
				g.Expect(item).To(BeNil())

				return rollback
			},
		},

		// Search:
		{
			name: "search",
			handler: func(g *GomegaWithT, tx orm.DB) error {
				scheme := Scheme{
					Version: 1,
					Tags:    []string{"a", "b", "c"},
					Data:    json.RawMessage(`{"hello": "world"}`),
				}

				err := (&schemes{db: tx}).Create(&scheme)
				g.Expect(err).NotTo(HaveOccurred())

				fixtures := []*Config{
					{
						SchemeID: scheme.ID,
						Version:  1,
						Tags:     []string{"a", "b", "c"},
					},

					{
						SchemeID: scheme.ID,
						Version:  2,
						Tags:     []string{"b", "c", "d"},
					},

					{
						SchemeID: scheme.ID,
						Version:  3,
						Tags:     []string{"c", "d", "e"},
					},
				}

				s := &configs{
					db: tx,
				}

				var ids []int64

				for _, item := range fixtures {
					item.Version = 1
					err := s.Create(item)
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

func TestConfigs(t *testing.T) {
	g := NewGomegaWithT(t)

	h, err := helium.New(&helium.Settings{
		File:   "../config.yml",
		Prefix: "TEST",
	}, testModule)

	g.Expect(err).NotTo(HaveOccurred())

	g.Expect(h.Invoke(func(db *pg.DB) {

		for _, suite := range testConfigSuites {
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
