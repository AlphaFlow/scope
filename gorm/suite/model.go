package suite

// Copied from https://github.com/gobuffalo/suite/blob/master/model.go
// Note: This replaces the SetupTest call to TruncateAll with a more performant option

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop/v5"
	buffaloSuite "github.com/gobuffalo/suite/v3"
	"github.com/gobuffalo/suite/v3/fix"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Model suite
type Model struct {
	suite.Suite
	*require.Assertions
	DB       *gorm.DB
	Fixtures buffaloSuite.Box
}

// DBDelta checks database table count change for a passed table name.
func (m *Model) DBDelta(delta int, name string, fn func()) {
	var sc, ec int64
	err := m.DB.Table(name).Count(&sc).Error
	m.NoError(err)
	fn()
	err = m.DB.Table(name).Count(&ec).Error
	m.NoError(err)
	m.Equal(sc+int64(delta), ec)
}

// LoadFixture loads a named fixture into the database.
func (m *Model) LoadFixture(name string) {
	sc, err := fix.Find(name)
	m.NoError(err)
	db := m.DB

	for _, table := range sc.Tables {
		for _, row := range table.Row {
			q := "insert into " + table.Name
			keys := []string{}
			skeys := []string{}
			for k := range row {
				keys = append(keys, k)
				skeys = append(skeys, "@"+k)
			}

			for k, v := range row {
				if arr, ok := v.([]interface{}); ok {
					asString := false
					if len(arr) > 0 {
						if _, ok := arr[0].(string); ok {
							asString = true
						}
					}
					b, err := json.Marshal(arr)
					if err != nil {
						m.NoError(err)
					}
					items := string(b)
					if asString {
						items = strings.ReplaceAll(items, "[\"", "{")
						items = strings.ReplaceAll(items, "\"]", "}")
						items = strings.ReplaceAll(items, "\",\"", ",")
						items = strings.ReplaceAll(items, "\", \"", ",")
					}
					row[k] = items
				}
			}

			q = q + fmt.Sprintf(" (%s) values (%s)", strings.Join(keys, ","), strings.Join(skeys, ","))
			err = db.Exec(q, row).Error
			m.NoError(err)
		}
	}
}

// NewModel creates a new model suite
func NewModel() *Model {
	m := &Model{}
	env := envy.Get("GO_ENV", "test")
	err := pop.LoadConfigFile()
	if err != nil {
		log.Fatalln(err)
	}
	popConn := pop.Connections[env]

	dsn := popConn.URL()
	c, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction:   true,
		DisableNestedTransaction: true,
		AllowGlobalUpdate:        false,
		FullSaveAssociations:     false,
		//Logger: logger.Default.LogMode(logger.Info),
	})
	if err == nil {
		m.DB = c
	} else {
		log.Fatal(err)
	}
	return m
}
