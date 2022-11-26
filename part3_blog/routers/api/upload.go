package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"blog/pkg/e"
	"blog/pkg/logging"
	"blog/pkg/upload"
)

func UploadImage(c *gin.Context) {
	code := e.SUCCESS
	data := make(map[string]string)

	file, image, err := c.Request.FormFile("image") //获取上传的图片
	if err != nil {
		logging.Warn(err)
		code = e.ERROR
		c.JSON(http.StatusOK, gin.H{
			"code": code,
			"msg":  e.GetMsg(code),
			"data": data,
		})
	}

	if image == nil {
		code = e.INVALID_PARAMS
	} else {
		imageName := upload.GetImageName(image.Filename) //保留后缀,文件名md5
		fullPath := upload.GetImageFullPath()
		savePath := upload.GetImagePath()

		src := fullPath + imageName
		if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) { //检查文件后缀,文件大小
			code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		} else {
			err = upload.CheckImage(fullPath) //判断路径存不存在,有没有权限,没问文件夹创建并赋予权限
			if err != nil {
				logging.Warn(err)
				code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
			} else if err = c.SaveUploadedFile(image, src); err != nil { //上传报错
				logging.Warn(err)
				code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
			} else { //上传成功,重复上传同名文件会覆盖
				data["image_url"] = upload.GetImageFullUrl(imageName)
				data["image_save_url"] = savePath + imageName
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"code": code,
		"msg":  e.GetMsg(code),
		"data": data,
	})
}
