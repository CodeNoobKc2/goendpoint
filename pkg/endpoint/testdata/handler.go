package testdata

import (
	"context"
	"errors"
	"fmt"
)

func MockAuth(ctx context.Context, req struct {
	AuthReq `body:"json"`
}) (resp AuthRes) {
	if req.Passwd == "bar" && req.Username == "foo" {
		resp.Success = "success"
	} else {
		resp.Error = errors.New("auth failed")
	}
	return
}

func MockListUsers(ctx context.Context, req struct {
	ListReq
	FuzzyQueryUser
}) (resp struct {
	ListUserResp `body:"json"`
}) {
	if req.PageSize == 0 {
		req.PageSize = 10
	}

	ret := make([]*User, req.PageSize)
	for i := 0; i < req.PageSize; i++ {
		if req.FuzzyQueryUser.BlurUsername == nil {
			ret[i] = &User{UserId: uint64(i), Username: fmt.Sprintf("user-%v", i)}
		} else {
			ret[i] = &User{UserId: uint64(i), Username: fmt.Sprintf("%v-%v", *req.FuzzyQueryUser.BlurUsername, i)}
		}
	}

	resp.List = ret
	resp.Tot = req.PageSize
	return
}
