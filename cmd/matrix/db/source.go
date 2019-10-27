package db

type Source struct {
	Id          string `xorm:"not null comment('唯一ID') VARCHAR(32)" json:"id"`
	Name        string `xorm:"not null comment('名称') VARCHAR(128)" json:"name"`
	Description string `xorm:"not null comment('详情描述') VARCHAR(1024)" json:"description"`
	Section     string `xorm:"not null comment('数据维护部门') VARCHAR(128)" json:"section"`
	Owner       string `xorm:"not null comment('所属者') VARCHAR(32)" json:"owner"`
	Service     string `xorm:"not null comment('服务URI地址') TEXT" json:"service"`
	Search      string `xorm:"not null comment('查询服务') TEXT" json:"search"`
	ApiDoc      string `xorm:"not null comment('API开发文档URI') TEXT" json:"api_doc"`
	Type        string `xorm:"not null comment('连通方式') VARCHAR(128)" json:"type"`
	Secret      string `xorm:"not null comment('加密字符串') VARCHAR(128)" json:"secret"`
	CreateAt    int    `xorm:"default NULL comment('创建时间') INT(11)" json:"create_at"`
}
type SourceEx struct {
	Id          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Section     string `json:"section"`
	Owner       string `json:"owner"`
	Service     string `json:"service"`
	Search      string `json:"search"`
	ApiDoc      string `json:"api_doc"`
	Type        string `json:"type"`
	Secret      string `json:"secret"`
	CreateAt    int    `json:"create_at"`
	Check       int    `json:"check"`
}
type SourceLabel struct {
	Id     string `json:"id"`
	Name   string `json:"name"`
	PinYin string `json:"pin_yin"`
}
