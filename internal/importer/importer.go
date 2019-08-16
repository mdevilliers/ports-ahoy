package importer

import (
	"encoding/json"
	"io"

	"github.com/mdevilliers/ports-ahoy/internal/store"
	"github.com/pkg/errors"
)

type importService struct {
	decoder *json.Decoder
	store   storer
}

var (
	errEOF         = errors.New("end of file")
	ErrInvalidJSON = errors.New("invalid JSON")
)

type storer interface {
	Save(store.Port) error
}

// New returns an importer instance or an error
func New(in io.Reader, store storer) (*importService, error) {

	decoder := json.NewDecoder(in)

	// burn the first '{'
	_, err := decoder.Token()

	if err != nil {
		return nil, err
	}

	return &importService{
		decoder: decoder,
		store:   store,
	}, nil
}

// Import reads one record at a time saving the record or returns the first error
func (i *importService) Import() (err error) {

	// json.Decoder tends to throw panics on invalid json
	// TODO : investigate some more...
	defer func() {
		if r := recover(); r != nil {
			err = ErrInvalidJSON
		}
	}()

	for {
		err = i.next()

		if err == errEOF {
			break
		}

		if err != nil {
			return err
		}
	}

	return nil
}

func (i *importService) next() error {

	// read record by record
	if !i.decoder.More() {
		return errEOF
	}

	n, err := i.decoder.Token()

	if err != nil {
		return errors.Wrap(err, "no name token")
	}

	key, ok := n.(string)

	if !ok {
		return errors.New("expected name to be a string")
	}

	p := store.Port{
		Key: key,
	}

	err = i.decoder.Decode(&p)

	if err != nil {
		return errors.Wrap(err, "error decoding details")
	}

	err = i.store.Save(p)

	if err != nil {
		return err
	}

	return nil

}
