package initializers

import (
	"fmt"

	"github.com/go-redis/redis"
)



var RC *redis.Client 

func ConnectToRedis(){
	opt, err := redis.ParseURL("rediss://default:AU9CAAIjcDFhNTBmMjYxMDllZjA0NTE5OTJkNTVkMWI5YWFiYTczY3AxMA@alert-jay-20290.upstash.io:6379")
	RC = redis.NewClient(opt)
  
	if err!=nil{
		fmt.Println("Error connecting redis")
		return
	}

	fmt.Println("Connected to redis")
}