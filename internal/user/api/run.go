package api

import (
	"im/config"
	"im/pkg/pb"
	"im/pkg/util"
	"net"

	"google.golang.org/grpc"
)

// StartRpcServer 启动rpc服务
func StartRpcServer() {
	go func() {
		defer util.RecoverPanic()

		intListen, err := net.Listen("tcp", config.UserConf.RPCIntListenAddr)
		if err != nil {
			panic(err)
		}
		intServer := grpc.NewServer(grpc.UnaryInterceptor(UserIntInterceptor))
		pb.RegisterUserIntServer(intServer, &UserIntServer{})
		err = intServer.Serve(intListen)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		defer util.RecoverPanic()

		extListen, err := net.Listen("tcp", config.UserConf.RPCExtListenAddr)
		if err != nil {
			panic(err)
		}
		extServer := grpc.NewServer(grpc.UnaryInterceptor(UserExtInterceptor))
		pb.RegisterUserExtServer(extServer, &UserExtServer{})
		err = extServer.Serve(extListen)
		if err != nil {
			panic(err)
		}
	}()

}
