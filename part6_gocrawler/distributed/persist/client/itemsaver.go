package client

import (
	"crawler/distributed/config"
	"crawler/distributed/rpcsupport"
	"crawler/model"
	"crawler/rediservice"
	"fmt"
	"log"
)

func ItemSeaver(host string) (chan model.SimpleInfo, error) {
	client, err := rpcsupport.NewClient(host)
	if err != nil {
		return nil, err
	}
	out := make(chan model.SimpleInfo)
	go func() {
		itemCount := 0
		for {
			item := <-out

			//判断是否已保存过
			id := fmt.Sprintf("%d", item.Id)
			unique, err := rediservice.RedisServer.ZhenaiDuplicate(id)
			if err != nil {
				log.Printf("Item Saver err: 【%v】", err)
				continue
			}
			if !unique {
				log.Printf("Item Saver info: 用户【%v】已存在", id)
				continue
			}

			result := ""
			err = client.Call(config.ItemSaverRpc, item, &result)
			if err != nil {
				log.Printf("Item Save fail: 【%v】", err)
				continue
			}
			log.Printf("Item Saver: got item: #%d：【%v】", itemCount, item)
			itemCount++
		}
	}()
	return out, nil
}
