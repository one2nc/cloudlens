package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/gcp"
)

type Disk struct {
}

// Header returns a header row.
func (disk Disk) Header() Header {
	return Header{
		HeaderColumn{Name: "Volume-Id", SortIndicatorIdx: 7, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Size", SortIndicatorIdx: 0, Align: tview.AlignRight, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Volume-Type", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Status", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Availability-Zone", SortIndicatorIdx: 13, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Creation-Time", SortIndicatorIdx: 9, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
	}
}

func (disk Disk) Render(o interface{}, ns string, row *Row) error {
	diskResp, ok := o.(gcp.DiskResp)

	if !ok {
		return fmt.Errorf("Expected DiskResp response, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		diskResp.Name,
		diskResp.Size,
		diskResp.Type,
		diskResp.Status,
		diskResp.Zone,
		// diskResp.Snapshot,
		diskResp.CreationTime,
	}

	return nil
}
