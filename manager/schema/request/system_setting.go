package request

type GetSystemSettingByInfo struct {
	Key               string `json:"key_name"`
	Category          string `json:"category"`
	ObjectReferenceId string `json:"object_reference_id"`
}

type CreateOrUpdateSystemSetting struct {
	Key               string `json:"key_name"`
	Value             string `json:"value"`
	Category          string `json:"category"`
	ObjectReferenceId string `json:"object_reference_id"`
}
