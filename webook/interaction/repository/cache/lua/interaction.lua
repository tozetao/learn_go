
-- key
local key = KEYS[1]

-- 要自增的hash的field
local field = ARGV[1]
local value = tonumber(ARGV[2])

if redis.call("exists", key) == 1 then
    -- 返回自增后的值
    return redis.call("hincrby", key, field, value)
else
    return nil
end