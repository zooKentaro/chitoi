package redis

import (
    "fmt"
    "os"

    "github.com/garyburd/redigo/redis"
)

func Connect() (redis.Conn, error) {
    return redis.Dial(
        "tcp",
        fmt.Sprintf(
            "%s:%s",
            os.Getenv("CHITOI_REDIS_HOST"),
            os.Getenv("CHITOI_REDIS_PORT"),
        ),
    )
}
