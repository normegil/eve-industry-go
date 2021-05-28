package model

type Portraits struct {
	URL64  string `json:"url64" bson:"url64"`
	URL128 string `json:"url128" bson:"url128"`
	URL256 string `json:"url256" bson:"url256"`
	URL512 string `json:"url512" bson:"url512"`
}

type Identity struct {
	ID           int64     `json:"id" bson:"id"`
	Name         string    `json:"name" bson:"name"`
	RefreshToken string    `json:"-" bson:"refresh_token"`
	Portraits    Portraits `json:"portraits" bson:"portraits"`
	Role         string    `json:"role" bson:"role"`
}

func IdentityAnonymous() *Identity {
	return &Identity{
		ID:           -1,
		Name:         "anonymous",
		RefreshToken: "",
		Portraits:    Portraits{},
	}
}
