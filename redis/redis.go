package redis

import (
	"delonakc.com/api/config"
	"encoding/json"
	"fmt"
	rd "github.com/go-redis/redis"
	"log"
	"time"
)

var Rd *rd.Client
var expireTime time.Duration

func Start(conf *config.Conf) {

	url := fmt.Sprintf("redis://%s:%d", conf.Redis.Host, conf.Redis.Port)

	options, err := rd.ParseURL(url)

	if err != nil {
		log.Fatal(err)
	}
	expireTime, err = time.ParseDuration(conf.Redis.ExpireTime)
	if err != nil {
		log.Println(err)
	}
	Rd = rd.NewClient(options)

	_, err = Rd.Ping().Result()

	log.Printf("url: %s\n", url)
	if err != nil {
		log.Println("redis connect failed!!!")
		log.Fatal(err)
	}
	fmt.Println("=====================")
	fmt.Println("redis connect success")
	fmt.Println("=====================")
}

func Set(key string, value interface{}) error {
	if str, ok := value.(string); ok {
		Rd.Set(key, str, expireTime)
		return nil
	}
	v, err := json.Marshal(value)

	if err == nil {
		Rd.Set(key, fmt.Sprintf("%s", v), expireTime)
		return nil
	}

	return err
}

func Get(key string) []byte {
	v := Rd.Get(key)

	res, err := v.Bytes()

	if err != nil {
		return nil
	}

	return res
}