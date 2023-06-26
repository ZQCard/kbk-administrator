package domain

type Administrator struct {
	Id            int64
	Username      string
	Password      string
	Salt          string
	Mobile        string
	Nickname      string
	Avatar        string
	Status        bool
	Role          string
	LastLoginTime string
	LastLoginIp   string
	CreatedAt     string
	UpdatedAt     string
	DeletedAt     string
}
