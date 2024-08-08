package fs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syfar-ls/document"

	"github.com/tliron/commonlog"
)

type FileStorage struct {
	workingDir string
	docStore   *DocumentStore
}

func NewFileStorage() *FileStorage {

	return &FileStorage{docStore: &DocumentStore{documents: make(map[string]*document.Document)}}
}

func NewDocument(uri string, content string) (*document.Document, error) {

	path, err := normalizePath(uri)
	if err != nil {
		return nil, err
	}
	doc := &document.Document{
		URI:     uri,
		Path:    path,
		Content: content,
	}
	return doc, nil
}

func (fs *FileStorage) GetDocument(pathOrURI string) (*document.Document, error) {
	commonlog.NewInfoMessage(0, fmt.Sprintf("Geting %s...", pathOrURI))
	return fs.docStore.GetDocument(pathOrURI)
}

func (fs *FileStorage) AddDocument(pathOrURI string, doc *document.Document) error {
	commonlog.NewInfoMessage(0, fmt.Sprintf("Adding %s...", pathOrURI))

	return fs.docStore.AddDocument(pathOrURI, doc)
}

func (fs *FileStorage) ColseDocument(pathOrURI string) error {
	commonlog.NewInfoMessage(0, fmt.Sprintf("Closing %s...", pathOrURI))

	return fs.docStore.ColseDocument(pathOrURI)
}

func (fs *FileStorage) WorkingDir() string {
	return fs.workingDir
}

func (fs *FileStorage) SetWorkingDir(path string) {
	fs.workingDir = path
}

func (fs *FileStorage) Abs(path string) (string, error) {
	var err error
	if !filepath.IsAbs(path) {
		path = filepath.Join(fs.workingDir, path)
		path, err = filepath.Abs(path)
		if err != nil {
			return path, err
		}
	}

	return path, nil
}

func (fs *FileStorage) Rel(path string) (string, error) {
	return filepath.Rel(fs.workingDir, path)
}

func (fs *FileStorage) Canonical(path string) (string, error) {
	path = filepath.Clean(path)

	resolvedPath, err := filepath.EvalSymlinks(path)
	if err != nil {
		return path, err
	} else {
		path = resolvedPath
	}

	return path, nil
}

func (fs *FileStorage) FileExists(path string) (bool, error) {
	fi, err := fs.fileInfo(path)
	if err != nil {
		return false, err
	} else {
		return fi != nil && (*fi).Mode().IsRegular(), nil
	}
}

func (fs *FileStorage) DirExists(path string) (bool, error) {
	fi, err := fs.fileInfo(path)
	return !os.IsNotExist(err) && fi != nil && (*fi).Mode().IsDir(), nil
}

func (fs *FileStorage) fileInfo(path string) (*os.FileInfo, error) {
	if fi, err := os.Stat(path); err == nil {
		return &fi, nil
	} else if os.IsNotExist(err) {
		return nil, nil
	} else {
		return nil, err
	}
}

func (fs *FileStorage) IsDescendantOf(dir string, path string) (bool, error) {
	dir, err := fs.Abs(dir)
	if err != nil {
		return false, err
	}
	dir, err = fs.Canonical(dir)
	if err != nil {
		return false, err
	}
	path, err = fs.Abs(path)
	if err != nil {
		return false, err
	}
	path, err = fs.Canonical(path)
	if err != nil {
		return false, err
	}
	path, err = filepath.Rel(dir, path)
	if err != nil {
		return false, err
	}

	return !strings.HasPrefix(path, ".."), nil
}

func (fs *FileStorage) Read(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (fs *FileStorage) Write(path string, content []byte) error {
	dir := filepath.Dir(path)
	if dir != "." && dir != ".." {
		err := os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()
	_, err = f.Write(content)
	return err
}
