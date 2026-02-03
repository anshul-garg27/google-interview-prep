package entry

type PreferenceKey string

const (
	NOTIFICATION  PreferenceKey = "NOTIFICATION"
	LOCALE        PreferenceKey = "LOCALE"
	WHATSAPPOPTIN PreferenceKey = "WHATSAPPOPTIN"
)

type Preference struct {
	Id       int           `json:"id,omitempty"`
	Key      PreferenceKey `json:"key,omitempty"`
	Value    string        `json:"value,omitempty"`
	UserId   *string       `json:"userId,omitempty"`
	DeviceId string        `json:"deviceId,omitempty"`
	ClientId string        `json:"clientId,omitempty"`
}
