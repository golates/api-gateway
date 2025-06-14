package handlers

import (
	"context"
	"fmt"
	"github.com/golates/api-gateway/internal/models"
	"github.com/golates/api-gateway/internal/utils"
	pb "github.com/golates/api-gateway/services/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"net/http"
	"os"
)

type AuthHandler struct {
	authServiceClient pb.AuthServiceClient
}

func NewAuthHandler(authServiceURL string) *AuthHandler {
	conn, err := grpc.NewClient(authServiceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
		os.Exit(-1)
	}

	authServiceClient := pb.NewAuthServiceClient(conn)

	return &AuthHandler{
		authServiceClient: authServiceClient,
	}
}

func (ah *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	body := new(models.LoginRequest)
	if err := utils.ValidateBody(w, r, body); err != nil {
		return
	}

	res, err := ah.authServiceClient.Login(context.Background(), &pb.LoginRequest{
		Email:    body.Email,
		Password: body.Password,
	})

	if err != nil {
		errCode, errMsg := utils.ParseGRPCError(err)
		utils.WriteJSON(w, errCode, models.MessageAPIResponseError{Message: errMsg})
		return
	}

	utils.WriteJSON(w, http.StatusOK, &models.LoginResponse{Success: res.GetSuccess()})
}

func (ah *AuthHandler) OAuthGoogleLogin(w http.ResponseWriter, r *http.Request) {
	body := new(models.OAuthGoogleLoginRequest)
	if err := utils.ValidateBody(w, r, body); err != nil {
		return
	}

	res, err := ah.authServiceClient.LoginWithOAuthGoogle(context.Background(), &pb.LoginWithOAuthGoogleRequest{
		Token: body.Token,
	})

	if err != nil {
		errCode, errMsg := utils.ParseGRPCError(err)
		utils.WriteJSON(w, errCode, models.MessageAPIResponseError{Message: errMsg})
		return
	}

	utils.WriteJSON(w, http.StatusOK, &models.OAuthGoogleLoginResponse{Success: res.GetSuccess()})
}

func (ah *AuthHandler) OAuthFacebookLogin(w http.ResponseWriter, r *http.Request) {
	body := new(models.OAuthFacebookLoginRequest)
	if err := utils.ValidateBody(w, r, body); err != nil {
		return
	}

	res, err := ah.authServiceClient.LoginWithOAuthFacebook(context.Background(), &pb.LoginWithOAuthFacebookRequest{
		Token: body.Token,
	})

	if err != nil {
		errCode, errMsg := utils.ParseGRPCError(err)
		utils.WriteJSON(w, errCode, models.MessageAPIResponseError{Message: errMsg})
		return
	}

	utils.WriteJSON(w, http.StatusOK, &models.OAuthFacebookLoginResponse{Success: res.GetSuccess()})
}

func (ah *AuthHandler) CheckEmail(w http.ResponseWriter, r *http.Request) {
	body := new(models.CheckEmailRequest)
	if err := utils.ValidateBody(w, r, body); err != nil {
		return
	}

	res, err := ah.authServiceClient.CheckEmail(context.Background(), &pb.CheckEmailRequest{
		Email: body.Email,
	})

	if err != nil {
		errCode, errMsg := utils.ParseGRPCError(err)
		utils.WriteJSON(w, errCode, models.MessageAPIResponseError{Message: errMsg})
		return
	}

	utils.WriteJSON(w, http.StatusOK, &models.CheckEmailResponse{AccountExist: res.GetAccountExists()})
}

func (ah *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	body := new(models.RegisterRequest)
	if err := utils.ValidateBody(w, r, body); err != nil {
		return
	}

	fmt.Println(body)
}

func (ah *AuthHandler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	body := new(models.ForgotPasswordRequest)
	if err := utils.ValidateBody(w, r, body); err != nil {
		return
	}

	fmt.Println(body)
}
