if (redis.call('exists', KEYS[1]) == 0) then
    redis.call('rpush', KEYS[1], ARGV[1]);
    return redis.call('pexpire', KEYS[1], ARGV[2]);
else
    return redis.call('rpushx', KEYS[1], ARGV[1]);
end;
