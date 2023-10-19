package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/gcp"
)

type SOBJ struct {
}

func (obj SOBJ) Header() Header {
	return Header{
		HeaderColumn{Name: "Name", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Type", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Last-Modified", SortIndicatorIdx: 5, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: true},
		HeaderColumn{Name: "Size", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Storage-Class", SortIndicatorIdx: 8, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: true},
		HeaderColumn{Name: "SizeInBytes", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: true, Wide: false, MX: false, Time: false},
	}
}

func (obj SOBJ) Render(o interface{}, ns string, row *Row) error {
	storageObjResp, ok := o.(gcp.StorageObjResp)
	if !ok {
		return fmt.Errorf("expected StorageObjResp, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		storageObjResp.Name,
		storageObjResp.ObjectType,
		storageObjResp.LastModified,
		storageObjResp.Size,
		storageObjResp.StorageClass,
		fmt.Sprint(storageObjResp.SizeInBytes),
	}
	return nil
}
