package middleware

import (
	"context"
	grpcAuth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/juanjiTech/jframe/core/logx"
	"github.com/juanjiTech/jframe/pkg/auth"
	"github.com/juanjiTech/jframe/pkg/ctxKey"
)

func AuthInterceptor(ctx context.Context) (context.Context, error) {
	token, err := grpcAuth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return ctx, nil //不强求auth 交给具体方法判断（后续使用CASBin控制
	}

	jwtClaims, err := auth.ParseToken(token)
	if err != nil {
		logx.NameSpace("grpc.middleware.auth").Infof("Error parsing token: %v", err)
		return ctx, nil
	}
	newCtx := context.WithValue(ctx, ctxKey.UID, jwtClaims.Info.UID)
	newCtx = context.WithValue(newCtx, ctxKey.OrgID, jwtClaims.Info.OrgID)
	return newCtx, nil
}
