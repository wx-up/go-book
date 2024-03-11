-- 缓存 key 例如： code:login:156***76
local key = KEYS[1]

-- 验证次数，最多验证 3 次
local cntKey = key..":cnt"

-- 验证码的值
local value = ARGV[1]

-- 过期时间，单位秒
local ttl = tonumber(redis.call("tll", key))

if (ttl == -1) then
    -- key 存在但是没有过期时间（ 可能是同事误操作了 ）
    -- 系统错误
    return -2
-- -2 表示 key 不存在
-- <540 表示发送验证码已经过去1分钟了（ 有效期10分钟，一分钟内只能发送一次 ）
elseif ttl == -2 or ttl < 540 then
    redis.call("set",key,value)
    redis.call("expire", key,600)
    redis.call("set",cntKey,3)
    redis.call("expire", cntKey,600)
    return 0
else
    -- 验证码发送过于频繁
    return -1
end