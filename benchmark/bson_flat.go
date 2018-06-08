package benchmark

import (
	"context"
	"io/ioutil"
	"path/filepath"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/pkg/errors"
)

const (
	perfDataDir  = "perf"
	bsonDataDir  = "extended_bson"
	flatBSONData = "flat_bson.json"
)

func FlatBSONEncodingDocument(ctx context.Context, iters int) error {
	data, err := ioutil.ReadFile(filepath.Join(getProjectRoot(), perfDataDir, bsonDataDir, flatBSONData))
	if err != nil {
		return err
	}
	doc, err := bson.ParseExtJSONObject(string(data))
	if err != nil {
		return err
	}
	if doc.Len() != 145 {
		return errors.New("bson parsing error")
	}

	for i := 0; i < iters; i++ {
		out, err := doc.MarshalBSON()
		if err != nil {
			return err
		}
		if len(out) == 0 {
			return errors.New("marshaling error")
		}
	}
	return nil
}
