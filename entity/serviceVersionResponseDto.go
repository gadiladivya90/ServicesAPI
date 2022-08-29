package entity

type ServiceVersionResponseDto struct {
	Id        string `json:"service_version_id"`
	Name      string `json:"name"`
	ServiceId string `json:"service_id"`
}
