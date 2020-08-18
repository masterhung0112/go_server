package filesstore

import (
	"io/ioutil"
	"net/http"
	"github.com/masterhung0112/go_server/mlog"
	"os"
	"io"
	"path/filepath"
	"bytes"
	"github.com/masterhung0112/go_server/model"
)

const (
	TEST_FILE_PATH = "/testfile"
)

type LocalFileBackend struct {
	directory string
}

func (b *LocalFileBackend) TestConnection() *model.AppError {
  f := bytes.NewReader([]byte("testingwrite"))
  if _, err := writeFileLocally(f, filepath.Join(b.directory, TEST_FILE_PATH)); err != nil{
		return model.NewAppError("TestFileConnection", "api.file.test_connection.local.connection.app_error", nil, err.Error(), http.StatusInternalServerError)
  }
  os.Remove(filepath.Join(b.directory, TEST_FILE_PATH))
	mlog.Debug("Able to write files to local storage.")
  return nil
}

func (b *LocalFileBackend) WriteFile(fr io.Reader, path string) (int64, *model.AppError) {
	return writeFileLocally(fr, filepath.Join(b.directory, path))
}

func writeFileLocally(fr io.Reader, path string) (int64, *model.AppError) {
  if err := os.MkdirAll(filepath.Dir(path), 0750); err != nil {
    directory, _ := filepath.Abs(filepath.Dir(path))
		return 0, model.NewAppError("WriteFile", "api.file.write_file_locally.create_dir.app_error", nil, "directory="+directory+", err="+err.Error(), http.StatusInternalServerError)
  }
  fw, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return 0, model.NewAppError("WriteFile", "api.file.write_file_locally.writing.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	defer fw.Close()
	written, err := io.Copy(fw, fr)
	if err != nil {
		return written, model.NewAppError("WriteFile", "api.file.write_file_locally.writing.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return written, nil
}

func (b *LocalFileBackend) RemoveFile(path string) *model.AppError {
	if err := os.Remove(filepath.Join(b.directory, path)); err != nil {
		return model.NewAppError("RemoveFile", "utils.file.remove_file.local.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return nil
}

func (b *LocalFileBackend) ReadFile(path string) ([]byte, *model.AppError) {
	f, err := ioutil.ReadFile(filepath.Join(b.directory, path))
	if err != nil {
		return nil, model.NewAppError("ReadFile", "api.file.read_file.reading_local.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return f, nil
}

func (b *LocalFileBackend) Reader(path string) (ReadCloseSeeker, *model.AppError) {
	f, err := os.Open(filepath.Join(b.directory, path))
	if err != nil {
		return nil, model.NewAppError("Reader", "api.file.reader.reading_local.app_error", nil, err.Error(), http.StatusInternalServerError)
	}
	return f, nil
}