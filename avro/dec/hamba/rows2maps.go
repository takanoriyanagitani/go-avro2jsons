package conv

import (
	"context"
	"io"
	"iter"
	"os"

	ha "github.com/hamba/avro/v2"
	ho "github.com/hamba/avro/v2/ocf"

	aj "github.com/takanoriyanagitani/go-avro2jsons"
	util "github.com/takanoriyanagitani/go-avro2jsons/util"
)

func ConfigToOptions(cfg aj.AvroConfig) []ho.DecoderFunc {
	var hcfg ha.Config = ha.Config{}
	hcfg.MaxByteSliceSize = int(cfg.BlobSizeMax())
	var hapi ha.API = hcfg.Freeze()
	return []ho.DecoderFunc{
		ho.WithDecoderConfig(hapi),
	}
}

func ReaderToMapsWithHambaOptions(
	r io.Reader,
	opts ...ho.DecoderFunc,
) util.IO[iter.Seq2[map[string]any, error]] {
	return func(_ context.Context) (iter.Seq2[map[string]any, error], error) {
		dec, e := ho.NewDecoder(r, opts...)
		if nil != e {
			return nil, e
		}
		return func(yield func(map[string]any, error) bool) {
			var buf map[string]any = map[string]any{}
			var err error = nil

			for dec.HasNext() {
				clear(buf)
				err = dec.Decode(&buf)

				if !yield(buf, err) {
					return
				}
			}
		}, nil
	}
}

func ReaderToMapsWithConfig(
	r io.Reader,
	cfg aj.AvroConfig,
) util.IO[iter.Seq2[map[string]any, error]] {
	var opts []ho.DecoderFunc = ConfigToOptions(cfg)
	return ReaderToMapsWithHambaOptions(r, opts...)
}

func StdinToMapsWithConfig(
	cfg aj.AvroConfig,
) util.IO[iter.Seq2[map[string]any, error]] {
	return ReaderToMapsWithConfig(os.Stdin, cfg)
}

func ReaderToMaps(r io.Reader) util.IO[iter.Seq2[map[string]any, error]] {
	return ReaderToMapsWithHambaOptions(r)
}

var StdinToMaps util.IO[iter.Seq2[map[string]any, error]] = ReaderToMaps(
	os.Stdin,
)
