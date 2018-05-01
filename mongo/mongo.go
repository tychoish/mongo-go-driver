// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package mongo

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"reflect"
	"strings"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// Dialer is used to make network connections.
type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

// DocumentTransformer provides a way for clients to modify the way
// that a client's collection methods convert input documents into
// bson documents in the form of byte slices ([]byte).
//
// These functions allow you to wrap other bson libraries or enforce
// different conversion strategies within the driver.
//
// These transformers are *not* called for well known types, including:
//
//  bson.Marshaler
//  bson.DocumentMarshaler
//  bson.Reader
//  []byte (must be a valid BSON document)
//  io.Reader (only 1 BSON document will be read)
//
type DocumentTransformer func(interface{}) ([]byte, error)

// TransformDocument is the default document transformer
// implementation which uses reflection to convert structs into bson
// documents (in the format of []byte objects.)
func TransformDocument(document interface{}) ([]byte, error) {
	var kind reflect.Kind
	if t := reflect.TypeOf(document); t.Kind() == reflect.Ptr {
		kind = t.Elem().Kind()
	}

	if reflect.ValueOf(document).Kind() == reflect.Struct || kind == reflect.Struct {
		return bson.NewDocumentEncoder().EncodeDocument(document).MarshalBSON()
	}

	if reflect.ValueOf(document).Kind() == reflect.Map &&
		reflect.TypeOf(document).Key().Kind() == reflect.String {
		return bson.NewDocumentEncoder().EncodeDocument(document).MarshalBSON()
	}

	return nil, fmt.Errorf("cannot transform type %s to a *bson.Document", reflect.TypeOf(document))
}

func convertToDocument(dt DocumentTransformer, document interface{}) (bson.Reader, error) {
	switch d := document.(type) {
	case nil:
		return bson.Reader([]byte{}), nil
	case []byte:
		return bson.Reader(d), nil
	case bson.Reader:
		return d, nil
	case io.Reader:
		return bson.NewFromIOReader(d)
	case *bson.Document, bson.Marshaler:
		return d.MarshalBSON()
	case bson.DocumentMarshaler:
		return d.MarshalBSONDocument().MarshalBSON()
	default:
		doc, err := dt(document)
		if err != nil {
			return nil, err
		}

		return bson.Reader(doc), nil
	}
}

func ensureID(d bson.Reader) (interface{}, error) {
	var id interface{}

	elem, err := d.Lookup("_id")
	switch {
	case err == bson.ErrElementNotFound:
		oid := objectid.New()
		d = append(d, bson.EC.ObjectID("_id", oid).MarshalBSON()...)
		id = oid
	case err != nil:
		return nil, err
	default:
		id = elem
	}
	return id, nil
}

func ensureDollarKey(doc *bson.Document) error {
	if elem, ok := doc.ElementAtOK(0); !ok || !strings.HasPrefix(elem.Key(), "$") {
		return errors.New("update document must contain key beginning with '$'")
	}
	return nil
}
