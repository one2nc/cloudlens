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
	"s3": {
		DAO:      &dao.S3{},
		Renderer: &render.S3{},
	},
	"sg": {
		DAO:      &dao.SG{},
		Renderer: &render.SG{},
	},
}
