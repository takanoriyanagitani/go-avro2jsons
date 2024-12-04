package main

import (
	"context"
	"iter"
	"log"

	util "github.com/takanoriyanagitani/go-avro2jsons/util"

	dh "github.com/takanoriyanagitani/go-avro2jsons/avro/dec/hamba"
	es "github.com/takanoriyanagitani/go-avro2jsons/json/enc/std"
)

var stdin2maps util.IO[iter.Seq2[map[string]any, error]] = dh.StdinToMaps

var stdin2maps2sink util.IO[util.Void] = util.Bind(
	stdin2maps,
	es.MapsToStdout,
)

func sub(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	_, e := stdin2maps2sink(ctx)
	return e
}

func main() {
	e := sub(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
