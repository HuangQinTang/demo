跟着这篇博客做的练习 
> [Go Web 框架 Gin 实践](https://segmentfault.com/a/1190000013297625#articleHeader5)
> 
> 作者：煎鱼
> 
> 项目地址：https://github.com/EDDYCJY/go-gin-example

项目结构
```
gin-blog/
├── conf
├──  middleware
├── models
├── pkg
├── routers
└──  runtime 
```
- conf：用于存储配置文件
- middleware：应用中间件
- models：应用数据库模型
- pkg：第三方包
- routers 路由逻辑处理
- runtime 应用运行时数据

启动项目

`docker-compose build`，构建go服务镜像

`docker-compose pull`，拉取镜像，主要是mysql,redis,nginx

mysql和redis的密码在.env文件

跑起mysql容器
```
docker-compose up mysql
```

创建数据库及表，复制/conf/blog.sql下的sql文件，执行
```
mysql -u root -p root
#然后粘贴sql执行...

#退出
exit;

#ctrl+c，退出容器
```

nginx的配置文件如下，nginx.conf
```
server {
  listen  80 default_server;
  location / {
    proxy_pass http://blog:10001;
  }
}
```
自行修改在docker-compose.yml中修改nginx配置文件，日志文件的挂载路径，还有mysql的。

启动服务

`docker-compose up -d`

如下表示容器启动成功，尝试访问`http://localhost/`
```
Starting gin_blog_nginx_1        ... done
Starting gin_blog_blog_1         ... done
Starting gin_blog_redis-server_1 ... done
Starting gin_blog_mysql-server_1 ... done
```

- 备注
- 这里我采用的是把mysql,redis容器名和密码作为环境变量定义在blog(go服务)容器，代码里面采用`os.Getenv(环境变量名)`来连接mysql和redis。
- go服务的构建文件是Dockerfile-new

go应用的镜像可以下载curl，添加如下健康检查
```azure
healthcheck:
  test: [ "CMD", "curl", "-f", "http://localhost:10001" ]
  interval: 30s
  timeout: 3s
  retries: 3
  start_period: 40s
```