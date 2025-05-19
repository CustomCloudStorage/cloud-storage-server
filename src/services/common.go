package services

type ServiceConfig struct {
	StorageDir     string
	Temp           string
	Secret         string
	Host           string
	TotalStorageGB int
	UserAllocGB    int
}

func (c ServiceConfig) TotalStorageBytes() int64 {
	return int64(c.TotalStorageGB) * 1024 * 1024 * 1024
}

func (c ServiceConfig) UserAllocBytes() int64 {
	return int64(c.UserAllocGB) * 1024 * 1024 * 1024
}
