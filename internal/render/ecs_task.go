package render

import (
	"fmt"
	"time"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/aws"
)

type EcsTasks struct {
}

func (ecs EcsTasks) Header() Header {
	return Header{
		HeaderColumn{Name: "Task", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "TaskArn", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: true, MX: false, Time: false},
		HeaderColumn{Name: "TaskDefinitionArn", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: true, MX: false, Time: false},
		HeaderColumn{Name: "ContainerInstanceArn", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: true, MX: false, Time: false},
		HeaderColumn{Name: "LastStatus", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "DesiredStatus", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "HealthStatus", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "StartedAt", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: true},
		HeaderColumn{Name: "LaunchType", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "PlatformVersion", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "vCPU", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Memory", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "GroupName", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "StartedBy", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
	}
}

func (ecs EcsTasks) Render(o interface{}, ns string, row *Row) error {
	task, ok := o.(aws.EcsTaskResp)
	if !ok {
		return fmt.Errorf("expected EcsServiceResp, but got %T", o)
	}

	var startedTime time.Time
	if task.StartedAt != nil {
		startedTime = *task.StartedAt
	}
	containerInstanceArn, taskDefinitionArn, lastStatus, desiredStatus, platformVersion, cpu, memory, group, startedBy := "", "", "", "", "", "", "", "", ""
	if task.ContainerInstanceArn != nil {
		containerInstanceArn = *task.ContainerInstanceArn
	}
	if task.TaskDefinitionArn != nil {
		taskDefinitionArn = *task.TaskDefinitionArn
	}
	if task.LastStatus != nil {
		lastStatus = *task.LastStatus
	}
	if task.DesiredStatus != nil {
		desiredStatus = *task.DesiredStatus
	}
	if task.PlatformVersion != nil {
		platformVersion = *task.PlatformVersion
	}
	if task.Cpu != nil {
		cpu = *task.Cpu
	}
	if task.Memory != nil {
		memory = *task.Memory
	}
	if task.Group != nil {
		group = *task.Group
	}
	if task.StartedBy != nil {
		startedBy = *task.StartedBy
	}
	row.ID = ns
	row.Fields = Fields{
		task.TaskId,
		*task.TaskArn,
		taskDefinitionArn,
		containerInstanceArn,
		lastStatus,
		desiredStatus,
		fmt.Sprintf("%v", task.HealthStatus),
		startedTime.String(),
		fmt.Sprintf("%v", task.LaunchType),
		platformVersion,
		cpu,
		memory,
		group,
		startedBy,
	}
	return nil
}
