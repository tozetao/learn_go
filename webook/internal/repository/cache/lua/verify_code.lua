local key = KEYS[1]
local cntKey = key..":cnt"

local input_code = ARGV[1]


-- code的验证次数是否用完
local cnt = tonumber(redis.call("get", cntKey))

if cnt == nil or cnt <= 0 then
    return -2
end

local code = redis.call("get", key)
if input_code == code then
    -- 注: 发送间隔1分钟的控制是通过send_code.lua脚本的key的过期时间来控制的，所以不能直接删除key，否则攻击方可以调用再次发送验证码接口进行攻击。
    redis.call("set", cntKey, 0)
    return 0
else
    redis.call("decr", cntKey)
    return -1
end
