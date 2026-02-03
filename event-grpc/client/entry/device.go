package entry

type DeviceAttribute struct {
	Id    int         `json:"id,omitempty"`
	Key   string      `json:"key,omitempty"`
	Value interface{} `json:"value,omitempty"`
}
