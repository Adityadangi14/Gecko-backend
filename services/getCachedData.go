package services

import (
	"fmt"
	"gecko_backend/initializers"

	"github.com/go-redis/redis"
)

func GetCachedData(key string)(string,error) {

	val,err:=	initializers.RC.Get(key).Result()


	if err == redis.Nil {
		fmt.Println("Key does not exist.")
		return "",fmt.Errorf("Key does not exist")
	} else if err != nil {
		fmt.Println("Error fetching value:", err)

		return "",fmt.Errorf("error fetching value: %v", err)
	} else {
		return val, nil
	}
}