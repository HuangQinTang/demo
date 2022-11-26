package routers

import (
	"blog/middleware/jwt"
	"blog/pkg/upload"
	"blog/routers/api"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"net/http"

	"blog/pkg/setting"
	"blog/routers/api/v1"
)

func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger()) //日志中间件

	r.Use(gin.Recovery()) //避免panic导致程序退出

	gin.SetMode(setting.ServerSetting.RunMode) //是否debug模式

	//swagger文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	//获取token
	r.GET("/auth", api.GetAuth)
	//上传文件
	r.StaticFS("/upload/images", http.Dir(upload.GetImageFullPath()))
	//用于健康检查
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"code": http.StatusOK,
			"msg":  "success",
			"data": nil,
		})
	})

	apiv1 := r.Group("/api/v1").Use(jwt.JWT())
	{
		//获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		//新建标签
		apiv1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)

		//获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		//获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		//新建文章
		apiv1.POST("/articles", v1.AddArticle)
		//更新指定文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		//删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)

		//文件上传
		apiv1.POST("/upload", api.UploadImage)
	}

	return r
}
