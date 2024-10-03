package symfile

import (
	"io"
	"sudonters/zootler/icearrow/zasm"
)

type WriterTo zasm.Assembly

func (wt WriterTo) WriteTo(w io.WriterAt) (int, error) {
	return 0, nil
}
