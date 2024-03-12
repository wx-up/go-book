local key = KEYS[1]
local cntKey = key..":cnt"
local inputCode = ARGV[1]

local cnt = redis.call("GET",cntKey)
-- key 不存在时返回 false
if cnt == false then
    return -3
end

local code = redis.call("GET",key)
if code == false then
    return -3
end

cnt = tonumber(cnt)
if (cnt <= 0) then
    -- 用户一直输错
    return -1
end

if (code == inputCode) then
    redis.call("DEL",key)
    redis.call("DEL",cntKey)
    return 0
else
    -- 验证码错误
    redis.call("DECR",cntKey)
    return -2
end
