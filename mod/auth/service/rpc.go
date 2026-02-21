package service

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net/smtp"
	"strconv"

	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authDao "github.com/GoldenSheep402/Hermes/mod/auth/dao"
	"github.com/GoldenSheep402/Hermes/mod/auth/model/codeValues"
	"github.com/GoldenSheep402/Hermes/mod/user/dao"
	"github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/pkg/auth"
	authV1 "github.com/GoldenSheep402/Hermes/pkg/proto/auth/v1"
	"github.com/GoldenSheep402/Hermes/pkg/randx"
	"github.com/GoldenSheep402/Hermes/pkg/utils/check"
	"github.com/GoldenSheep402/Hermes/pkg/utils/crypto"
)

var _ authV1.AuthServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	authV1.UnimplementedAuthServiceServer
}

// RegisterSendEmail sends an email verification code.
func (s *S) RegisterSendEmail(ctx context.Context, req *authV1.RegisterSendEmailRequest) (*authV1.RegisterSendEmailResponse, error) {
	// TODO: Read SMTP settings from system KV store
	// For now, use placeholder values until system DAO is updated
	smtpEnable := false
	if !smtpEnable {
		return nil, status.Error(codes.PermissionDenied, "SMTP is not enabled")
	}

	smtpHost := ""
	smtpPort := "465"
	smtpUser := ""
	smtpPassword := ""
	senderEmail := smtpUser
	smtpTO := req.Email
	subject := "Subject: 欢迎来到HERMES\r\n"
	code := randx.String(6)
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h1>你好!</h1>
			<p>欢迎来到 HERMES！您的验证码是 <strong>%s</strong>，请在页面中输入此验证码。</p>
		</body>
		</html>
	`, code)

	msg := []byte(subject + mime + body)
	smtpAuth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         smtpHost,
	}

	conn, err := tls.Dial("tcp", smtpHost+":"+smtpPort, tlsconfig)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create TLS connection")
	}

	client, err := smtp.NewClient(conn, smtpHost)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to create SMTP client")
	}
	defer client.Close()

	if err = client.Auth(smtpAuth); err != nil {
		return nil, status.Error(codes.Internal, "SMTP authentication failed")
	}
	if err = client.Mail(senderEmail); err != nil {
		return nil, status.Error(codes.Internal, "Failed to set sender email")
	}
	if err = client.Rcpt(smtpTO); err != nil {
		return nil, status.Error(codes.Internal, "Failed to set recipient email")
	}

	wc, err := client.Data()
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to write email data")
	}
	if _, err = wc.Write(msg); err != nil {
		return nil, status.Error(codes.Internal, "Failed to write message")
	}
	if err = wc.Close(); err != nil {
		return nil, status.Error(codes.Internal, "Failed to close write connection")
	}
	if err = client.Quit(); err != nil {
		return nil, status.Error(codes.Internal, "Failed to close SMTP connection")
	}

	if err = authDao.Code.SetCodeWithEmail(ctx, req.Email, code); err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &authV1.RegisterSendEmailResponse{}, nil
}

// RegisterWithEmail registers a new user with email.
func (s *S) RegisterWithEmail(ctx context.Context, req *authV1.RegisterWithEmailRequest) (*authV1.RegisterWithEmailResponse, error) {
	// TODO: Check register enable/smtp enable from system KV store

	if req.EmailToken != "" {
		_status, err := authDao.Code.CheckCodeWithAttempts(ctx, req.Email, req.EmailToken)
		if err != nil {
			return nil, status.Error(codes.Internal, "Internal error")
		}
		switch _status {
		case codeValues.Wrong:
			return nil, status.Error(codes.InvalidArgument, "Email token error")
		case codeValues.TooManyAttempts:
			return nil, status.Error(codes.InvalidArgument, "Too many attempts")
		case codeValues.Right:
			// OK
		}
	}

	email := req.Email
	password := req.Password

	if !check.VerifyEmailFormat(email) {
		return nil, status.Error(codes.InvalidArgument, "Email format error")
	}
	if len(password) < 6 {
		return nil, status.Error(codes.InvalidArgument, "Password too short")
	}
	if req.Username == "" {
		return nil, status.Error(codes.InvalidArgument, "Username is empty")
	}

	salt, err := crypto.GenerateSalt(16)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	newUser := &model.User{
		Username:  req.Username,
		Email:     email,
		Password:  crypto.Md5CryptoWithSalt(password, salt),
		Salt:      salt,
		Passkey:   ulid.Make().String(),
		IsAdmin:   false,
		IsEnabled: true,
	}

	_, err = dao.User.CreateUser(ctx, newUser)
	if err != nil {
		if errors.Is(err, dao.ErrEmailAlreadyUsed) {
			return nil, status.Error(codes.InvalidArgument, "Email already used")
		}
		if errors.Is(err, dao.ErrUsernameAlreadyUsed) {
			return nil, status.Error(codes.InvalidArgument, "Username already used")
		}
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &authV1.RegisterWithEmailResponse{}, nil
}

func (s *S) Login(ctx context.Context, req *authV1.LoginRequest) (*authV1.LoginResponse, error) {
	email := req.Email
	password := req.Password

	if !check.VerifyEmailFormat(email) {
		return nil, status.Error(codes.InvalidArgument, "Email format error")
	}

	user, err := dao.User.GetByEmail(ctx, email)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "Email error")
	}

	if crypto.Md5CryptoWithSalt(password, user.Salt) != user.Password {
		return nil, status.Error(codes.InvalidArgument, "Password error")
	}

	refreshToken, err := auth.GenToken(auth.Info{
		UID:            user.ID,
		IsRefreshToken: true,
	}, auth.RefreshTokenExpireIn)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	accessToken, err := auth.GenToken(auth.Info{
		UID:            user.ID,
		IsRefreshToken: false,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &authV1.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *S) RefreshToken(_ context.Context, request *authV1.RefreshTokenRequest) (*authV1.RefreshTokenResponse, error) {
	entity, err := auth.ParseToken(request.RefreshToken)
	if err != nil || !entity.Info.IsRefreshToken {
		return nil, status.Error(codes.Unauthenticated, "invalid refresh token")
	}

	accessToken, err := auth.GenToken(auth.Info{
		UID: entity.Info.UID,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authV1.RefreshTokenResponse{
		AccessToken:  accessToken,
		RefreshToken: request.RefreshToken,
	}, nil
}

// smtpPortStr is a helper to convert port to string.
func smtpPortStr(port int) string {
	return strconv.Itoa(port)
}
