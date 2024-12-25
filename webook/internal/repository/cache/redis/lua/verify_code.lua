local key = KEYS[1]
local cntKey = key..":cnt"
-- 用户输入的 code
local expectedCode = ARGV[1]

local code = redis.call("get", key)
--转成一个数字
local cnt = tonumber(redis.call("get", cntKey))
if cnt == nil or cnt <= 0 then
    --说明用户出错超过限值,或者已经用过了
    return -1
elseif expectedCode == code then
    --输入正确
    --删除对应记录
    redis.call("set", cntKey, 0)
    return 0
else
    --输错了
    redis.call("decr", cntKey)
    return -2
end