package config

import (
	"path/filepath"
	"time"
  "io/ioutil"
  "os"
  "testing"

  "github.com/stretchr/testify/require"
)

func TestWatcherInvalidDirectory(t *testing.T) {
  if testing.Short() {
    t.Skip("skipping watcher test in short mode")
  }

  callback := func() {}
  _, err := newWatcher("/does/not/exist", callback)
  require.Error(t, err, "should have failed to watch a non-existent directory")
}

func TestWatcher(t *testing.T) {
  if testing.Short() {
    t.Skip("skipping watcher test in short mode")
  }

  tempDir, err := ioutil.TempDir("", "TestWatcher")
  require.NoError(t, err)
  defer os.RemoveAll(tempDir)

  f, err := ioutil.TempFile(tempDir, "TestWatcher")
  require.NoError(t, err)
  defer f.Close()
  defer os.Remove(f.Name())

  called := make(chan bool)
  callback := func() {
    called <- true
  }
  watcher, err := newWatcher(f.Name(), callback)
  require.NoError(t, err)
  defer watcher.Close()

  // Write to a different file
  ioutil.WriteFile(filepath.Join(tempDir, "unrelated"), []byte("data"), 0644)
  select {
  case <-called:
    require.Fail(t, "callback should not have been called for unrelated file")
  case <-time.After(1 * time.Second):
  }

  ioutil.WriteFile(f.Name(), []byte("data"), 0644)
  select {
  case <-called:
  case <-time.After(1 * time.Second):
    require.Fail(t, "callback should have been called when file written")
  }
}
