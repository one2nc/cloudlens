package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/aws"
)

type IamUserGroupPloicy struct {
}

func (iugp IamUserGroupPloicy) Header() Header {
	return Header{
		HeaderColumn{Name: "Policy-ARN", SortIndicatorIdx: 7, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Policy-Name", SortIndicatorIdx: 7, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
	}
}

func (iugp IamUserGroupPloicy) Render(o interface{}, ns string, row *Row) error {
	usrGroupPolicy, ok := o.(aws.IAMUSerGroupPolicyResponse)
	if !ok {
		return fmt.Errorf("expected usrGroupPolicy, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		usrGroupPolicy.PolicyArn,
		usrGroupPolicy.PolicyName,
	}
	return nil
}
