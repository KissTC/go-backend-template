package models

// en modelos le decimos a nivel de Go como funciona el flujo de datos
// como funcionan los modelos, que tipos usan, como se integran

// modelo de usuario a nivel de Go (se tiene que ver a nivel de go, patron de diseño, bd, etc)
type User struct {
	Id       string `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
