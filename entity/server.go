package entity

type Server struct {
	ID        uint64 `json:"id" gorm:"primary_key;auto_increment"`
	Name      string `json:"name" gorm:"UNIQUE"`
	Ipv4      string `json:"ipv4" binding:"required" validate:"ipv4"`
	Status    bool   `json:"status" gorm:"default:false"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"update_at"`
}
