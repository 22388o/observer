package cmd

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"log"
	"obwallet/obrpc/user"
	"os"
)

var gClient user.UserServiceClient

func init() {
	//certData, _, err := cert.LoadCert(
	//	"/home/wxf/git/lnd/docker/obtest/volumes/lnd/bob/tls.cert", "/home/wxf/git/lnd/docker/obtest/volumes/lnd/bob/tls.key",
	//)
	//tlsCfg := &tls.Config{Certificates: []tls.Certificate{certData}}
	//tlsCfg.InsecureSkipVerify = true
	//credentials.NewTLS(tlsCfg)

	conn, err := grpc.DialContext(context.TODO(), "localhost:19090",
		grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{Certificates: []tls.Certificate{}, InsecureSkipVerify: true})),
	)
	if err != nil {
		log.Fatal("error grpc conn")
	}
	gClient = user.NewUserServiceClient(conn)
}

func printJSON(resp interface{}) {
	b, err := json.Marshal(resp)
	if err != nil {
		log.Fatalln(err)
	}

	var out bytes.Buffer
	json.Indent(&out, b, "", "\t")
	out.WriteString("\n")
	out.WriteTo(os.Stdout)
}

func printRespJSON(resp proto.Message) {
	jsonMarshaler := &jsonpb.Marshaler{
		EmitDefaults: true,
		OrigName:     true,
		Indent:       "    ",
	}

	jsonStr, err := jsonMarshaler.MarshalToString(resp)
	if err != nil {
		fmt.Println("unable to decode response: ", err)
		return
	}

	fmt.Println(jsonStr)
}
