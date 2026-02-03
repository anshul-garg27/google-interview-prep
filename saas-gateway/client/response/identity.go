package response

import (
	"init.bulbul.tv/bulbul-backend/saas-gateway/client/entry"
)

type UserAuthResponse struct {
	Status            entry.Status             `json:"status,omitempty"`
	Token             string                   `json:"token,omitempty"`
	RefreshTime       int64                    `json:"refreshTime,omitempty"`
	RefreshToken      *string                  `json:"refreshToken,omitempty"`
	IsAnonymous       bool                     `json:"isAnonymous,omitempty"`
	SocialAccounts    []entry.SocialAccount    `json:"socialAccounts,omitempty"`
	Client            *entry.Client            `json:"client,omitempty"`
	UserClientAccount *entry.UserClientAccount `json:"userClientAccount,omitempty"`
	BbDeviceId        int                      `json:"bbDeviceId,omitempty"`
	IsSignup          bool                     `json:"isSignup,omitempty"`
}

type ClientAuthResponse struct {
	Status entry.Status `json:"status,omitempty"`
	Id     string       `json:"id,omitempty"`
	Token  string       `json:"token,omitempty"`
	Client entry.Client `json:"client,omitempty"`
}

type UserResponse struct {
	Status entry.Status `json:"status,omitempty"`
	Users  []entry.User `json:"users,omitempty"`
}

type PayoutResponse struct {
	Status  bool    `json:"status"`
	Amount  int64   `json:"amount_paid"`
	Message Message `json:"extra"`
}

type Message struct {
	Status StatusMessage `json:"status"`
}

type StatusMessage struct {
	Message string `json:"message"`
}

type UserClientAccountResponse struct {
	Status   entry.Status              `json:"status,omitempty"`
	Accounts []entry.UserClientAccount `json:"accounts,omitempty"`
}

type UserAccountListResponse struct {
	Status   entry.Status              `json:"status,omitempty"`
	Accounts []entry.UserClientAccount `json:"list,omitempty"`
}

type UserDevicePreferenceResponse struct {
	Status entry.Status `json:"status,omitempty"`
}

type UserSSOResponse struct {
	Status entry.Status `json:"status,omitempty"`
	Token  string       `json:"token,omitempty"`
}
