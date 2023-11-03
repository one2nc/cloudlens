package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/aws"
)

type EcsServices struct {
}

func (ecs EcsServices) Header() Header {
	return Header{
		HeaderColumn{Name: "Name", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Status", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		//HeaderColumn{Name: "Type", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "TaskDefinition", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Arn", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
	}
}

func (ecs EcsServices) Render(o interface{}, ns string, row *Row) error {
	ecsServiceResp, ok := o.(aws.EcsServiceResp)
	if !ok {
		return fmt.Errorf("expected EcsServiceResp, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		ecsServiceResp.ServiceName,
		ecsServiceResp.Status,
		//ecsServiceResp.,
		ecsServiceResp.TaskDefinition,
		ecsServiceResp.ServiceArn,
	}
	return nil
}
