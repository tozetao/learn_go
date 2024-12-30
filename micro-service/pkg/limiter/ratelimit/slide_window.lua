local key = KEYS[1]

-- 窗口大小
local window = tonumber(ARGV[1])
-- 速率
local rate = tonumber(ARGV[2])
-- 当前毫秒时间戳
local now = tonumber(ARGV[3])

-- 窗口起始时间
local min = now - window

-- 删除窗口中过期的元素
redis.call('ZREMRANGEBYSCORE', key, '-inf', min)

local cnt = redis.call('ZCOUNT', key, '-inf', '+inf')
if (cnt >= rate) then
    return "true"
else
    -- 把score和member设置成now
    redis.call('ZADD', key, now, now)
    redis.call('PEXPIRE', key, window)
    return "false"
end
