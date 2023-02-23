package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/aws"
)

type EC2S struct {
}

func (es EC2S) Header() Header {
	return Header{
		HeaderColumn{Name: "Snapshot-Id", SortIndicatorIdx: 9, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Owner-Id", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Volume-Id", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Volume-Size", SortIndicatorIdx: 7, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "State", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Start-Time", SortIndicatorIdx: 6, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: true},
	}
}

func (es EC2S) Render(o interface{}, ns string, row *Row) error {
	esResp, ok := o.(aws.Snapshot)
	if !ok {
		return fmt.Errorf("Expected esResp, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		esResp.SnapshotId,
		esResp.OwnerId,
		esResp.VolumeId,
		esResp.VolumeSize,
		esResp.State,
		esResp.StartTime,
	}
	return nil
}
