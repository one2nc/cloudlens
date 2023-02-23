package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/aws"
)

type BObj struct {
}

func (obj BObj) Header() Header {
	return Header{
		HeaderColumn{Name: "Name", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Type", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Last-Modified", SortIndicatorIdx: 5, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: true},
		HeaderColumn{Name: "Size", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Storage-Class", SortIndicatorIdx: 8, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: true},
	}
}

func (obj BObj) Render(o interface{}, ns string, row *Row) error {
	s3Resp, ok := o.(aws.S3Object)
	if !ok {
		return fmt.Errorf("expected S3Resp, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		s3Resp.Name,
		s3Resp.ObjectType,
		s3Resp.LastModified,
		s3Resp.Size,
		s3Resp.StorageClass,
	}
	return nil
}
