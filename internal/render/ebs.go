package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloud-lens/internal/aws"
)

type EBS struct {
}

// Header returns a header row.
func (ebs EBS) Header() Header {
	return Header{
		HeaderColumn{Name: "Volume-Id", SortIndicatorIdx: 7, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Size", SortIndicatorIdx: 0, Align: tview.AlignRight, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Volume-Type", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "State", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Availability-Zone", SortIndicatorIdx: 13, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Snapshot", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: true},
		HeaderColumn{Name: "Creation-Time", SortIndicatorIdx: 9, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
	}
}

func (ebs EBS) Render(o interface{}, ns string, row *Row) error {
	ebsResp, ok := o.(aws.EBSResp)

	if !ok {
		return fmt.Errorf("Expected EBS response, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		ebsResp.VolumeId,
		ebsResp.Size,
		ebsResp.VolumeType,
		ebsResp.State,
		ebsResp.AvailabilityZone,
		ebsResp.Snapshot,
		ebsResp.CreationTime,
	}

	return nil
}
