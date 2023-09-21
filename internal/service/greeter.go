package service

import (
	"context"

	"github.com/omalloc/contrib/protobuf"
	"github.com/samber/lo"

	v1 "github.com/omalloc/kratos-layout/api/helloworld/v1"
	"github.com/omalloc/kratos-layout/internal/biz"
)

// GreeterService is a greeter service.
type GreeterService struct {
	v1.UnimplementedGreeterServer

	uc *biz.GreeterUsecase
}

// NewGreeterService new a greeter service.
func NewGreeterService(uc *biz.GreeterUsecase) *GreeterService {
	return &GreeterService{uc: uc}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, in *v1.HelloRequest) (*v1.HelloReply, error) {
	g := &biz.Greeter{Hello: in.Name}
	if err := s.uc.CreateGreeter(ctx, g); err != nil {
		return nil, err
	}
	return &v1.HelloReply{Message: "Hello " + g.Hello}, nil
}

func (s *GreeterService) List(ctx context.Context, in *v1.HeelloListRequest) (*v1.HelloListReply, error) {
	pagination := protobuf.PageWrap(in.Pagination)
	data, err := s.uc.List(ctx, pagination, in.Name)
	if err != nil {
		return nil, err
	}

	return &v1.HelloListReply{
		Pagination: pagination.Resp(),
		Data: lo.Map(data, func(item *biz.Greeter, _ int) *v1.GreeterInfo {
			return &v1.GreeterInfo{
				Id:        item.ID,
				Hello:     item.Hello,
				CreatedAt: item.CreatedAt.Unix(),
				UpdatedAt: item.UpdatedAt.Unix(),
			}
		}),
	}, nil
}
