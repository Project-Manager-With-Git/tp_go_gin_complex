package setting

import "time"

var RedisQueryTimeout = time.Duration(50) * time.Millisecond

func Init(redis_query_timeout_ms int) {
	RedisQueryTimeout = time.Duration(redis_query_timeout_ms) * time.Millisecond
}
