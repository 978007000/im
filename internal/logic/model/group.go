package model

import "time"

const (
	GroupTypeGroup    = 1 // 群组
	GroupTypeChatRoom = 2 // 聊天室
)

// Group 群组
type Group struct {
	Id           int64     // 群组id
	Name         string    // 组名
	Introduction string    // 群简介
	UserNum      int32     // 群组人数
	Type         int32     // 群组类型
	Extra        string    // 附加字段
	CreateTime   time.Time // 创建时间
	UpdateTime   time.Time // 更新时间
}

type GroupUserUpdate struct {
	GroupId int64  `json:"group_id"` // 群组id
	UserId  int64  `json:"user_id"`  // 用户id
	Label   string `json:"label"`    // 用户标签
	Extra   string `json:"extra"`    // 群组用户附件属性
}
