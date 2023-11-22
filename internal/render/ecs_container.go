package render

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/ecs/types"
	"github.com/derailed/tview"
)

type EcsContainers struct {
}

func (ecs EcsContainers) Header() Header {
	return Header{
		HeaderColumn{Name: "ContainerName", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "RuntimeId", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Image URI", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: true, MX: false, Time: false},
		HeaderColumn{Name: "ImageDigest", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: true, MX: false, Time: false},
		HeaderColumn{Name: "LastStatus", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "HealthStatus", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Cpu", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "MemoryHardLimit", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
	}

}

func (ecs EcsContainers) Render(o interface{}, ns string, row *Row) error {
	container, ok := o.(types.Container)
	if !ok {
		return fmt.Errorf("expected EcsServiceResp, but got %T", o)
	}
	row.ID = ns
	row.Fields = Fields{
		*container.Name,
		*container.RuntimeId,
		*container.Image,
		*container.ImageDigest,
		*container.LastStatus,
		fmt.Sprintf("%v", container.HealthStatus),
		*container.Cpu,
		*container.Memory,
	}
	return nil
}
