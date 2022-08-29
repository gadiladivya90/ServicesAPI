package entity

type ServiceVersion struct {
	Name      string `db:"name"`
	Id        string `db:"service_version_id"`
	ServiceId string `db:"service_package_id"`
}

type ServiceVersionTmp struct {
	Name string
}

func (s ServiceVersion) ToSVDto() ServiceVersionResponseDto {

	return ServiceVersionResponseDto{
		Id:        s.Id,
		Name:      s.Name,
		ServiceId: s.ServiceId,
	}
}

// Convert ServiceVersionRequestDto to ServiceVersion
func NewServiceVersion(d ServiceVersionRequestDto) *ServiceVersion {
	return &ServiceVersion{
		Name:      d.Name,
		ServiceId: d.ServiceId,
	}
}
