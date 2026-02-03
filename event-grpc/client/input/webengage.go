package input

type WEBENGAGE_APP_TYPE string

const (
	WEBENGAGE_APP_TYPE_USER  WEBENGAGE_APP_TYPE = "USER"
	WEBENGAGE_APP_TYPE_BRAND WEBENGAGE_APP_TYPE = "BRAND"
)

type WebengageEvent struct {
	AppType     *WEBENGAGE_APP_TYPE    `json:"appType,omitempty"`
	UserId      *string                `json:"userId,omitempty"`
	AnonymousId *string                `json:"anonymousId,omitempty"`
	EventName   string                 `json:"eventName,omitempty"`
	EventTime   string                 `json:"eventTime,omitempty"`
	EventData   map[string]interface{} `json:"eventData,omitempty"`
}

type WebengageUserEvent struct {
	AppType       *WEBENGAGE_APP_TYPE    `json:"appType,omitempty"`
	UserId        *string                `json:"userId,omitempty"`
	Phone         *string                `json:"phone,omitempty"`
	Email         *string                `json:"email,omitempty"`
	FirstName     *string                `json:"firstName,omitempty"`
	LastName      *string                `json:"lastName,omitempty"`
	WhatsappOptIn *bool                  `json:"whatsappOptIn,omitempty"`
	EmailOptIn    *bool                  `json:"emailOptIn,omitempty"`
	Attributes    map[string]interface{} `json:"attributes,omitempty"`
}
