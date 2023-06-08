package service

import (
	"context"
	pb "kratos-base-template/api/provider/v1"
	"kratos-base-template/internal/biz"
)

type ProviderService struct {
	pb.UnimplementedProviderServer
	proUseCase *biz.ProviderUseCase
}

func NewProviderService(proUseCase *biz.ProviderUseCase) *ProviderService {
	return &ProviderService{proUseCase: proUseCase}
}

func (s *ProviderService) BaserService(ctx context.Context, req *pb.Request) (*pb.Response, error) {
	return s.proUseCase.CreateProvider(ctx, req)
}
