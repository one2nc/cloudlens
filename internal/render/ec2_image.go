package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/aws"
)

type EC2I struct {
}

func (ei EC2I) Header() Header {
	return Header{
		HeaderColumn{Name: "Image-Id", SortIndicatorIdx: 6, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Owner-Id", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Image-Location", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Name", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Image-Type", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
	}
}

func (ei EC2I) Render(o interface{}, ns string, row *Row) error {
	eiResp, ok := o.(aws.ImageResp)
	if !ok {
		return fmt.Errorf("Expected eiResp, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		eiResp.ImageId,
		eiResp.OwnerId,
		eiResp.ImageLocation,
		eiResp.Name,
		eiResp.ImageType,
	}
	return nil
}
