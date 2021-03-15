package api

import (
	"net/http"
)

func Authorize(res http.ResponseWriter, req *http.Request) {

	addr := req.FormValue("__wal")
	key := req.FormValue("__key")
	authRequest := &AuthorizeRequest{
		addr,
		key,
	}

	JSONResponse(res, 200, authRequest)

}

type AuthorizeRequest struct {
	WalletAddress string
	AddressKey    string
}

type AuthorizeResponse struct {
	WalletInfo string `json:"wallet_info"`
	Error      string `json:"auth_error"`
	Code       int    `json:"auth_code"`
}
