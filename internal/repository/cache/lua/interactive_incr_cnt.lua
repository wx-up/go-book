local key = KEYS[1]
local cntKey = ARGV[1]
local delta = tonumber(ARGV[2])

local exists = redis.call("EXISTS", cntKey)
if exists == 1 then
    redis.call("HINCRBY", key,cntKey,delta)
    return 1
else
    return 0
end