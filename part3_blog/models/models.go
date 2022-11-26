package models

import (
	"blog/pkg/logging"
	"blog/pkg/setting"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var once sync.Once

type Model struct {
	ID         int `gorm:"primary_key" json:"id"`
	CreatedOn  int `json:"created_on"`
	ModifiedOn int `json:"modified_on"`
	DeletedOn  int `json:"deleted_on"`
}

func Setup() {
	once.Do(func() {
		host := os.Getenv(setting.DatabaseSetting.Host) + ":" + setting.DatabaseSetting.Port
		psw := os.Getenv(setting.DatabaseSetting.Password)
		var err error
		db, err = gorm.Open(setting.DatabaseSetting.Type, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
			setting.DatabaseSetting.User,
			psw,
			host,
			setting.DatabaseSetting.Name))

		if err != nil {
			fmt.Println(psw)
			fmt.Println(err.Error())
			logging.Fatal(err)
		}

		gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
			return setting.DatabaseSetting.TablePrefix + defaultTableName
		}

		db.SingularTable(true)       //gorm默认使用复数映射，true表示严格匹配不走默认复数映射
		db.LogMode(true)             //打印sql
		db.DB().SetMaxIdleConns(10)  //空闲连接数
		db.DB().SetMaxOpenConns(100) //最大连接数

		//注册回调方法，用于更新 创建时间、更新时间、删除时间 字段
		db.Callback().Create().Replace("gorm:update_time_stamp", updateTimeStampForCreateCallback)
		db.Callback().Update().Replace("gorm:update_time_stamp", updateTimeStampForUpdateCallback)
		db.Callback().Delete().Replace("gorm:delete", deleteCallback)
	})
}

func CloseDB() {
	defer db.Close()
}

// updateTimeStampForCreateCallback will set `CreatedOn`, `ModifiedOn` when creating
func updateTimeStampForCreateCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		nowTime := time.Now().Unix()
		if createTimeField, ok := scope.FieldByName("CreatedOn"); ok { //是否存在字段,created_on也行
			if createTimeField.IsBlank { //是否为空
				createTimeField.Set(nowTime)
			}
		}

		if modifyTimeField, ok := scope.FieldByName("ModifiedOn"); ok {
			if modifyTimeField.IsBlank {
				modifyTimeField.Set(nowTime)
			}
		}
	}
}

// updateTimeStampForUpdateCallback will set `ModifyTime` when updating
func updateTimeStampForUpdateCallback(scope *gorm.Scope) {
	if _, ok := scope.Get("gorm:update_column"); !ok { //根据入参获取设置了字面值的参数
		scope.SetColumn("ModifiedOn", time.Now().Unix()) //没有指定 update_column 的字段值，这里补充上
	}
}

func deleteCallback(scope *gorm.Scope) {
	if !scope.HasError() {
		var extraOption string
		if str, ok := scope.Get("gorm:delete_option"); ok { //是否手动指定了delete_option
			extraOption = fmt.Sprint(str)
		}
		deletedOnField, hasDeletedOnField := scope.FieldByName("DeletedOn") //获取我们约定的删除字段

		if !scope.Search.Unscoped && hasDeletedOnField { //若存在则 UPDATE 软删除
			scope.Raw(fmt.Sprintf(
				"UPDATE %v SET %v=%v%v%v",
				scope.QuotedTableName(),                            //返回引用的表名
				scope.Quote(deletedOnField.DBName),                 //字段名
				scope.AddToVars(time.Now().Unix()),                 //添加值作为SQL的参数，也可防范SQL注入
				addExtraSpaceIfExist(scope.CombinedConditionSql()), //scope.CombinedConditionSql() 返回组合好的条件SQL
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		} else { //不存在则 DELETE 硬删除
			scope.Raw(fmt.Sprintf(
				"DELETE FROM %v%v%v",
				scope.QuotedTableName(),
				addExtraSpaceIfExist(scope.CombinedConditionSql()),
				addExtraSpaceIfExist(extraOption),
			)).Exec()
		}
	}
}

func addExtraSpaceIfExist(str string) string {
	if str != "" {
		return " " + str
	}
	return ""
}
