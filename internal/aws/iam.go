package aws

import (
	"context"
	"fmt"
	"time"

	awsV2 "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/one2nc/cloudlens/internal/config"
	"github.com/rs/zerolog/log"
)

func GetUsers(cfg awsV2.Config) []IAMUSerResp {
	iamSrv := iam.NewFromConfig(cfg)
	result, err := iamSrv.ListUsers(context.Background(), &iam.ListUsersInput{})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching Iam users: ,  err: %v", err))
		return nil
	}
	var users []IAMUSerResp
	for _, u := range result.Users {
		launchTime := u.CreateDate
		localZone, err := config.GetLocalTimeZone() // Empty string loads the local timezone
		if err != nil {
			fmt.Println("Error loading local timezone:", err)
			return nil
		}
		loc, _ := time.LoadLocation(localZone)
		IST := launchTime.In(loc)
		user := &IAMUSerResp{
			UserId:       *u.UserId,
			UserName:     *u.UserName,
			ARN:          *u.Arn,
			CreationTime: IST.Format("Mon Jan _2 15:04:05 2006"),
		}
		users = append(users, *user)
	}
	return users
}

func GetUserGroups(cfg awsV2.Config) []IAMUSerGroupResp {
	iamSrv := iam.NewFromConfig(cfg)
	result, err := iamSrv.ListGroups(context.Background(), &iam.ListGroupsInput{})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching Iam Groups: , err: %v", err))
		return nil
	}
	var userGroups []IAMUSerGroupResp
	for _, u := range result.Groups {
		userGroup := &IAMUSerGroupResp{
			GroupId:   *u.GroupId,
			GroupName: *u.GroupName,
			ARN:       *u.Arn,
			// CreationTime: fmt.Sprintf("%v",*u.CreateDate), Created time is not coming from sdk
		}
		userGroups = append(userGroups, *userGroup)
	}
	return userGroups
}

func GetGroupUsers(cfg awsV2.Config, grpName string) []IAMUSerResp {
	iamSrv := iam.NewFromConfig(cfg)
	result, err := iamSrv.GetGroup(context.Background(), &iam.GetGroupInput{
		GroupName: &grpName,
	})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching Iam users of the Group: %s,  err: %v", grpName, err))
		return nil
	}

	var users []IAMUSerResp
	for _, u := range result.Users {
		user := &IAMUSerResp{
			UserId:   *u.UserId,
			UserName: *u.UserName,
			ARN:      *u.Arn,
			// CreationTime: , Created time is not coming from sdk
		}
		users = append(users, *user)
	}
	return users
}

func GetPoliciesOfGrp(cfg awsV2.Config, grpName string) []IAMUSerGroupPolicyResponse {
	imaSrv := iam.NewFromConfig(cfg)
	result, err := imaSrv.ListAttachedGroupPolicies(context.Background(), &iam.ListAttachedGroupPoliciesInput{
		GroupName: &grpName,
	})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching Iam policies of the Group: %s,  err: %v", grpName, err))
		return nil
	}
	var grpPolicies []IAMUSerGroupPolicyResponse
	for _, up := range result.AttachedPolicies {
		grpPolicy := &IAMUSerGroupPolicyResponse{
			PolicyArn:  *up.PolicyArn,
			PolicyName: *up.PolicyName,
		}
		grpPolicies = append(grpPolicies, *grpPolicy)
	}
	return grpPolicies
}

// If a user belong to a Group then we can't see the user's attached policy here,
// their policies are governed on the top of the group
func GetPoliciesOfUser(cfg awsV2.Config, usrName string) []IAMUSerPolicyResponse {
	imaSrv := iam.NewFromConfig(cfg)
	result, err := imaSrv.ListAttachedUserPolicies(context.Background(), &iam.ListAttachedUserPoliciesInput{
		UserName: &usrName,
	})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching Iam policies of the User: %s,  err: %v", usrName, err))
		return nil
	}
	var usersPolicy []IAMUSerPolicyResponse
	for _, up := range result.AttachedPolicies {
		userPolicy := &IAMUSerPolicyResponse{
			PolicyArn:  *up.PolicyArn,
			PolicyName: *up.PolicyName,
		}
		usersPolicy = append(usersPolicy, *userPolicy)
	}
	return usersPolicy
}

func GetIamRoles(cfg awsV2.Config) []IamRoleResp {
	iamSrv := iam.NewFromConfig(cfg)
	result, err := iamSrv.ListRoles(context.Background(), &iam.ListRolesInput{})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching Iam roles,  err: %v", err))
		return nil
	}
	var roles []IamRoleResp
	for _, r := range result.Roles {
		launchTime := r.CreateDate
		localZone, err := config.GetLocalTimeZone() // Empty string loads the local timezone
		if err != nil {
			fmt.Println("Error loading local timezone:", err)
			return nil
		}
		loc, _ := time.LoadLocation(localZone)
		IST := launchTime.In(loc)
		role := &IamRoleResp{
			RoleId:       *r.RoleId,
			RoleName:     *r.RoleName,
			ARN:          *r.Arn,
			CreationTime: IST.Format("Mon Jan _2 15:04:05 2006"),
		}
		roles = append(roles, *role)
	}
	return roles
}

func GetPoliciesOfRoles(cfg awsV2.Config, roleName string) []IamRolePolicyResponse {
	imaSrv := iam.NewFromConfig(cfg)
	result, err := imaSrv.ListAttachedRolePolicies(context.Background(), &iam.ListAttachedRolePoliciesInput{
		RoleName: &roleName,
	})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching Iam policies of the User: %v  err: %v", roleName, err))
		return nil
	}
	var Policies []IamRolePolicyResponse
	for _, up := range result.AttachedPolicies {
		userPolicy := &IamRolePolicyResponse{
			PolicyArn:  *up.PolicyArn,
			PolicyName: *up.PolicyName,
		}
		Policies = append(Policies, *userPolicy)
	}
	return Policies
}
