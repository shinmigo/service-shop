package main

import (
	grpcserver "goshop/service-shop/pkg/grpc/server"
)

func initService() {
	go grpcserver.Run()
	//go user.Hello()
}
