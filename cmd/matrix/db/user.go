package db

type User struct {
	Id       string `xorm:"not null comment('ID') VARCHAR(32)" json:"id"`
	Username string `xorm:"not null comment('用户名') VARCHAR(128)" json:"username"`
	Password string `xorm:"not null comment('密码') VARCHAR(128)" json:"password"`
	Salt     string `xorm:"not null comment('加密盐') VARCHAR(128)" json:"salt"`
	CreateAt int64  `xorm:"not null comment('创建时间') BIGINT(20)" json:"create_at"`
}
