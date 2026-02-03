package input

import "init.bulbul.tv/bulbul-backend/saas-gateway/client/entry"

type Token struct {
	Token        string  `json:"token,omitempty"`
	RefreshToken *string `json:"refreshToken,omitempty"`
}

type UserAuth struct {
	Mode                    AuthMode                 `json:"mode,omitempty"`
	Id                      string                   `json:"id,omitempty"`
	Secret                  string                   `json:"secret,omitempty"`
	ClientAccount           *entry.UserClientAccount `json:"clientAccount,omitempty"`
	TrueCallerSignatureAlgo *string                  `json:"trueCallerSignatureAlgo,omitempty"`
	CountryCode             *string                  `json:"countryCode,omitempty"`
	WhatsappOptin           *bool                    `json:"whatsappOptin"`
	Password                *string                  `json:"password"`
}

type PhoneVerifyInput struct {
	Id            string `json:"id"`
	Otp           string `json:"otp"`
	WhatsappOptIn *bool  `json:"whatsappOptin"`
}

type InvoiceInput struct {
	InvoiceId string `json:"invoice_id"`
}

type ClientAuth struct {
	Id     string `json:"id,omitempty"`
	Secret string `json:"secret,omitempty"`
}

type AuthMode string

const (
	AuthModeAnonymous     AuthMode = "ANONYMOUS"
	AuthModeEmail         AuthMode = "EMAIL"
	AuthModeEmailOTP      AuthMode = "EMAIL_OTP"
	AuthModePhone         AuthMode = "PHONE"
	AuthModeGoogle        AuthMode = "GOOGLE"
	AuthModeGoogleOauth   AuthMode = "GOOGLE_OAUTH"
	AuthModeYoutbe        AuthMode = "YOUTBE"
	AuthModeFb            AuthMode = "FB"
	AuthModeInstagram     AuthMode = "INSTAGRAM"
	AuthModeTruecaller    AuthMode = "TRUECALLER"
	AuthModeTruecallerWeb AuthMode = "TRUECALLER_WEB"
	AuthModeMx            AuthMode = "MX"
	AuthModePaytm         AuthMode = "PAYTM"
	AuthModeBulBulSSO     AuthMode = "BULBUL_SSO"
)

type OTPRequestInput struct {
	Operation     OTPRequestOperation `json:"operation,omitempty"`
	TargetMode    OTPRequestTarget    `json:"targetMode,omitempty"`
	Identity      string              `json:"identity"`
	CountryCode   string              `json:"countryCode"`
	SocialAccount *SocialAccountInput `json:"socialAccount"`
}

type SocialAccountInput struct {
	Platform entry.SocialNetworkType `json:"platform"`
	Handle   string                  `json:"handle"`
}

type OTPRequestOperation string

const (
	OTPRequestOperationLink   OTPRequestOperation = "LINK"
	OTPRequestOperationVerify OTPRequestOperation = "VERIFY"
	OTPRequestOperationLogin  OTPRequestOperation = "LOGIN"
	OTPRequestOperationSignup OTPRequestOperation = "SIGNUP"
)

type OTPRequestTarget string

const (
	OTPRequestTargetPhone OTPRequestTarget = "PHONE"
	OTPRequestTargetEmail OTPRequestTarget = "EMAIL"
)
