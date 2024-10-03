package service

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	authDao "github.com/GoldenSheep402/Hermes/mod/auth/dao"
	"github.com/GoldenSheep402/Hermes/mod/auth/model/codeValues"
	systemDao "github.com/GoldenSheep402/Hermes/mod/system/dao"
	"github.com/GoldenSheep402/Hermes/mod/user/dao"
	userDao "github.com/GoldenSheep402/Hermes/mod/user/dao"
	"github.com/GoldenSheep402/Hermes/mod/user/model"
	"github.com/GoldenSheep402/Hermes/mod/user/model/bindType"
	"github.com/GoldenSheep402/Hermes/pkg/auth"
	authV1 "github.com/GoldenSheep402/Hermes/pkg/proto/auth/v1"
	"github.com/GoldenSheep402/Hermes/pkg/randx"
	"github.com/GoldenSheep402/Hermes/pkg/utils/check"
	"github.com/GoldenSheep402/Hermes/pkg/utils/crypto"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/smtp"
	"strconv"
)

var _ authV1.AuthServiceServer = (*S)(nil)

type S struct {
	Log *zap.SugaredLogger
	authV1.UnimplementedAuthServiceServer
}

// RegisterSendEmail TODO: SMTP
func (s *S) RegisterSendEmail(ctx context.Context, req *authV1.RegisterSendEmailRequest) (*authV1.RegisterSendEmailResponse, error) {
	settings, _, _, err := systemDao.Setting.GetSettings(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	if !settings.SmtpEnable {
		return nil, status.Error(codes.PermissionDenied, "Smtp is not allowed")
	}

	senderEmail := settings.SmtpUser
	smtpHost := settings.SmtpHost
	smtpPort := strconv.Itoa(settings.SmtpPort)
	smtpUser := settings.SmtpUser
	smtpPassword := settings.SmtpPass
	smtpTO := req.Email
	subject := "Subject: 欢迎来到HERMES\r\n"
	code := randx.String(6)
	mime := "MIME-version: 1.0;\r\nContent-Type: text/html; charset=\"UTF-8\";\r\n\r\n"
	// TODO: HTML template
	body := fmt.Sprintf(`
		<html>
		<body>
			<h1>你好!</h1>
			<p>欢迎来到 HERMES！您的验证码是 <strong>%s</strong>，请在页面中输入此验证码。</p>
		</body>
		</html>
	`, code)

	msg := []byte(subject + mime + body)

	auth := smtp.PlainAuth("", smtpUser, smtpPassword, smtpHost)

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

	if err = client.Auth(auth); err != nil {
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

	_, err = wc.Write(msg)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to write message")
	}

	err = wc.Close()
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to close write connection")
	}

	err = client.Quit()
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed to close SMTP connection")
	}

	err = authDao.Code.SetCodeWithEmail(ctx, req.Email, code)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &authV1.RegisterSendEmailResponse{}, nil
}

// RegisterWithEmail TODO: SMTP
func (s *S) RegisterWithEmail(ctx context.Context, req *authV1.RegisterWithEmailRequest) (*authV1.RegisterWithEmailResponse, error) {
	settings, _, _, err := systemDao.Setting.GetSettings(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")
	}

	if !settings.RegisterEnable {
		return nil, status.Error(codes.PermissionDenied, "Register is not allowed")
	}

	if settings.SmtpEnable {
		//	check email
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

	_, err = dao.User.NewUserWithBind(ctx,
		&model.User{
			Name:     req.Username,
			IsAdmin:  false,
			Salt:     salt,
			Password: crypto.Md5CryptoWithSalt(password, salt),
		},
		&model.Bind{
			OpenID:   email,
			Platform: bindType.Email,
		},
	)

	if err != nil {
		if errors.Is(err, userDao.ErrBindInfoAlreadyUsed) {
			return nil, status.Error(codes.InvalidArgument, "Email already used")
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

	bind := &model.Bind{
		OpenID: email,
	}

	if result := dao.Bind.DB().Where("open_id = ?", email).First(bind); result.Error != nil {
		return nil, status.Error(codes.InvalidArgument, "Email error")
	}

	user := &model.User{}

	if result := dao.User.DB().Where("id = ?", bind.UID).First(user); result.Error != nil {
		return nil, status.Error(codes.InvalidArgument, "Email error")
	}

	if crypto.Md5CryptoWithSalt(password, user.Salt) != user.Password {
		return nil, status.Error(codes.InvalidArgument, "Password error")
	}

	refreshToken, err := auth.GenToken(auth.Info{
		UID:            bind.UID,
		IsRefreshToken: true,
	}, auth.RefreshTokenExpireIn)
	if err != nil {
		return nil, status.Error(codes.Internal, "Internal error")

	}

	accessToken, err := auth.GenToken(auth.Info{
		UID:            bind.UID,
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

	accessToken, err := auth.GenToken(
		auth.Info{
			// following the same UID
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
