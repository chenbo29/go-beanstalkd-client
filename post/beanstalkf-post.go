package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/chenbo29/go-beanstalkd-client"
	_ "github.com/chenbo29/gostorage"
	_ "github.com/go-redis/redis"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	executeFunc := beans.JobExecuteFunc{
		Execute: func(id uint64, body []byte) bool {
			logger := NewLog()
			logger.Printf("jobID(%d) start to do jobBody(%s)", id, string(body))
			var param struct {
				Num  int
				Type string
				Id   int
			}
			err := json.Unmarshal(body, &param)
			if err != nil {
				logger.Printf("error is %s", err)
				return true
			}
			logger.Printf("Type %s Id[%d] Num[%d]", param.Type, param.Id, param.Num)
			var start = 12600
			var end = 105000
			var num = 0
			db, err := sql.Open("mysql", "linkcar:linkcar!)QP@tcp(47.100.225.249:3306)/linkcar_api")
			//db, err := sql.Open("mysql", "linkcar:Linux007@tcp(rm-rj94ov6252s8p44dxuo.mysql.rds.aliyuncs.com:3306)/linkcar_main")
			if err != nil {
				log.Fatal(err)
			}

			switch param.Type {
			case "support":
				for i := 0; i < param.Num; i++ {
					rand.Seed(time.Now().UnixNano())
					var userId = rand.Intn(end-start) + start
					_, err = db.Exec("insert into post_support (uid,post_id) values (?, ?)", userId, param.Id)
					if err != nil {
						logger.Printf("error is %s", err)
					}
					num++
				}
				_, err = db.Exec("update post set support_count = (select count(*) from post_support where post_id = ?) where id = ?", param.Id, param.Id)
				if err != nil {
					logger.Printf("update support_count error is %s", err)
				}
			case "comment":
				for i := 0; i < param.Num; i++ {
					rand.Seed(time.Now().UnixNano())
					var userId = rand.Intn(end-start) + start
					_, err = db.Exec("insert into post_comment (uid,post_id,content) values (?, ?, (select content from post_comment_template order by RAND() limit 1))", userId, param.Id)
					if err != nil {
						logger.Printf("error is %s", err)
					}
					num++
				}
				_, err = db.Exec("update post set comment_count = (select count(*) from post_comment where post_id = ?) where id = ?", param.Id, param.Id)
				if err != nil {
					logger.Printf("update support_count error is %s", err)
				}
			case "share":
				for i := 0; i < param.Num; i++ {
					_, err = db.Exec("update post set repost_count = repost_count + 1 where id = ?", param.Id)
					if err != nil {
						logger.Printf("update support_count error is %s", err)
					}
					num++
				}
			default:

			}
			logger.Printf("the post[%d] actual get %s num is %d", param.Id, param.Type, num)
			db.Close()
			return true
		},
	}
	beans.RunByTubeName(&executeFunc, "post")
}

func NewLog() *log.Logger {
	var w io.Writer
	logFileName := fmt.Sprintf("./beanstalkf-post-finish-%s.log", time.Now().Format("2006-01-02"))
	w, _ = os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	return log.New(w, "go beanstalk client ", log.LstdFlags)
}
