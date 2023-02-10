package aws

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/rs/zerolog/log"
)

func GetUsers(sess session.Session) []IAMUSerResp {
	iamSrv := iam.New(&sess)
	result, err := iamSrv.ListUsers(&iam.ListUsersInput{})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching Iam users: ,  err: %v", err))
		return nil
	}
	var users []IAMUSerResp
	for _, u := range result.Users {
		launchTime := u.CreateDate
		loc, _ := time.LoadLocation("Asia/Kolkata")
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

func GetUserGroups(sess session.Session) []IAMUSerGroupResp {
	iamSrv := iam.New(&sess)
	result, err := iamSrv.ListGroups(&iam.ListGroupsInput{})
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

func GetGroupUsers(sess session.Session, grpName string) []IAMUSerResp {
	iamSrv := iam.New(&sess)
	result, err := iamSrv.GetGroup(&iam.GetGroupInput{
		GroupName: aws.String(grpName),
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

func GetPoliciesOfGrp(sess session.Session, grpName string) []IAMUSerGroupPolicyResponse {
	imaSrv := iam.New(&sess)
	result, err := imaSrv.ListAttachedGroupPolicies(&iam.ListAttachedGroupPoliciesInput{
		GroupName: aws.String(grpName),
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
func GetPoliciesOfUser(sess session.Session, usrName string) []IAMUSerPolicyResponse {
	imaSrv := iam.New(&sess)
	result, err := imaSrv.ListAttachedUserPolicies(&iam.ListAttachedUserPoliciesInput{
		UserName: aws.String(usrName),
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

func GetIamRoles(sess session.Session) []IamRoleResp {
	iamSrv := iam.New(&sess)
	result, err := iamSrv.ListRoles(&iam.ListRolesInput{})
	if err != nil {
		log.Info().Msg(fmt.Sprintf("Error in fetching Iam roles,  err: %v", err))
		return nil
	}
	var roles []IamRoleResp
	for _, r := range result.Roles {
		launchTime := r.CreateDate
		loc, _ := time.LoadLocation("Asia/Kolkata")
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

func GetPoliciesOfRoles(sess session.Session, roleName string) []IamRolePolicyResponse {
	imaSrv := iam.New(&sess)
	result, err := imaSrv.ListAttachedRolePolicies(&iam.ListAttachedRolePoliciesInput{
		RoleName: aws.String(roleName),
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
