package conv

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"iter"
	"os"

	util "github.com/takanoriyanagitani/go-avro2jsons/util"
)

func MapsToWriter(
	w io.Writer,
) func(iter.Seq2[map[string]any, error]) util.IO[util.Void] {
	return func(m iter.Seq2[map[string]any, error]) util.IO[util.Void] {
		return func(ctx context.Context) (util.Void, error) {
			var bw *bufio.Writer = bufio.NewWriter(w)
			defer bw.Flush()

			var enc *json.Encoder = json.NewEncoder(bw)
			var err error = nil

			for row, e := range m {
				select {
				case <-ctx.Done():
					return util.Empty, ctx.Err()
				default:
				}

				if nil != e {
					return util.Empty, e
				}

				err = enc.Encode(row)
				if nil != err {
					return util.Empty, err
				}
			}
			return util.Empty, nil
		}
	}
}

var MapsToStdout func(
	iter.Seq2[map[string]any, error],
) util.IO[util.Void] = MapsToWriter(os.Stdout)
