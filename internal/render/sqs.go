package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloud-lens/internal/aws"
)

type SQS struct {
}

// Header returns a header row.
func (sqs SQS) Header() Header {
	return Header{
		HeaderColumn{Name: "URL", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Name", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Type", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Created", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: true},
		HeaderColumn{Name: "Messages-Available", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Encryption", SortIndicatorIdx: -1, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Max-Message-Size", SortIndicatorIdx: -1, Align: tview.AlignCenter, Hide: false, Wide: true, MX: false, Time: false},
	}
}

func (sqs SQS) Render(o interface{}, ns string, row *Row) error {
	sqsResp, ok := o.(aws.SQSResp)
	if !ok {
		return fmt.Errorf("Expected sqsResp, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		sqsResp.URL,
		sqsResp.Name,
		sqsResp.Type,
		sqsResp.Created,
		sqsResp.MessagesAvailable,
		sqsResp.Encryption,
		sqsResp.MaxMessageSize,
	}

	return nil
}
