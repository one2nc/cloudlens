package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/aws"
)

type EcsClusters struct {
}

func (ecs EcsClusters) Header() Header {
	return Header{
		HeaderColumn{Name: "Name", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Status", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "RunningTasksCount", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Arn", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
	}
}

func (ecs EcsClusters) Render(o interface{}, ns string, row *Row) error {
	ecsResp, ok := o.(aws.EcsClusterResp)
	if !ok {
		return fmt.Errorf("expected EcsResp, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		ecsResp.ClusterName,
		ecsResp.Status,
		ecsResp.RunningTasksCount,
		ecsResp.ClusterArn,
	}
	return nil
}
