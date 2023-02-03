package api

import (
	"database/sql"
	"errors"
	db "github.com/SemmiDev/simpeg/db/mysql"
	"github.com/SemmiDev/simpeg/util"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

type createUserRequest struct {
	Name     string `json:"name"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Photo    string `json:"photo"`
}

type userResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Email string `json:"email"`
	Photo string `json:"photo"`
}

type updateUserPasswordRequest struct {
	ID          string `json:"id"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func (server *Server) changePassword(ctx *gin.Context) {
	var req updateUserPasswordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	if err := util.CheckPassword(req.OldPassword, user.Password); err != nil {
		err = errors.New("Password salah")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.NewPassword)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.UpdateUserPasswordParams{
		ID:       req.ID,
		Password: hashedPassword,
	}

	err = server.store.UpdateUserPassword(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, nil)
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	arg := db.CreateUserParams{
		ID:       uuid.New().String(),
		Name:     req.Name,
		Role:     req.Role,
		Email:    req.Email,
		Password: hashedPassword,
		Photo:    req.Photo,
	}

	err = server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := userResponse{
		ID:    arg.ID,
		Name:  req.Name,
		Role:  req.Role,
		Email: req.Email,
		Photo: req.Photo,
	}

	ctx.JSON(http.StatusCreated, rsp)
}

type loginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginUserResponse struct {
	AccessToken string       `json:"access_token"`
	User        userResponse `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.New("Akun tidak ditemukan")
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		err = errors.New("Server sedang bermasalah")
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	err = util.CheckPassword(req.Password, user.Password)
	if err != nil {
		err = errors.New("Email atau password salah")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}

	accessToken, _, err := server.tokenMaker.CreateToken(
		user.Email,
		server.config.AccessTokenDuration,
	)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := loginUserResponse{
		AccessToken: accessToken,
		User: userResponse{
			ID:    user.ID,
			Name:  user.Name,
			Role:  user.Role,
			Email: user.Email,
			Photo: user.Photo,
		},
	}

	ctx.JSON(http.StatusOK, rsp)
}

type updateUserRequest struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Role  string `json:"role"`
	Email string `json:"email"`
	Photo string `json:"photo"`
}

func (server *Server) updateUserProfile(ctx *gin.Context) {
	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.UpdateUserProfileParams{
		ID:    req.ID,
		Name:  req.Name,
		Role:  req.Role,
		Email: req.Email,
		Photo: req.Photo,
	}

	err := server.store.UpdateUserProfile(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	rsp := userResponse{
		ID:    arg.ID,
		Name:  req.Name,
		Role:  req.Role,
		Email: req.Email,
		Photo: req.Photo,
	}

	ctx.JSON(http.StatusCreated, rsp)
}
