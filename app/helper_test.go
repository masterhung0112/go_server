package app

import (
	"time"
	"sync"
	"bytes"
	"testing"

	"github.com/masterhung0112/go_server/config"
	"github.com/masterhung0112/go_server/model"
	"github.com/masterhung0112/go_server/store"
)

type TestHelper struct {
	App       *App
	Server    *Server
	LogBuffer *bytes.Buffer

	tempWorkspace string
}

func setupTestHelper(dbStore store.Store, tb testing.TB, configSet func(*model.Config)) *TestHelper {

	memoryStore, err := config.NewMemoryStoreWithOptions(&config.MemoryStoreOptions{IgnoreEnvironmentOverrides: true})
	if err != nil {
		panic("failed to initialize memory store: " + err.Error())
	}
	config := memoryStore.Get()
	if configSet != nil {
		configSet(config)
	}
	// *config.PluginSettings.Directory = filepath.Join(tempWorkspace, "plugins")
	// *config.PluginSettings.ClientDirectory = filepath.Join(tempWorkspace, "webapp")
	// *config.LogSettings.EnableSentry = false // disable error reporting during tests
	memoryStore.Set(config)

	buffer := &bytes.Buffer{}

	var options []Option
	options = append(options, ConfigStore(memoryStore))
	options = append(options, StoreOverride(dbStore))
	// options = append(options, SetLogger(mlog.NewTestingLogger(tb, buffer)))

	s, err := NewServer(options...)
	if err != nil {
		panic(err)
	}

	th := &TestHelper{
		App:       New(ServerConnector(s)),
		Server:    s,
		LogBuffer: buffer,
	}

	// th.App.UpdateConfig(func(cfg *model.Config) { *cfg.TeamSettings.MaxUsersPerTeam = 50 })
	// th.App.UpdateConfig(func(cfg *model.Config) { *cfg.RateLimitSettings.Enable = false })
	// prevListenAddress := *th.App.Config().ServiceSettings.ListenAddress
	// th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.ListenAddress = ":0" })
	// Start HTTP Server and other stuff
	if err := th.Server.Start(); err != nil {
		panic(err)
	}

	// th.App.UpdateConfig(func(cfg *model.Config) { *cfg.ServiceSettings.ListenAddress = prevListenAddress })

	// th.App.Srv().SearchEngine = mainHelper.SearchEngine

	// th.App.Srv().Store.MarkSystemRanUnitTests()

	// th.App.UpdateConfig(func(cfg *model.Config) { *cfg.TeamSettings.EnableOpenServer = true })

	// Disable strict password requirements for test
	th.App.UpdateConfig(func(cfg *model.Config) {
		*cfg.PasswordSettings.MinimumLength = 5
		*cfg.PasswordSettings.Lowercase = false
		*cfg.PasswordSettings.Uppercase = false
		*cfg.PasswordSettings.Symbol = false
		*cfg.PasswordSettings.Number = false
	})

	th.App.InitServer()

	return th
}

func Setup(tb testing.TB) *TestHelper {
	if testing.Short() {
		tb.SkipNow()
	}
	dbStore := mainHelper.GetStore()
	dbStore.DropAllTables()
	dbStore.MarkSystemRanUnitTests()

	return setupTestHelper(dbStore, tb, nil)
}

var initBasicOnce sync.Once
func (me *TestHelper) InitBasic() *TestHelper {
  // create users once and cache them because password hashing is slow
  initBasicOnce.Do(func() {
  })
  return me
}

func (me *TestHelper) ShutdownApp() {
	done := make(chan bool)
	go func() {
		me.Server.Shutdown()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(30 * time.Second):
		// panic instead of fatal to terminate all tests in this package, otherwise the
		// still running App could spuriously fail subsequent tests.
		panic("failed to shutdown App within 30 seconds")
	}
}

func (me *TestHelper) TearDown() {
	// if me.IncludeCacheLayer {
	// 	// Clean all the caches
	// 	me.App.Srv().InvalidateAllCaches()
	// }
	me.ShutdownApp()
	// if me.tempWorkspace != "" {
	// 	os.RemoveAll(me.tempWorkspace)
	// }
}