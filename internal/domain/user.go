package domain

import "time"

// User : 领域对象(业务意义上的用户对象)
type User struct {
	Id       int64
	Email    string
	Password string
	Ctime    time.Time
}
