package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/aws"
)

type VPC struct {
}

func (v VPC) Header() Header {
	return Header{
		HeaderColumn{Name: "VPC-Id", SortIndicatorIdx: 4, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Owner-Id", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Cidr Block", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Instance Tenancy", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "VPC-State", SortIndicatorIdx: 4, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
	}
}

func (v VPC) Render(o interface{}, ns string, row *Row) error {
	vpcResp, ok := o.(aws.VpcResp)
	if !ok {
		return fmt.Errorf("expected vpc, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		vpcResp.VpcId,
		vpcResp.OwnerId,
		vpcResp.CidrBlock,
		vpcResp.InstanceTenancy,
		vpcResp.State,
	}
	return nil
}
