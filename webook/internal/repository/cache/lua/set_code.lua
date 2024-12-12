--你的验证码在 Redis 上的 key
-- phone_code:login:152xxxxxxxx
local key = KEYS[1]
-- 记录验证次数
-- phone_code:login:152xxxxxxxx:cnt
local cntKey = key..":cnt"
-- 你的验证码 123456
local val = ARGV[1]
-- 过期时间
local ttl = tonumber(redis.call("ttl", key))
if ttl == -1 then
    -- key 存在，但没有过期时间
    -- 系统错误，没有过期时间
    return -2
    -- 540 = 600 - 60 即 9 分钟
elseif ttl == -2 or ttl <540 then
    redis.call("set", key, val)
    redis.call("expire", key, 600)
    redis.call("expire", cntKey, 600)
    redis.call("set", cntKey, 3)
    -- 符合预期
    return 0
else
    -- 发送太频繁
    return -1
end