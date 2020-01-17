package cache

import (
	"im/pkg/db"
	"im/pkg/gerrors"
	"strconv"
)

const (
	UserSeqKey  = "user_seq:"
	GroupSeqKey = "group_seq"
)

type seqCache struct{}

var SeqCache = new(seqCache)

func (*seqCache) UserKey(userId int64) string {
	return UserSeqKey + strconv.FormatInt(userId, 10)
}

func (*seqCache) GroupKey(groupId int64) string {
	return GroupSeqKey + strconv.FormatInt(groupId, 10)
}

// Incr 将序列号增加1
func (c *seqCache) Incr(key string) (int64, error) {
	seq, err := db.RedisCli.Incr(key).Result()
	if err != nil {
		return 0, gerrors.WrapError(err)
	}
	return seq, nil
}
