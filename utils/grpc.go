package utils

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
	"strconv"
)

func GetGrpcAuthHeader(ctx context.Context, header string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no headers in request")
	}

	authHeaders, ok := md[header]
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no header in request")
	}

	if len(authHeaders) != 1 {
		return "", status.Error(codes.Unauthenticated, "more than 1 header in request")
	}

	return authHeaders[0], nil
}
func GetRpcCtxUserid(ctx context.Context) int64 {
	useridStr, err := GetGrpcAuthHeader(ctx, "userid")
	if err != nil {
		log.Println("getCtxUserid err :", err)
		return 0
	}
	userid, _ := strconv.Atoi(useridStr)
	return int64(userid)
}
func SetGrpcHeader(ctx context.Context, header, value string) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	//mdCopy := md.Copy()
	md[header] = []string{value}
	return metadata.NewIncomingContext(ctx, md)
}
