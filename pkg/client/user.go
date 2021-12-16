package client

import (
	"context"
	"git.internal.yunify.com/qxp/misc/client"
	"net/http"
)

const (
	host          = "http://org/api/v1/org"
	usersInfoURI  = "/usersInfo"
	departmentURI = "/depByIDs"
	userDEPIDURI  = "/otherGetUserList"
)

// User User
type User interface {
	GetInfo(ctx context.Context, userIDs ...string) ([]UserInfo, error)
	GetDepartment(ctx context.Context, ids ...string) ([]Department, error)
	GetUsersByDEPID(ctx context.Context, depID string, includeChildDEPChild, page, limit int) ([]UserInfo, error)
}

type user struct {
	client http.Client
}

// NewUser NewUser
func NewUser(conf Config) User {
	return &user{
		client: New(conf),
	}
}

// UserInfo 用户信息
type UserInfo struct {
	ID          string `json:"id"`
	UserName    string `json:"userName"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	Address     string `json:"address"`
	LeaderID    string `json:"leaderID"`
	CompanyID   string `json:"companyID"`
	Avatar      string `json:"avatar"`
	IsDEPLeader int    `json:"isDEPLeader,omitempty"` //是否部门领导 1是，-1不是
	DEP         struct {
		ID                 string `json:"id"`
		DepartmentName     string `json:"departmentName"`
		DepartmentLeaderID string `json:"departmentLeaderID"`
		UseStatus          int    `json:"useStatus"`
		PID                string `json:"pid"`
		SuperPID           string `json:"superID"`
		CompanyID          string `json:"companyID"`
		Grade              int    `json:"grade"`
	} `json:"dep"`
}

// GetInfo GetInfo
func (u *user) GetInfo(ctx context.Context, userIDs ...string) ([]UserInfo, error) {
	params := struct {
		IDS []string `json:"ids"`
	}{
		IDS: userIDs,
	}

	userInfo := make([]UserInfo, 0)
	err := POST(ctx, &u.client, host+usersInfoURI, params, &userInfo)
	return userInfo, err
}

// Department Department
type Department struct {
	ID                 string `json:"id"`
	DepartmentName     string `json:"departmentName"`
	DepartmentLeaderID string `json:"departmentLeaderID"`
	UseStatus          int    `json:"useStatus"`
	PID                string `json:"pid"`
	SuperPID           string `json:"superID"`
	CompanyID          string `json:"companyID"`
	Grade              int    `json:"grade"`
}

// GetDepartment GetDepartment
func (u *user) GetDepartment(ctx context.Context, ids ...string) ([]Department, error) {
	params := struct {
		IDS []string `json:"ids"`
	}{
		IDS: ids,
	}

	deparment := make([]Department, 0)
	err := client.POST(ctx, &u.client, host+departmentURI, params, &deparment)
	return deparment, err
}

// GetUsersByDEPID GetUsersByDEPID
func (u *user) GetUsersByDEPID(ctx context.Context, depID string, includeChildDEPChild, page, limit int) ([]UserInfo, error) {
	params := struct {
		DepID                string `json:"depID"`
		IncludeChildDEPChild int    `json:"includeChildDEPChild"`
		Page                 int    `json:"page"`
		Limit                int    `json:"limit"`
	}{
		DepID:                depID,
		IncludeChildDEPChild: includeChildDEPChild,
		Page:                 page,
		Limit:                limit,
	}

	res := make([]UserInfo, 0)
	err := client.POST(ctx, &u.client, host+userDEPIDURI, params, &res)
	return res, err
}
