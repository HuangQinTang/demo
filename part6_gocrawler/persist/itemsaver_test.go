package persist

import (
	"context"
	"crawler/model"
	"encoding/json"
	"github.com/olivere/elastic/v7"
	"testing"
)

func TestSave(t *testing.T) {
	client, err := elastic.NewClient(
		//es服务不是跑在本地的，swtSniff=false, 不维护集群状态
		elastic.SetSniff(false))
	if err != nil {
		t.Errorf(err.Error())
	}

	examples := model.SimpleInfo{
		Url:           "http://album.zhenai.com/u/1968078839",
		Nickname:      "山水有相逢",
		Gender:        "男士",
		Income:        "5001-8000元",
		Place:         "四川",
		Age:           18,
		EducationMate: "大专",
		Marriage:      "未婚",
		Height:        180,
		Introduce:     "本人是一个打工人生于1987年，普普通通的农民家庭",
	}

	id, err := Save(client, model.Zhenai, examples)
	if err != nil {
		t.Errorf(err.Error())
	}

	data, err := client.Get().Index(model.Zhenai).Id(id).Do(context.Background())
	if err != nil {
		t.Errorf(err.Error())
	}

	var result model.SimpleInfo
	err = json.Unmarshal(data.Source, &result)

	if result.Nickname != examples.Nickname {
		t.Errorf("写入数据与查询数据不一致")
	}
}
