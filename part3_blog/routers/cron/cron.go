package main

import (
	"blog/pkg/setting"
	"github.com/robfig/cron"
	"log"
	"time"

	"blog/models"
)

func main() {
	setting.Setup() //初始化配置
	models.Setup()  //连接mysql
	log.Println("Cron Starting...")

	c := cron.New()
	c.AddFunc("* * * * * *", func() {
		log.Println("Run models.CleanAllTag...")
		models.CleanAllTag()
	})
	c.AddFunc("* * * * * *", func() {
		log.Println("Run models.CleanAllArticle...")
		models.CleanAllArticle()
	})

	c.Start()

	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:
			t1.Reset(time.Second * 10)
		}
	}
}
