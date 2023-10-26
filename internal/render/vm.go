package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/gcp"
)

type VM struct {
}

// Header returns a header row.
func (vm VM) Header() Header {
	return Header{
		HeaderColumn{Name: "Instance-Id", SortIndicatorIdx: 9, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Instance-State", SortIndicatorIdx: 9, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		// HeaderColumn{Name: "Instance-Type", SortIndicatorIdx: 9, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		// HeaderColumn{Name: "Monitoring-State", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Launch-Time", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: true},
		// HeaderColumn{Name: "Public-DNS", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: true, MX: false, Time: false},
		HeaderColumn{Name: "Availability-Zone", SortIndicatorIdx: -1, Align: tview.AlignCenter, Hide: false, Wide: false, MX: false, Time: false},
	}
}

func (vm VM) Render(o interface{}, ns string, row *Row) error {
	vmResp, ok := o.(gcp.VMResp)

	if !ok {
		return fmt.Errorf("Expected EC2Resp, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		vmResp.InstanceId,
		vmResp.InstanceState,
		// vmResp.InstanceType,
		// ec2Resp.MonitoringState,
		vmResp.LaunchTime,
		// ec2Resp.PublicDNS,
		vmResp.AvailabilityZone,
	}

	return nil
}
