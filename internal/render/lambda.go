package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloud-lens/internal/aws"
)

type Lambda struct {
}

// Header returns a header row.
func (l Lambda) Header() Header {
	return Header{
		HeaderColumn{Name: "Function-Name", SortIndicatorIdx: 9, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Description", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Role", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Function-Arn", SortIndicatorIdx: 9, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Code-Size", SortIndicatorIdx: 5, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Last-Modified", SortIndicatorIdx: 5, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: true},
	}
}

func (l Lambda) Render(o interface{}, ns string, row *Row) error {
	lambdaResp, ok := o.(aws.LambdaResp)
	if !ok {
		return fmt.Errorf("Expected sqsResp, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		lambdaResp.FunctionName,
		lambdaResp.Description,
		lambdaResp.Role,
		lambdaResp.FunctionArn,
		lambdaResp.CodeSize,
		lambdaResp.LastModified,
	}

	return nil
}
