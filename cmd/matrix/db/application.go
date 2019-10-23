package db

type Application struct {
	Id          string `xorm:"not null pk comment('唯一ID') VARCHAR(34)" json:"id"`
	Name        string `xorm:"not null comment('应用名称') VARCHAR(128)" json:"name"`
	Type        string `xorm:"not null comment('应用类型') VARCHAR(128)" json:"type"`
	Description string `xorm:"not null comment('详情描述') VARCHAR(1024)" json:"description"`
	Section     string `xorm:"not null comment('数据维护部门') VARCHAR(128)" json:"section"`
	Sources     string `xorm:"not null comment('数据来源项目IDs') TEXT" json:"sources"`
	Owner       string `xorm:"not null comment('所属者') VARCHAR(32)" json:"owner"`
	SecretKey   string `xorm:"not null comment('secret key') VARCHAR(128)" json:"secret_key"`
}
