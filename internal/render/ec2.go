package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/aws"
)

type EC2 struct {
}

// Header returns a header row.
func (e EC2) Header() Header {
	return Header{
		HeaderColumn{Name: "Instance-Id", SortIndicatorIdx: 9, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Instance-State", SortIndicatorIdx: 9, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Instance-Type", SortIndicatorIdx: 9, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Monitoring-State", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Launch-Time", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: true},
		HeaderColumn{Name: "Public-DNS", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: true, MX: false, Time: false},
		HeaderColumn{Name: "Availability-Zone", SortIndicatorIdx: -1, Align: tview.AlignCenter, Hide: false, Wide: false, MX: false, Time: false},
	}
}

func (e EC2) Render(o interface{}, ns string, row *Row) error {
	ec2Resp, ok := o.(aws.EC2Resp)

	if !ok {
		return fmt.Errorf("Expected EC2Resp, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		ec2Resp.InstanceId,
		ec2Resp.InstanceState,
		ec2Resp.InstanceType,
		ec2Resp.MonitoringState,
		ec2Resp.LaunchTime,
		ec2Resp.PublicDNS,
		ec2Resp.AvailabilityZone,
	}

	return nil
}
