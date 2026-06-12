package booking

import (
    "context"
    "fmt"
    "time"
    "github.com/redis/go-redis/v9"
)

const lockTTL = 5 * time.Minute

func lockKey(showtimeID, seat string) string {
    return fmt.Sprintf("lock:seat:%s:%s", showtimeID, seat)
}

func AcquireLock(ctx context.Context, rdb *redis.Client, showtimeID, seat, userID string) (bool, error) {
    key := lockKey(showtimeID, seat)

    ok, err := rdb.SetNX(ctx, key, userID, lockTTL).Result()
    if err != nil {
        return false, fmt.Errorf("redis error: %w", err)
    }
    return ok, nil
}

func ReleaseLock(ctx context.Context, rdb *redis.Client, showtimeID, seat, userID string) error {
    key := lockKey(showtimeID, seat)

    script := redis.NewScript(`
        if redis.call("GET", KEYS[1]) == ARGV[1] then
            return redis.call("DEL", KEYS[1])
        else
            return 0
        end
    `)
    return script.Run(ctx, rdb, []string{key}, userID).Err()
}

func GetLockOwner(ctx context.Context, rdb *redis.Client, showtimeID, seat string) (string, error) {
    return rdb.Get(ctx, lockKey(showtimeID, seat)).Result()
}