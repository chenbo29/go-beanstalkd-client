package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/chenbo29/go-beanstalkd-client"
	_ "github.com/chenbo29/gostorage"
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
				FollowerNum       int
				FollowerNumActual int
				InfluencerUserId  int
				InfluencerPFId    int
				InfluencerPFDId   int
			}
			err := json.Unmarshal(body, &param)
			if err != nil {
				logger.Printf("error is %s", err)
				return true
			}
			logger.Printf("follower_num is %d", param.FollowerNum)
			var start = 12600
			var end = 105000
			var userNum = 0
			//db, err := sql.Open("mysql", "linkcar:linkcar!)QP@tcp(47.100.225.249:3306)/linkcar_api")
			db, err := sql.Open("mysql", "linkcar:Linux007@tcp(rm-rj94ov6252s8p44dxuo.mysql.rds.aliyuncs.com:3306)/linkcar_main")
			if err != nil {
				log.Fatal(err)
			}
			for i := 0; i < param.FollowerNumActual; i++ {
				rand.Seed(time.Now().UnixNano())
				var userId = rand.Intn(end-start) + start
				_, err = db.Exec("insert into user_follow (uid,follower_id) values (?, ?)", userId, param.InfluencerUserId)
				if err != nil {
					logger.Printf("error is %s", err)
				}
				_, err = db.Exec("update user set follower_num = follower_num + 100 where id = ?", param.InfluencerUserId)
				if err != nil {
					logger.Printf("user error is %s", err)
				}
				userNum++
			}
			_, err = db.Exec("update influencer_post_follower set follower_num_actual = follower_num_actual+?, successful_num=successful_num+1 where id = ?", param.FollowerNumActual, param.InfluencerPFId)
			if err != nil {
				logger.Printf("influencer_post_follower error is %s", err)
			}
			_, err = db.Exec("update influencer_pf_detail set status = 'success', update_time = current_timestamp where id = ?", param.InfluencerPFDId)
			if err != nil {
				logger.Printf("influencer_pf_detail error is %s", err)
			}

			logger.Printf("the user[%d] actual distribution user num is %d", param.InfluencerUserId, userNum)
			db.Close()
			return true
		},
	}
	beans.Run(&executeFunc)
}

func NewLog() *log.Logger {
	var w io.Writer
	logFileName := fmt.Sprintf("./beanstalkf-finish-%s.log", time.Now().Format("2006-01-02"))
	w, _ = os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	return log.New(w, "go beanstalk client ", log.LstdFlags)
}
