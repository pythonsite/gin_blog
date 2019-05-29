module gin_blog

go 1.12

require (
	github.com/astaxie/beego v1.11.1
	github.com/gin-contrib/sessions v0.0.0-20190512062852-3cb4c4f2d615
	github.com/gin-gonic/contrib v0.0.0-20190510065052-87e961e51ccc
	github.com/gin-gonic/gin v1.4.0
	github.com/go-sql-driver/mysql v1.4.1
	github.com/gorilla/feeds v1.1.1
	github.com/jinzhu/gorm v1.9.8
	github.com/kr/pretty v0.1.0 // indirect
	github.com/microcosm-cc/bluemonday v1.0.2
	github.com/pkg/errors v0.8.0
	github.com/pythonsite/iniConfig v0.0.0-20190522083436-7bb1daf30131
	github.com/russross/blackfriday v2.0.0+incompatible
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/snluu/uuid v0.0.0-20130306162636-1dd34a9ad6c0 // indirect
)

replace (
	golang.org/x/net v0.0.0-20181220203305-927f97764cc3 => github.com/golang/net v0.0.0-20181220203305-927f97764cc3
	golang.org/x/net v0.0.0-20190503192946-f4e77d36d62c => github.com/golang/net v0.0.0-20190503192946-f4e77d36d62c
)
