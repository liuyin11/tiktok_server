// Code generated by hertz generator.

package api

import (
	"context"

	interaction "github.com/PCBismarck/tiktok_server/cmd/api/biz/model/interaction"
	"github.com/PCBismarck/tiktok_server/cmd/api/biz/model/shared"
	"github.com/PCBismarck/tiktok_server/cmd/api/biz/mw"
	"github.com/PCBismarck/tiktok_server/cmd/api/biz/rpc"
	"github.com/PCBismarck/tiktok_server/cmd/favorite/kitex_gen/favorite"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

// FavrotieAction .
// @router /douyin/favorite/action/ [POST]
func FavrotieAction(ctx context.Context, c *app.RequestContext) {
	var err error
	var req interaction.FavoriteActionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	mw.JwtMiddleware.MiddlewareFunc()(ctx, c)
	user, ok := c.Get(mw.JwtMiddleware.IdentityKey)
	if !ok {
		return
	}
	uid := user.(*shared.User).ID
	resp, _ := rpc.FavoriteAction(ctx, &favorite.FavoriteActionRequest{
		UserId:     uid,
		VideoId:    req.VideoId,
		ActionType: req.ActionType,
	})
	c.JSON(consts.StatusOK, resp)
}

// FavoriteList .
// @router /douyin/favorite/list/ [GET]
func FavoriteList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req interaction.FavoriteListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}
	mw.JwtMiddleware.MiddlewareFunc()(ctx, c)
	_, ok := c.Get(mw.JwtMiddleware.IdentityKey)
	if !ok {
		return
	}
	resp, _ := rpc.FavoriteList(ctx, &favorite.FavoriteListRequest{
		UserId: req.UserId,
	})

	c.JSON(consts.StatusOK, resp)
}

// CommentAction .
// @router /douyin/comment/action/ [POST]
func CommentAction(ctx context.Context, c *app.RequestContext) {
	var err error
	var req interaction.CommentActionRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	// 鉴权并获取uid
	// mw.JwtMiddleware.MiddlewareFunc()(ctx, c)
	// user, ok := c.Get(mw.JwtMiddleware.IdentityKey)
	// if !ok {
	// 	return
	// }
	// uid := user.(*shared.User).ID

	resp := new(interaction.CommentActionResponse)

	c.JSON(consts.StatusOK, resp)
}

// CommentList .
// @router /douyin/comment/list/ [GET]
func CommentList(ctx context.Context, c *app.RequestContext) {
	var err error
	var req interaction.CommentListRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		c.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp := new(interaction.CommentListResponse)

	c.JSON(consts.StatusOK, resp)
}
