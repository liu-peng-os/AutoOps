package websocket

import (
	"dodevops-api/common/util"

	"gorm.io/gorm"
)

type sshRecord struct {
	gorm.Model
	ConnectID   string     `gorm:"comment:'connect id';size:64" json:"connect_id"`
	UserName    string     `gorm:"comment:'system user';size:128" json:"user_name"`
	HostName    string     `gorm:"comment:'host name';size:128" json:"host_name"`
	ConnectTime util.HTime `gorm:"index;comment:'connect time'" json:"connect_time"`
	LogoutTime  util.HTime `gorm:"index;comment:'logout time'" json:"logout_time"`
	Records     []byte     `json:"records" gorm:"type:longblob;comment:'terminal recording';size:128"`
	HostId      uint       `gorm:"comment:'host id'" json:"host_id"`
}

func (s sshRecord) TableName() string {
	return "ssh_record"
}
