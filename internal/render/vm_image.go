package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/gcp"
)

type VMI struct {
}

func (vmi VMI) Header() Header {
	return Header{
		HeaderColumn{Name: "Image-Id", SortIndicatorIdx: 6, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Image-Location", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Status", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Creation-Time", SortIndicatorIdx: 9, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: true},
	}
}

func (vmi VMI) Render(o interface{}, ns string, row *Row) error {
	iResp, ok := o.(gcp.ImageResp)
	if !ok {
		return fmt.Errorf("Expected ImageResp, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		iResp.Name,
		iResp.Location,
		iResp.Status,
		iResp.CreatedAt,
	}
	return nil
}
