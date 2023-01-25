uwsgi-linux-amd64

修改$(go env GOROOT)/src/net/http/status.go

文件中的StatusBadRequest

然后再进行编译

把400 Bad Request改为204 Nothing Yet

原版status.go地址：
https://github.com/golang/go/blob/master/src/net/http/status.go


------
termux里面status.go文件路径为：
/data/data/com.termux/files/usr/lib/go/src/net/http


修改status.go后，直接复制替换：

cp status.go $(go env GOROOT)/src/net/http
