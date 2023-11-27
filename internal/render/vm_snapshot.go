package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/gcp"
)

type VMS struct {
}

func (vms VMS) Header() Header {
	return Header{
		HeaderColumn{Name: "Snapshot-Id", SortIndicatorIdx: 9, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Volume-Size", SortIndicatorIdx: 7, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Start-Time", SortIndicatorIdx: 6, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: true},
	}
}

func (vms VMS) Render(o interface{}, ns string, row *Row) error {
	sResp, ok := o.(gcp.SnapshotResp)
	if !ok {
		return fmt.Errorf("Expected SnapshotResp, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		sResp.Name,

		sResp.Size,
		sResp.CreatedAt,
	}
	return nil
}
