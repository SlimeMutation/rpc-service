package client

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
)

var errWalletHttpError = errors.New("wallet http error")

type Address struct {
	PublicKey string
	Address   string
}

type WalletClient interface {
	GetSupportCoins(chain, network string) (bool, error)
	GetWalletAddress(chain, network string) (*Address, error)
}

type Client struct {
	client *resty.Client
}

// GetSupportCoins implements WalletClient.
func (c *Client) GetSupportCoins(chain string, network string) (bool, error) {
	res, err := c.client.R().SetQueryParams(map[string]string{
		"chain":   chain,
		"network": network,
	}).SetResult(&SupportChainResponse{}).Get("/api/v1/support_chain")
	if err != nil {
		return false, errors.New("failed to request GetSupportCoins")
	}
	spt, ok := res.Result().(*SupportChainResponse)
	if !ok {
		return false, errors.New("failed to parse GetSupportCoins response")
	}
	return spt.Support, nil
}

// GetWalletAddress implements WalletClient.
func (c *Client) GetWalletAddress(chain string, network string) (*Address, error) {
	res, err := c.client.R().SetQueryParams(map[string]string{
		"chain":   chain,
		"network": network,
	}).SetResult(&WalletAddressResponse{}).Get("/api/v1/wallet_address")
	if err != nil {
		return nil, errors.New("failed to request GetWalletAddress")
	}
	addr, ok := res.Result().(*WalletAddressResponse)
	if !ok {
		return nil, errors.New("failed to parse GetWalletAddress response")
	}
	return &Address{
		PublicKey: addr.PublicKey,
		Address:   addr.Address,
	}, nil
}

func NewWalletClient(url string) WalletClient {
	client := resty.New()
	client.SetBaseURL(url)
	client.OnAfterResponse(func(c *resty.Client, r *resty.Response) error {
		statusCode := r.StatusCode()
		if statusCode >= 400 {
			method := r.Request.Method
			baseUrl := r.Request.URL
			return fmt.Errorf("%s %s failed with status code %d: %w", method, baseUrl, statusCode, errWalletHttpError)
		}
		fmt.Println("method: ", r.Request.Method)
		fmt.Println("baseUrl: ", r.Request.URL)
		fmt.Println("params: ", r.Request.QueryParam)
		return nil
	})
	return &Client{
		client: client,
	}
}
