package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/piriya-muaithaisong/testgolang_mongo/db"
	"github.com/piriya-muaithaisong/testgolang_mongo/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type createUserRequest struct {
	Username       string `json:"username" bson:"username"`
	FistName       string `json:"first_name" bson:"first_name"`
	LastName       string `json:"last_name" bson:"last_name"`
	Email          string `json:"email" bson:"email"`
	Org            string `json:"org" bson:"org"`
	BillingAddress string `json:"billing_address" bson:"billing_address"`
	Password       string `json:"password"`
}

type userResponse struct {
	Username       string    `json:"username" bson:"username"`
	FistName       string    `json:"first_name" bson:"first_name"`
	LastName       string    `json:"last_name" bson:"last_name"`
	Email          string    `json:"email" bson:"email"`
	Org            string    `json:"org" bson:"org"`
	BillingAddress string    `json:"billing_address" bson:"billing_address"`
	CreateAt       time.Time `json:"create_at" bson:"create_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:       user.Username,
		FistName:       user.FistName,
		LastName:       user.LastName,
		Email:          user.Email,
		CreateAt:       user.CreateAt,
		Org:            user.Org,
		BillingAddress: user.BillingAddress,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	id := primitive.NewObjectID()

	user := db.User{
		ID:                id,
		Username:          req.Username,
		HashedPassword:    hashedPassword,
		FistName:          req.FistName,
		LastName:          req.LastName,
		Email:             req.Email,
		CreateAt:          time.Now(),
		PasswordChangedAt: time.Now(),
		Role:              "client",
		Org:               req.Org,
		BillingAddress:    req.BillingAddress,
	}

	err = server.store.UserClient.CreateUser(ctx, &user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	user, err = server.store.UserClient.GetUserByID(ctx, id.String())
	rep := newUserResponse(user)

	ctx.JSON(http.StatusOK, rep)

}

type loginUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginUserResponse struct {
	SessionID             primitive.ObjectID `json:"session_id"`
	AccessToken           string             `json:"access_token"`
	AccessTokenExpiresAt  time.Time          `json:"access_token_expires_at"`
	RefreshToken          string             `json:"refresh_token"`
	RefreshTokenExpiresAt time.Time          `json:"refresh_token_expires_at"`
	User                  userResponse       `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.UserClient.GetUserByUsername(ctx, req.Username)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = utils.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, accessPayload, err := server.tokenMaker.CreateToken(
		user.ID,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	refreshToken, refreshPayload, err := server.tokenMaker.CreateToken(
		user.ID,
		server.config.AccessTokenDuration,
	)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	session := db.Session{
		ID:           refreshPayload.ID,
		UserID:       user.ID,
		RefreshToken: refreshToken,
		IsBlocked:    false,
		ExpiresAt:    refreshPayload.ExpiredAt,
	}

	err = server.store.SessionClient.CreateSession(ctx, &session)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	res := loginUserResponse{
		SessionID:             session.ID,
		AccessTokenExpiresAt:  accessPayload.ExpiredAt,
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: refreshPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}

	ctx.JSON(http.StatusOK, res)

}
