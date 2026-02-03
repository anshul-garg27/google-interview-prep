package response

import (
	"init.bulbul.tv/bulbul-backend/event-grpc/client/entry"
)

type UserAuthResponse struct {
	Status            entry.Status             `json:"status,omitempty"`
	Token             string                   `json:"token,omitempty"`
	RefreshTime       int64                    `json:"refreshTime,omitempty"`
	RefreshToken      *string                  `json:"refreshToken,omitempty"`
	IsAnonymous       bool                     `json:"isAnonymous,omitempty"`
	Client            *entry.Client            `json:"client,omitempty"`
	UserClientAccount *entry.UserClientAccount `json:"userClientAccount,omitempty"`
	BbDeviceId        int                      `json:"bbDeviceId,omitempty"`
	Preferences       []entry.Preference       `json:"preferences,omitempty"`
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

type UserClientAccountResponse struct {
	Status   entry.Status              `json:"status,omitempty"`
	Accounts []entry.UserClientAccount `json:"list,omitempty"`
}

type UserDevicePreferenceResponse struct {
	Status      entry.Status       `json:"status,omitempty"`
	Preferences []entry.Preference `json:"preferences,omitempty"`
}

type OTPResponse struct {
	Status    entry.Status              `json:"status,omitempty"`
	RequestID int                       `json:"requestId,omitempty"`
	Operation entry.OTPRequestOperation `json:"operation,omitempty"`
}

type TrueCallerResponse struct {
	Status  entry.Status            `json:"status,omitempty"`
	Request entry.TrueCallerRequest `json:"request,omitempty"`
	Success bool                    `json:"success,omitempty"`
}
