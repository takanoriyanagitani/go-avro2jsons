package main

import (
	"context"
	"fmt"
	"iter"
	"log"
	"os"
	"strconv"

	aj "github.com/takanoriyanagitani/go-avro2jsons"
	util "github.com/takanoriyanagitani/go-avro2jsons/util"

	dh "github.com/takanoriyanagitani/go-avro2jsons/avro/dec/hamba"
	es "github.com/takanoriyanagitani/go-avro2jsons/json/enc/std"
)

func EnvByKey(key string) util.IO[string] {
	return func(_ context.Context) (string, error) {
		val, found := os.LookupEnv(key)
		switch found {
		case true:
			return val, nil
		default:
			return "", fmt.Errorf("env var %s missing", key)
		}
	}
}

var blobSizeMaxString util.IO[string] = EnvByKey("ENV_BLOB_SIZE_MAX").
	OrElse(util.Of("1048576"))

var StringToBlobSizeMax func(string) (aj.AvroBlobSizeMax, error) = util.
	ComposeErr(
		strconv.Atoi,
		func(i int) (aj.AvroBlobSizeMax, error) {
			return aj.AvroBlobSizeMax(i), nil
		},
	)

var blobSizeMax util.IO[aj.AvroBlobSizeMax] = util.Bind(
	blobSizeMaxString,
	util.Lift(StringToBlobSizeMax),
)

var config util.IO[aj.AvroConfig] = util.Bind(
	blobSizeMax,
	util.Lift(func(b aj.AvroBlobSizeMax) (aj.AvroConfig, error) {
		return aj.AvroConfigDefault.
			WithBlobSizeMax(b), nil
	}),
)

var stdin2maps util.IO[iter.Seq2[map[string]any, error]] = util.Bind(
	config,
	func(c aj.AvroConfig) util.IO[iter.Seq2[map[string]any, error]] {
		return dh.StdinToMapsWithConfig(c)
	},
)

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
