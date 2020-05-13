package schemas

import (
	"context"

	"github.com/yankooo/go-chassis/examples/schemas/helloworld"
)

//HelloServer is a struct
type HelloServer struct {
}

//SayHello is a method used to reply message
func (s *HelloServer) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "Go Hello  " + in.Name}, nil
}
