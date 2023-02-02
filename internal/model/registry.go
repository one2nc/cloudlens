package model

import (
	"github.com/one2nc/cloud-lens/internal/dao"
	"github.com/one2nc/cloud-lens/internal/render"
)

var Registry = map[string]ResourceMeta{
	"EC2": {
		DAO:      &dao.EC2{},
		Renderer: &render.EC2{},
	},
	"S3": {
		DAO:      &dao.S3{},
		Renderer: &render.S3{},
	},
	"SG": {
		DAO:      &dao.SG{},
		Renderer: &render.SG{},
	},
}
