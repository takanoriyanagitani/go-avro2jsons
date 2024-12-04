package conv

import (
	"context"
	"io"
	"iter"
	"os"

	ho "github.com/hamba/avro/v2/ocf"

	util "github.com/takanoriyanagitani/go-avro2jsons/util"
)

func ReaderToMaps(r io.Reader) util.IO[iter.Seq2[map[string]any, error]] {
	return func(_ context.Context) (iter.Seq2[map[string]any, error], error) {
		dec, e := ho.NewDecoder(r)
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

var StdinToMaps util.IO[iter.Seq2[map[string]any, error]] = ReaderToMaps(
	os.Stdin,
)
