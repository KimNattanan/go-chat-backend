package v1

import (
	"net/http"

	"github.com/KimNattanan/go-chat-backend/internal/chat/entity"
	"github.com/KimNattanan/go-chat-backend/pkg/responses"
	"github.com/google/uuid"
	"github.com/labstack/echo/v5"
)

func (r *V1) findMembershipByRoomID(c *echo.Context) error {
	ctx := c.Request().Context()
	roomID := c.Param("roomID")

	memberships, err := r.membershipUseCase.FindByRoomID(ctx, roomID)
	if err != nil {
		r.l.Error(err, "rest - v1 - findMembershipByRoomID")
		return responses.ErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, toMembershipResponseList(memberships))
}

func (r *V1) findMembershipByUserID(c *echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Param("userID")

	memberships, err := r.membershipUseCase.FindByRoomID(ctx, userID)
	if err != nil {
		r.l.Error(err, "rest - v1 - findMembershipByUserID")
		return responses.ErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, toMembershipResponseList(memberships))
}

func (r *V1) findMembershipByRoomIDAndUserID(c *echo.Context) error {
	ctx := c.Request().Context()
	roomID := c.Param("roomID")
	userID := c.Param("userID")

	membership, err := r.membershipUseCase.FindByRoomIDAndUserID(ctx, roomID, userID)
	if err != nil {
		r.l.Error(err, "rest - v1 - findMembership")
		return responses.ErrorResponse(c, err)
	}

	return c.JSON(http.StatusOK, toMembershipResponse(membership))
}

func (r *V1) createMembership(c *echo.Context) error {
	ctx := c.Request().Context()

	roomIDStr := c.Param("roomID")
	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		r.l.Error(err, "rest - v1 - createMembership")
		return responses.ErrorResponse(c, err)
	}

	userIDStr := c.Param("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		r.l.Error(err, "rest - v1 - createMembership")
		return responses.ErrorResponse(c, err)
	}

	membership := &entity.Membership{
		RoomID: roomID,
		UserID: userID,
	}
	if err := r.membershipUseCase.Create(ctx, membership); err != nil {
		r.l.Error(err, "rest - v1 - createMembership")
		return responses.ErrorResponse(c, err)
	}

	return c.JSON(http.StatusCreated, toMembershipResponse(membership))
}

func (r *V1) deleteMembershipByUserID(c *echo.Context) error {
	ctx := c.Request().Context()
	userID := c.Param("userID")

	if err := r.membershipUseCase.DeleteByUserID(ctx, userID); err != nil {
		r.l.Error(err, "rest - v1 - deleteMembershipByUserID")
		return responses.ErrorResponse(c, err)
	}

	return responses.MessageResponse(c, http.StatusOK, "memberships deleted")
}

func (r *V1) deleteMembershipByRoomIDAndUserID(c *echo.Context) error {
	ctx := c.Request().Context()
	roomID := c.Param("roomID")
	userID := c.Param("userID")

	if err := r.membershipUseCase.Delete(ctx, roomID, userID); err != nil {
		r.l.Error(err, "rest - v1 - deleteMembership")
		return responses.ErrorResponse(c, err)
	}

	return responses.MessageResponse(c, http.StatusOK, "membership deleted")
}
