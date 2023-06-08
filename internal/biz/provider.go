package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	v1 "kratos-base-template/api/provider/v1"
)

var (
	ErrEventNotSupport = errors.NotFound(v1.ErrorReason_EVENT_NOT_SUPPORT.String(), "event not support")
)

// Greeter is a Greeter model.
type User struct {
	Hello string
}

// ProviderBaseInterface is a Provider repo.
type ProviderBaseInterface interface {
	Save(context.Context, *User) (*User, error)
	Update(context.Context, *User) (*User, error)
	Get(context.Context, int64) (*User, error)
	ListUser(context.Context) ([]*User, error)
}

// ProviderUseCase is a  Provider UseCase.
type ProviderUseCase struct {
	pro ProviderBaseInterface
	log *log.Helper
}

// NewProviderUseCase new a Provider UseCase.
func NewProviderUseCase(pro ProviderBaseInterface, logger log.Logger) *ProviderUseCase {
	return &ProviderUseCase{pro: pro, log: log.NewHelper(logger)}
}

// CreateProvider creates a Provider, and returns the new Provider.
func (uc *ProviderUseCase) CreateProvider(ctx context.Context, req *v1.Request) (*v1.Response, error) {
	//uc.log.WithContext(ctx).Infof("CreateProvider: %#v", req)
	switch req.GetEvent() {
	case v1.EVENT_CODE_ADD_USER:
		return uc.save(ctx, req)
	case v1.EVENT_CODE_UPDATE_USER:
		return uc.update(ctx, req)
	case v1.EVENT_CODE_GET_USER:
		return uc.get(ctx, req)
	case v1.EVENT_CODE_LIST_USER:
		return uc.update(ctx, req)
	default:
		return &v1.Response{
				Event: req.GetEvent(),
				ReqId: req.GetReqId(),
				Code:  v1.ErrorReason_EVENT_NOT_SUPPORT,
				Desc:  "Not Supported Event",
			},
			nil
	}
}

func (uc *ProviderUseCase) save(ctx context.Context, req *v1.Request) (*v1.Response, error) {
	//uc.log.WithContext(ctx).Infof("Save: %#v", req)
	uc.pro.Save(ctx, &User{})
	return &v1.Response{
		ReqId: req.GetReqId(),
		Event: req.GetEvent(),
		Code:  v1.ErrorReason_SUCCESS,
		Desc:  "success",
	}, nil
}

func (uc *ProviderUseCase) update(ctx context.Context, req *v1.Request) (*v1.Response, error) {
	uc.log.WithContext(ctx).Infof("update: %#v", req)
	uc.pro.Update(ctx, &User{})
	return nil, nil
}

func (uc *ProviderUseCase) get(ctx context.Context, req *v1.Request) (*v1.Response, error) {
	uc.log.WithContext(ctx).Infof("Get: %#v", req)
	uc.pro.Get(ctx, 0)
	return nil, nil
}

func (uc *ProviderUseCase) listUser(ctx context.Context, req *v1.Request) (*v1.Response, error) {
	uc.log.WithContext(ctx).Infof("ListUser: %#v", req)
	uc.pro.ListUser(ctx)
	return nil, nil
}
