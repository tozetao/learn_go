local key = KEYS[1]
local cntKey = key..":cnt"

local code = ARGV[1]

-- key的有效期，单位秒
local expiration = ARGV[2]

-- 默认验证码间隔的发送时间为60s。注：过期时间必须大于间隔时间。
local limit = expiration - 60

-- code的验证次数
local cnt = 5

local ttl = redis.call("ttl", key)

if ttl == -2 or ttl < limit then
    -- key不存在 或者 code的过期时间小于9分钟，可以再次进行发送。
    redis.call("set", key, code, "EX", expiration)
    redis.call("set", cntKey, cnt, "EX", expiration)
    return 0
elseif ttl == -1 then
    -- key没有设置过期时间，属于系统异常
    return -2
else
    -- access too many times
    return -1
end