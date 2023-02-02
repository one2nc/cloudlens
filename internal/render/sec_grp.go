package render

import (
	"fmt"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/derailed/tview"
)

type SG struct {
}

func (sg SG) Header() Header {
	return Header{
		HeaderColumn{Name: "Group-Id", SortIndicatorIdx: 6, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Group-Name", SortIndicatorIdx: 6, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: true},
		HeaderColumn{Name: "Desription", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: true},
		HeaderColumn{Name: "Owner-Id", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: true},
		HeaderColumn{Name: "VPC-Id", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: true},
	}
}

func (sg SG) Render(o interface{}, ns string, row *Row) error {
	sgResp, ok := o.(*ec2.SecurityGroup)
	if !ok {
		return fmt.Errorf("Expected S3Resp, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		*sgResp.GroupId,
		*sgResp.GroupName,
		*sgResp.Description,
		*sgResp.OwnerId,
		*sgResp.VpcId,
	}
	return nil
}
