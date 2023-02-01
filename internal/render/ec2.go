package render

import (
	"fmt"

	"github.com/one2nc/cloud-lens/internal/aws"
)

type EC2 struct {
}

// Header returns a header row.
func (e EC2) Header() Header {
	return Header{
		HeaderColumn{Name: "Instance-Id"},
		HeaderColumn{Name: "Instance-State"},
		HeaderColumn{Name: "Instance-Type"},
		HeaderColumn{Name: "Monitoring-State"},
		HeaderColumn{Name: "Launch-Time"},
		HeaderColumn{Name: "Public-DNS", Wide: true},
		HeaderColumn{Name: "Availability-Zone"},
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
