package main

// BasicUser базовые поля
type BasicUser struct {
	Username  string `json:"username"`
	PhotoUUID string `json:"photo_uuid"`
}

// InfoUser BasicUser, расширенный служебной инфой
type InfoUser struct {
	BasicUser
	ID     int64 `json:"id"`
	Active bool  `json:"active"`
}

// FormUser BasicUser, расширенный паролем, используется для входа и регистрации
type FormUser struct {
	BasicUser
	Password string `json:"password"`
}
