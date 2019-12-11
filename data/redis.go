package data

import (
	"SecKill/conf"
	"github.com/go-redis/redis/v7"
)

var client *redis.Client

// 开启redis连接池
func initRedisConnection(config conf.AppConfig)  {
	client = redis.NewClient(&redis.Options{
		Addr:     config.App.Redis.Address,
		Password: config.App.Redis.Password, // It's ok if password is "".
		DB:       0,  // use default DB
	})
}

// 确保redis加载lua脚本，若未加载则加载
func PrepareScript(script string) string {
	// sha := sha1.Sum([]byte(script))
	scriptsExists, err := client.ScriptExists(script).Result()
	if err != nil {
		panic("Failed to check if script exists: " + err.Error())
	}
	if !scriptsExists[0] {
		scriptSHA, err := client.ScriptLoad(script).Result()
		if err != nil {
			panic("Failed to load script " + script + " err: " + err.Error())
		}
		return scriptSHA
	}
	print("Script Exists.")
	return ""
}

// 执行lua脚本
func EvalSHA(sha string, args []string) (interface{}, error) {
	val, err := client.EvalSha(sha, args).Result()
	if err != nil {
		print("Error executing evalSHA... " + err.Error())
		return nil, err
	}
	return val, nil
}

// redis operation SET
func SetForever(key string, value interface{}) (string, error) {
	val, err := client.Set(key, value, 0).Result()  // expiration表示无过期时间
	return val, err
}

// redis operation hmset
func SetMapForever(key string, field map[string]interface{}) (string, error) {
	return client.HMSet(key, field).Result()
}

// redis operation hmget
func GetMap(key string, fields ...string) ([]interface{}, error) {
	return client.HMGet(key, fields...).Result()
}