package model

import (
	"github.com/one2nc/cloudlens/internal"
	"github.com/one2nc/cloudlens/internal/dao"
	"github.com/one2nc/cloudlens/internal/render"
)

var Registry = map[string]ResourceMeta{
	internal.LowercaseEc2: {
		DAO:      &dao.EC2{},
		Renderer: &render.EC2{},
	},
	internal.LowercaseS3: {
		DAO:      &dao.S3{},
		Renderer: &render.S3{},
	},
	internal.LowercaseSg: {
		DAO:      &dao.SG{},
		Renderer: &render.SG{},
	},
	internal.Object: {
		DAO:      &dao.BObj{},
		Renderer: &render.BObj{},
	},
	internal.LowercaseIamUser: {
		DAO:      &dao.IAMU{},
		Renderer: &render.IAMU{},
	},
	internal.LowercaseIamGroup: {
		DAO:      &dao.IAMUG{},
		Renderer: &render.IAMUG{},
	},
	internal.LowercaseIamRole: {
		DAO:      &dao.IamRole{},
		Renderer: &render.IamRole{},
	},
	internal.UserPolicy: {
		DAO:      &dao.IAMUP{},
		Renderer: &render.IamUserPloicy{},
	},
	internal.LowercaseEBS: {
		DAO:      &dao.EBS{},
		Renderer: &render.EBS{},
	},
	internal.UserGroupPolicy: {
		DAO:      &dao.IAMUGP{},
		Renderer: &render.IamUserGroupPloicy{},
	},
	internal.RolePolicy: {
		DAO:      &dao.IamRolePloicy{},
		Renderer: &render.IamRolePloicy{},
	},
	internal.GroupUsers: {
		DAO:      &dao.IamGroupUser{},
		Renderer: &render.IamGroupUser{},
	},
	internal.LowercaseEc2Snapshot: {
		DAO:      &dao.EC2S{},
		Renderer: &render.EC2S{},
	},
	internal.LowercaseEc2Image: {
		DAO:      &dao.EC2I{},
		Renderer: &render.EC2I{},
	},
	internal.LowercaseSQS: {
		DAO:      &dao.SQS{},
		Renderer: &render.SQS{},
	},
	internal.LowercaseVPC: {
		DAO:      &dao.VPC{},
		Renderer: &render.VPC{},
	},
	internal.LowercaseSubnet: {
		DAO:      &dao.Subnet{},
		Renderer: &render.Subnet{},
	},
	internal.LowercaseLamda: {
		DAO:      &dao.Lambda{},
		Renderer: &render.Lambda{},
	},
	internal.LowercaseEcsCluster: {
		DAO:      &dao.ECSClusters{},
		Renderer: &render.EcsClusters{},
	},
}
