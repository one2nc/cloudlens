package aws

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
)

func GetUsers(sess session.Session) []IAMUSerResp {
	iamSrv := iam.New(&sess)
	result, err := iamSrv.ListUsers(&iam.ListUsersInput{})
	if err != nil {
		fmt.Println("Error in fetching Iam Users: ", " err: ", err)
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

func GetUserGroups(sess session.Session) []*iam.Group {
	iamSrv := iam.New(&sess)
	result, err := iamSrv.ListGroups(&iam.ListGroupsInput{})
	if err != nil {
		fmt.Println("Error in fetching Iam Groups: ", " err: ", err)
		return nil
	}
	return result.Groups
}

func GetGroupUsers(sess session.Session, grpName string) []*iam.User {
	iamSrv := iam.New(&sess)
	result, err := iamSrv.GetGroup(&iam.GetGroupInput{
		GroupName: &grpName,
	})
	if err != nil {
		fmt.Println("Error in fetching Iam users of the Group: ", grpName, " err: ", err)
		return nil
	}
	return result.Users
}

func GetPoliciesOfGrp(sess session.Session, grpName string) []*iam.AttachedPolicy {
	imaSrv := iam.New(&sess)
	result, err := imaSrv.ListAttachedGroupPolicies(&iam.ListAttachedGroupPoliciesInput{
		GroupName: &grpName,
	})
	if err != nil {
		fmt.Println("Error in fetching Iam policies of the Group: ", grpName, " err: ", err)
		return nil
	}
	return result.AttachedPolicies
}

// If a user belong to a Group then we can't see the user's attached policy here,
// their policies are governed on the top of the group
func GetPoliciesOfUser(sess session.Session, usrName string) []*iam.AttachedPolicy {
	imaSrv := iam.New(&sess)
	result, err := imaSrv.ListAttachedUserPolicies(&iam.ListAttachedUserPoliciesInput{
		UserName: &usrName,
	})
	if err != nil {
		fmt.Println("Error in fetching Iam policies of the User: ", usrName, " err: ", err)
		return nil
	}
	return result.AttachedPolicies
}

func GetIamRoles(sess session.Session) []*iam.Role {
	iamSrv := iam.New(&sess)
	result, err := iamSrv.ListRoles(&iam.ListRolesInput{})
	if err != nil {
		fmt.Println("Error in fetching Iam Roles: ", " err: ", err)
		return nil
	}
	return result.Roles
}
