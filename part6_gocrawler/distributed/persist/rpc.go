package persist

import (
	"crawler/model"
	"crawler/persist"
	"github.com/olivere/elastic/v7"
	"log"
)

type ItemSaverService struct {
	Client *elastic.Client
	EsIndex  string
}

func (s *ItemSaverService) Save(item model.SimpleInfo, result *string) error {
	docId, err := persist.Save(s.Client, s.EsIndex, item)
	log.Printf("Save Item %v", item)
	if err == nil {
		*result = docId
	} else {
		log.Printf("Error saving item %v: %v", item, err)
	}
	return err
}
