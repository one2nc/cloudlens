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
	"OBJ": {
		DAO:      &dao.BObj{},
		Renderer: &render.BObj{},
	},
	"iam:u": {
		DAO:      &dao.IAMU{},
		Renderer: &render.IAMU{},
	},
	"iam:g": {
		DAO:      &dao.IAMUG{},
		Renderer: &render.IAMUG{},
	},
	"iam:r": {
		DAO:      &dao.IamRole{},
		Renderer: &render.IamRole{},
	},
	"User Policy": {
		DAO:      &dao.IAMUP{},
		Renderer: &render.IamUserPloicy{},
	},
	"ebs": {
		DAO:      &dao.EBS{},
		Renderer: &render.EBS{},
	},
	"User Group Policy": {
		DAO:      &dao.IAMUGP{},
		Renderer: &render.IamUserGroupPloicy{},
	},
	"Role Policy": {
		DAO:      &dao.IamRolePloicy{},
		Renderer: &render.IamRolePloicy{},
	},
	"Group Users": {
		DAO:      &dao.IamGroupUser{},
		Renderer: &render.IamGroupUser{},
	},
	"ec2:s": {
		DAO:      &dao.EC2S{},
		Renderer: &render.EC2S{},
	},
	"ec2:i": {
		DAO:      &dao.EC2I{},
		Renderer: &render.EC2I{},
	},

	"sqs": {
		DAO:      &dao.SQS{},
		Renderer: &render.SQS{},
	},

	"vpc": {
		DAO:      &dao.VPC{},
		Renderer: &render.VPC{},
	},
	"lambda": {
		DAO:      &dao.Lambda{},
		Renderer: &render.Lambda{},
	},
}
