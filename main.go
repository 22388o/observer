package main

import (
	"context"
	"crypto/tls"
	"flag"
	"github.com/lightningnetwork/lnd/cert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"net"
	"obwallet/obrpc"
	"obwallet/obrpc/user"
	"obwallet/utils"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

func main() {
	var serverPort = "38332"
	var dbConnstr = ""
	var mailServer string
	var mailUser string
	var mailPwd string
	var certDir string
	var kycUrl string
	var payUrl string
	var merNo string
	var merKey string

	flag.StringVar(&serverPort, "server_port", "19090", " service port")
	flag.StringVar(&dbConnstr, "db_conn", "root:password@tcp(127.0.0.1:3306)/obwallet?charset=utf8&parseTime=True&loc=Local", "mysql connstr : root:password@tcp(127.0.0.1:3306)/obwallet?charset=utf8&parseTime=True&loc=Local")
	flag.StringVar(&mailServer, "mailServer", "smtp.163.com:25", " ")
	flag.StringVar(&mailUser, "mailuser", "", "")
	flag.StringVar(&mailPwd, "mailpwd", "", "")
	flag.StringVar(&certDir, "cert_dir", "./cert", "tls cert dir")

	flag.StringVar(&kycUrl, "kyc_url", "http://116.204.78.80:9567", "")
	flag.StringVar(&payUrl, "pay_url", "https://payment.flyzeroc.xyz", "")
	flag.StringVar(&merNo, "mer_no", "104001001", "")
	flag.StringVar(&merKey, "mer_key", "_bndqOD^", "")

	flag.Parse()

	if dbConnstr != "" {
		utils.InitDb(dbConnstr)
	} else {
		panic("miss dbConnstr")
	}
	if mailUser == "" || mailPwd == "" {
		panic("miss mailUser or v")
	}
	utils.InitMailAuth(mailServer, mailUser, mailPwd)
	utils.InitCapcha()

	//grpcs
	lis1, err := net.Listen("tcp", ":"+serverPort)
	if err != nil {
		log.Fatalln("Failed to listen:", err)
	}
	tlsCfg := getTlsCfg(certDir)
	serverCreds := credentials.NewTLS(tlsCfg)
	serverOpts := []grpc.ServerOption{grpc.Creds(serverCreds),
		grpc.UnaryInterceptor(useridInterceptor),
	}
	gserver := grpc.NewServer(serverOpts...)
	us := obrpc.NewUserServer(kycUrl, payUrl, merNo, merKey)
	user.RegisterUserServiceServer(gserver, us)
	log.Println("Serving gRPCs on 0.0.0.0:" + serverPort)
	log.Fatalln(gserver.Serve(lis1))
}

func getTlsCfg(certDir string) *tls.Config {
	certPath := filepath.Join(certDir, "tls.cert")
	keyPath := filepath.Join(certDir, "tls.key")
	if !fileExists(certPath) {
		log.Println("Generating TLS certificates...")
		err := cert.GenCertPair(
			"lnd autogenerated cert", certPath,
			keyPath, nil, nil, true,
			1000*24*time.Hour,
		)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Done generating TLS certificates")
	}
	certData, _, err := cert.LoadCert(
		certPath, keyPath,
	)
	if err != nil {
		log.Println("err load TLS certificates")
	}
	//tlsCfg := cert.TLSConfFromCert(certData)
	return &tls.Config{Certificates: []tls.Certificate{certData}}
	//return tlsCfg
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func useridInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Println("useridInterceptor ", info.FullMethod)
	if info.FullMethod == "/user.UserService/SignIn" {
		return handler(ctx, req)
	}
	if info.FullMethod == "/user.UserService/SignUp" {
		return handler(ctx, req)
	}
	if info.FullMethod == "/user.UserService/VerifyCode" {
		return handler(ctx, req)
	}
	token, err := utils.GetGrpcAuthHeader(ctx, "token")
	if err != nil {
		return nil, err
	}
	userid, err := obrpc.MustGetTokenUserId(token)
	if err != nil {
		return nil, err
	}
	newCtx := utils.SetGrpcHeader(ctx, "userid", strconv.Itoa(int(userid)))
	return handler(newCtx, req)
}

// WrappedServerStream is a thin wrapper around grpc.ServerStream that allows modifying context.
type WrappedServerStream struct {
	grpc.ServerStream
	// WrappedContext is the wrapper's own Context. You can assign it.
	WrappedContext context.Context
}

// Context returns the wrapper's WrappedContext, overwriting the nested grpc.ServerStream.Context()
func (w *WrappedServerStream) Context() context.Context {
	return w.WrappedContext
}

// WrapServerStream returns a ServerStream that has the ability to overwrite context.
func WrapServerStream(stream grpc.ServerStream) *WrappedServerStream {
	if existing, ok := stream.(*WrappedServerStream); ok {
		return existing
	}
	return &WrappedServerStream{ServerStream: stream, WrappedContext: stream.Context()}
}
