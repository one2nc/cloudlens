package render

import (
	"fmt"

	"github.com/derailed/tview"
	"github.com/one2nc/cloudlens/internal/aws"
)

type S3 struct {
}

func (s3 S3) Header() Header {
	return Header{
		HeaderColumn{Name: "Bucket-Name", SortIndicatorIdx: 0, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: false},
		HeaderColumn{Name: "Creation-Time", SortIndicatorIdx: 9, Align: tview.AlignLeft, Hide: false, Wide: false, MX: false, Time: true},
	}
}

func (s3 S3) Render(o interface{}, ns string, row *Row) error {
	s3Resp, ok := o.(aws.BucketResp)
	if !ok {
		return fmt.Errorf("Expected S3Resp, but got %T", o)
	}

	row.ID = ns
	row.Fields = Fields{
		s3Resp.BucketName,
		s3Resp.CreationTime,
	}
	return nil
}
