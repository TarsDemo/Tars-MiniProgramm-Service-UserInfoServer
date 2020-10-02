package main

import (
	"context"
	"errors"
	"strconv"
	"time"

	"LifeService"
)

// UserInfoServiceImp servant implementation
type UserInfoServiceImp struct {
	app     *LifeService.DataService
	objName string
}

// Init servant init
func (imp *UserInfoServiceImp) Init() error {
	//initialize servant here:
	imp.app = new(LifeService.DataService)
	imp.objName = "LifeService.DataServer.DataServiceObj"
	comm.StringToProxy(imp.objName, imp.app)
	return nil
}

// Destroy servant destory
func (imp *UserInfoServiceImp) Destroy() {
	//destroy servant here:
	//...
}

//HasPhone 判断号码是否存在
func (imp *UserInfoServiceImp) HasPhone(ctx context.Context, phone string, phoneExist *bool) (int32, error) {
	_, err := imp.app.HasPhone(phone, phoneExist)
	if err != nil {
		SLOG.Error("Call Remote DataServer::HasPhone error: ", err.Error())
		return -1, err
	}
	return 0, nil
}

//SignUp Create a new account
func (imp *UserInfoServiceImp) SignUp(ctx context.Context, wxID string, userInfo *LifeService.UserInfo, errCode *LifeService.ErrorCode) (int32, error) {
	// 判断号码是否存在
	var hasPhone bool
	_, err1 := imp.app.HasPhone(userInfo.Phone, &hasPhone)
	if err1 != nil {
		SLOG.Error("Create user error with error message: ", err1)
		*errCode = LifeService.ErrorCode_SERVERERROR
		return 0, nil
	}
	if hasPhone {
		*errCode = LifeService.ErrorCode_PHONEEXIST
		return 0, nil
	}

	// 获取当前时间, 据说这个时间是Go语言诞生的时间
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	userInfo.Group = 0
	userInfo.Registration_time = currentTime
	// 创建用户
	iRet, err := imp.app.CreateUser(wxID, userInfo)
	if err != nil {
		SLOG.Error("Create user error with error message: ", err)
		*errCode = LifeService.ErrorCode_SERVERERROR
	} else {
		if iRet == 0 {
			SLOG.Debug("Create success")
			*errCode = LifeService.ErrorCode_SUCCESS
		} else {
			SLOG.Debug("Create user fail: user exist")
			*errCode = LifeService.ErrorCode_USEREXIST
		}
	}
	return 0, nil
}

//SignIn 登录，获取用户信息
func (imp *UserInfoServiceImp) SignIn(ctx context.Context, wxID string, userInfo *LifeService.UserInfo, errCode *LifeService.ErrorCode) (int32, error) {
	var HasUser bool
	SLOG.Debug("SignIn")
	// 判断用户是否存在
	_, err := imp.app.HasUser(wxID, &HasUser)

	if err != nil {
		*errCode = LifeService.ErrorCode_SERVERERROR
		return 0, err
	}
	// 如果用户存在, 返回用户信息; 否则返回用户不存在错误USERNOTEXIST
	if HasUser {
		_, err := imp.app.GetUserInfo(wxID, userInfo)
		if err != nil {
			SLOG.Error("Call error: ", err)
			*errCode = LifeService.ErrorCode_SERVERERROR
		} else {
			SLOG.Debug("Call success")
			*errCode = LifeService.ErrorCode_SUCCESS
		}
		return 0, nil
	}
	*errCode = LifeService.ErrorCode_USERNOTEXIST
	return 0, nil
}

//GetGroupList 获取权限组列表
func (imp *UserInfoServiceImp) GetGroupList(ctx context.Context, groupInfo *map[int32]string) (int32, error) {
	SLOG.Debug("getGroupInfo")
	_, err := imp.app.GetGroupInfo(groupInfo)
	return 0, err
}

//IsClubManager 判断是否是社团管理员
func (imp *UserInfoServiceImp) IsClubManager(ctx context.Context, wxID string, clubID string, isClubManager *bool) (int32, error) {
	var count int32
	_, err := imp.app.GetClubManagerCount(wxID, clubID, &count)

	if err != nil {
		SLOG.Error("UserInfoServer::IsClubManager error")
	} else {
		SLOG.Debug("UserInfoServer::IsClubManager: count=" + strconv.Itoa(int(count)))
		if count > 0 {
			*isClubManager = true
		} else if count == 0 {
			*isClubManager = false
		} else {
			*isClubManager = false
			return -1, errors.New("Get user manager info failed")
		}
	}

	return 0, nil
}

//IsInClub 判断用户是否在社团中或已经申请社团
func (imp *UserInfoServiceImp) IsInClub(ctx context.Context, wxID string, clubID string, justInClub bool, isIn *bool) (int32, error) {
	var applyStatus int32
	if justInClub {
		applyStatus = 1
	} else {
		applyStatus = -2
	}
	var count int32
	_, err := imp.app.GetApplyCount(wxID, clubID, applyStatus, &count)

	if err != nil {
		SLOG.Error("UserInfoServer::IsInClub error")
		return -1, err
	}
	SLOG.Debug("UserInfoServer::IsInClub: " + strconv.Itoa(int(count)))
	if count > 0 {
		*isIn = true
	} else if count == 0 {
		*isIn = false
	} else {
		*isIn = false
		return -1, errors.New("Get user club info failed")
	}
	return 0, nil
}

//IsAppliedActivity 判断用户是否已经参加活动
func (imp *UserInfoServiceImp) IsAppliedActivity(ctx context.Context, wxID string, activityID string, isApplied *bool) (int32, error) {
	var count int32
	_, err := imp.app.GetRecordCount(wxID, activityID, &count)

	if err != nil {
		SLOG.Error("UserInfoServer::IsAppliedActivity error")
		return -1, err
	}
	SLOG.Debug("UserInfoServer::IsAppliedActivity " + strconv.Itoa(int(count)))
	if count > 0 {
		*isApplied = true
	} else if count == 0 {
		*isApplied = false
	} else {
		*isApplied = false
		return -1, errors.New("Get Info for applied activity error")
	}
	return 0, nil
}

//Test test
func (imp *UserInfoServiceImp) Test(ctx context.Context, testStr *string) (int32, error) {
	*testStr = "Test Successfull"
	return 0, nil
}
