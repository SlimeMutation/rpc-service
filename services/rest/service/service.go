package service

import (
	"github.com/SlimeMutation/rpc-service/database"
	"github.com/SlimeMutation/rpc-service/services/rest/model"
)

type Service interface {
	GetSupportCoins(*model.ChainRequest) (*model.SupportChainResponse, error)
	GetWalletAddress(*model.ChainRequest) (*model.WalletAddressResponse, error)
}

type HandleSrv struct {
	v        *Validator
	keysView database.KeysView
}

func (h HandleSrv) GetSupportCoins(req *model.ChainRequest) (*model.SupportChainResponse, error) {
	ok := h.v.VerifyWalletAddress(req.Chain, req.Network)
	return &model.SupportChainResponse{
		Support: ok,
	}, nil
}

func (h HandleSrv) GetWalletAddress(*model.ChainRequest) (*model.WalletAddressResponse, error) {
	return &model.WalletAddressResponse{
		PublicKey: "public key",
		Address:   "0x00000000",
	}, nil
}

func NewHandleSrv(v *Validator, keysView database.KeysView) Service {
	return &HandleSrv{
		v:        v,
		keysView: keysView,
	}
}
