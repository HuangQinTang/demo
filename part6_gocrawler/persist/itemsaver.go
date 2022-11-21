package persist

import (
	"context"
	"crawler/model"
	"crawler/rediservice"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
)

func SimpleInfoSeaver(client *elastic.Client) chan model.SimpleInfo {
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

			//保存数据至es
			if _, err = Save(client, model.Zhenai, item); err != nil {
				log.Printf("Item Save fail: 【%v】", err)
			}
			log.Printf("Item Saver: got item: #%d：【%v】", itemCount, item)
			itemCount++
		}
	}()
	return out
}

// Save 存储简单用户信息
func Save(client *elastic.Client, index string, userInfo model.SimpleInfo) (docId string, err error) {
	resp, err := client.Index().Index(index).Id(fmt.Sprintf("%d", userInfo.Id)).BodyJson(userInfo).Do(context.Background())
	if err != nil {
		return "", err
	}

	return resp.Id, nil
}
