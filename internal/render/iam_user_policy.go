package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/aws"
)

type IamUserPloicy struct {
}

func (iup IamUserPloicy) Header() Header {
	return Header{
		HeaderColumn{Name: "Policy-ARN", SortIndicatorIdx: 7, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Policy-Name", SortIndicatorIdx: 7, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
	}
}

func (iup IamUserPloicy) Render(o interface{}, ns string, row *Row) error {
	usrPolicy, ok := o.(aws.IAMUSerPolicyResponse)
	if !ok {
		return fmt.Errorf("expected usrPolicy, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		usrPolicy.PolicyArn,
		usrPolicy.PolicyName,
	}
	return nil
}
