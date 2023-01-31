package model

import (
	"github.com/one2nc/cloud-lens/internal/dao"
	"github.com/one2nc/cloud-lens/internal/render"
)

var Registry = map[string]ResourceMeta{
	"ec2": {
		DAO:      &dao.EC2{},
		Renderer: &render.EC2{},
	},
}
