package fs

import (
	"fmt"
	"syfar-ls/document"
)

type DocumentStore struct {
	documents map[string]*document.Document
}

func (s *DocumentStore) GetDocument(pathOrURI string) (*document.Document, error) {
	path, err := normalizePath(pathOrURI)
	if err != nil {

		return nil, err
	}
	d, ok := s.documents[path]
	if ok {
		return d, nil
	}
	return nil, fmt.Errorf("error when geting: %s", pathOrURI)
}

func (s *DocumentStore) AddDocument(pathOrURI string, doc *document.Document) error {
	path, err := normalizePath(pathOrURI)
	if err != nil {

		return err
	}
	s.documents[path] = doc
	return nil
}

func (s *DocumentStore) ColseDocument(pathOrURI string) error {
	path, err := normalizePath(pathOrURI)
	if err != nil {

		return err
	}
	delete(s.documents, path)
	return nil
}
