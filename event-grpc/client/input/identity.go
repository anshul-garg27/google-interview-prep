package input

import "init.bulbul.tv/bulbul-backend/event-grpc/client/entry"

type Token struct {
	Token        string  `json:"token,omitempty"`
	RefreshToken *string `json:"refreshToken,omitempty"`
}

type UserAuth struct {
	Mode                    AuthMode                 `json:"mode,omitempty"`
	Id                      string                   `json:"id,omitempty"`
	Secret                  string                   `json:"secret,omitempty"`
	ClientAccount           *UpdateUserClientAccount `json:"clientAccount,omitempty"`
	TrueCallerSignatureAlgo *string                  `json:"trueCallerSignatureAlgo,omitempty"`
	CountryCode             *string                  `json:"countryCode,omitempty"`
}

type ClientAuth struct {
	Id     string `json:"id,omitempty"`
	Secret string `json:"secret,omitempty"`
}

type UpdateUser struct {
	Status *string       `json:"status,omitempty"`
	Dob    *int64        `json:"dob,omitempty"`
	Gender *entry.Gender `json:"gender,omitempty"`
}

type UpdateUserClientAccount struct {
	Name *string `json:"name,omitempty"`
	Bio  *string `json:"bio,omitempty"`

	ProfileImageId *int `json:"profileImageId,omitempty"`
	CoverImageId   *int `json:"coverImageId,omitempty"`

	ReferredViaCode     *string `json:"referredViaCode,omitempty"`
	ReferredViaClientId *string `json:"referredViaClientId,omitempty"`

	Permissions []string `json:"permissions,omitempty"`

	User *UpdateUser `json:"user,omitempty"`

	//
	//ClientId               string                `json:"clientId,omitempty"`
	//ClientAppType          string                `json:"clientAppType,omitempty"`
	//UserId                 int                   `json:"userId,omitempty"`
	//Status                 UserClientAccountStatus                `json:"status,omitempty"`
	//OnboardingStep         UserClientAccountOnboardingStep                `json:"onboardingStep,omitempty"`
	//Remark *string               `json:"remark,omitempty"`
	//
	//
	//HostProfile *HostProfile `json:"hostProfile,omitempty"`
	//CustomerProfile *CustomerProfile `json:"customerProfile,omitempty"`
}

type AuthMode string

const (
	AuthModeAnonymous     AuthMode = "ANONYMOUS"
	AuthModeEmail         AuthMode = "EMAIL"
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
)

type UserStatus string

const (
	UserStatusTempBlocked UserStatus = "TEMP_BLOCKED"
	UserStatusActive      UserStatus = "ACTIVE"
	UserStatusInactive    UserStatus = "INACTIVE"
	UserStatusBlocked     UserStatus = "BLOCKED"
	UserStatusSuspended   UserStatus = "SUSPENDED"
	UserStatusBlacklisted UserStatus = "BLACKLISTED"
)

type OTPRequestInput struct {
	Operation   OTPRequestOperation `json:"operation,omitempty"`
	TargetMode  OTPRequestTarget    `json:"targetMode,omitempty"`
	Identity    string              `json:"identity"`
	CountryCode string              `json:"countryCode"`
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
