package sqlstore

import (
	"github.com/masterhung0112/hk_server/mlog"
	"github.com/masterhung0112/hk_server/model"
	"github.com/masterhung0112/hk_server/store"
	"github.com/masterhung0112/hk_server/store/storetest"
	"github.com/stretchr/testify/suite"
	"os"
	"sync"
	"testing"
)

type StoreTestBaseSuite interface {
	suite.TestingSuite

	InitInitializeStore()

	SetStore(store store.Store)
	Store() store.Store

	SetSqlStore(sqlStore *SqlStore)
	SqlStore() *SqlStore
}

type StoreTestSuite struct {
	// suite.Suite

	store    store.Store
	sqlStore *SqlStore
	// sqlSupplier storetest.SqlSupplier
}

/***
 * StoreTestSuite implements interface StoreTestBaseSuite
 ***/
func (s *StoreTestSuite) InitInitializeStore() {
	if len(StoreTypes) >= 1 && (s.Store() == nil || s.SqlStore() == nil) {
		s.SetStore(StoreTypes[0].Store)
		s.SetSqlStore(StoreTypes[0].SqlStore)
	}
}

func (s *StoreTestSuite) SetStore(store store.Store) {
	s.store = store
}

func (s *StoreTestSuite) Store() store.Store {
	return s.store
}

func (s *StoreTestSuite) SetSqlStore(sqlStore *SqlStore) {
	s.sqlStore = sqlStore
}

func (s *StoreTestSuite) SqlStore() *SqlStore {
	return s.sqlStore
}

type storeType struct {
	Name        string
	SqlSettings *model.SqlSettings
	SqlStore    *SqlStore
	Store       store.Store
}

var StoreTypes []*storeType
var tearDownStoresOnce sync.Once

func initStores() {
	if testing.Short() {
		return
	}

	// In CI, we already run the entire test suite for both mysql and postgres in parallel.
	// So we just run the tests for the current database set.
	if os.Getenv("IS_CI") == "true" {
		panic("Not implement IS_CI yet")
	} else {
		StoreTypes = append(StoreTypes, &storeType{
			Name:        "MySQL",
			SqlSettings: storetest.MakeSqlSettings(model.DATABASE_DRIVER_MYSQL),
		})
		StoreTypes = append(StoreTypes, &storeType{
			Name:        "PostgreSQL",
			SqlSettings: storetest.MakeSqlSettings(model.DATABASE_DRIVER_POSTGRES),
		})
	}

	defer func() {
		if err := recover(); err != nil {
			tearDownStores()
			panic(err)
		}
	}()

	var wg sync.WaitGroup
	for _, st := range StoreTypes {
		st := st
		wg.Add(1)
		go func() {
			defer wg.Done()
			st.SqlStore = New(*st.SqlSettings, nil)
			st.Store = st.SqlStore
			st.Store.DropAllTables()
			st.Store.MarkSystemRanUnitTests()
		}()
	}
	wg.Wait()
}

func tearDownStores() {
	if testing.Short() {
		return
	}
	tearDownStoresOnce.Do(func() {
		var wg sync.WaitGroup
		wg.Add(len(StoreTypes))
		for _, st := range StoreTypes {
			st := st
			go func() {
				if st.Store != nil {
					st.Store.Close()
				}
				if st.SqlSettings != nil {
					storetest.CleanupSqlSettings(st.SqlSettings)
				}
				wg.Done()
			}()
		}
		wg.Wait()
	})
}

func StoreTestSuiteWithSqlSupplier(t *testing.T, testSuite StoreTestBaseSuite, executeFunc func(t *testing.T, testSuite StoreTestBaseSuite)) {
	for _, st := range StoreTypes {
		st := st
		t.Run(st.Name, func(t *testing.T) {
			if testing.Short() {
				t.SkipNow()
			}
			testSuite.SetStore(st.Store)
			testSuite.SetSqlStore(st.SqlStore)
			// suite.Run(t, testSuite)
			executeFunc(t, testSuite)
		})
	}
}

func StoreTestMysqlTestSuite(t *testing.T, testSuite *suite.Suite) {
	// Setup a global logger to catch tests logging outside of app context
	// The global logger will be stomped by apps initializing but that's fine for testing. Ideally this won't happen.
	mlog.InitGlobalLogger(mlog.NewLogger(&mlog.LoggerConfiguration{
		EnableConsole: true,
		ConsoleJson:   true,
		ConsoleLevel:  "error",
		EnableFile:    false,
	}))

	// dir, err := ioutil.TempDir("", "")
	// require.NoError(t, err)
	// defer os.RemoveAll(dir)

	suite.Run(t, testSuite)

	// suite.Run(t, &FileBackendTestSuite{
	// 	settings: model.FileSettings{
	// 		DriverName: model.NewString(model.IMAGE_DRIVER_LOCAL),
	// 		Directory:  &dir,
	// 	},
	// })
}
