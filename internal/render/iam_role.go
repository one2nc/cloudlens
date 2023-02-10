package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloud-lens/internal/aws"
)

type IamRole struct {
}

func (ir IamRole) Header() Header {
	return Header{
		HeaderColumn{Name: "User-Id", SortIndicatorIdx: 5, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "User-Name", SortIndicatorIdx: 5, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "ARN", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Created-Date", SortIndicatorIdx: 8, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: true},
	}
}

func (ir IamRole) Render(o interface{}, ns string, row *Row) error {
	irResp, ok := o.(aws.IamRoleResp)
	if !ok {
		return fmt.Errorf("expected iam role didn't receive, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		irResp.RoleId,
		irResp.RoleName,
		irResp.ARN,
		irResp.CreationTime,
	}
	return nil
}
