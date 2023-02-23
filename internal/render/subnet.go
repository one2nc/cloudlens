package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/aws"
)

type Subnet struct {
}

func (sn Subnet) Header() Header {
	return Header{
		HeaderColumn{Name: "Subnet-Id", SortIndicatorIdx: 7, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Owner-Id", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Cidr Block", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Availability Zone", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Subnet-State", SortIndicatorIdx: 7, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
	}
}

func (sn Subnet) Render(o interface{}, ns string, row *Row) error {
	snResp, ok := o.(aws.SubnetResp)
	if !ok {
		return fmt.Errorf("expected vpc, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		snResp.SubnetId,
		snResp.OwnerId,
		snResp.CidrBlock,
		snResp.AvailabilityZone,
		snResp.State,
	}
	return nil
}
