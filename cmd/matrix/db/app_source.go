package db

type AppSource struct {
	Id       string `xorm:"not null comment('ID') VARCHAR(34)"`
	AppId    string `xorm:"not null comment('程序ID') VARCHAR(34)"`
	SourceId string `xorm:"not null comment('数据源ID') VARCHAR(32)"`
}
