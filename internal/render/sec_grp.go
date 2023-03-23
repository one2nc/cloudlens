package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/aws"
)

type SG struct {
}

func (sg SG) Header() Header {
	return Header{
		HeaderColumn{Name: "Group-Id", SortIndicatorIdx: 6, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Group-Name", SortIndicatorIdx: 6, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Desription", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Owner-Id", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "VPC-Id", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
	}
}

func (sg SG) Render(o interface{}, ns string, row *Row) error {
	sgResp, ok := o.(aws.SGResp)
	if !ok {
		return fmt.Errorf("Expected SGResp, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		sgResp.GroupId,
		sgResp.GroupName,
		sgResp.Description,
		sgResp.OwnerId,
		sgResp.VpcId,
	}

	return nil
}
