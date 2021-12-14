package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"red/mysql"
	"time"
)

func main() {
	fmt.Println("Start clearing deleted data")

	dbClient := mysql.GetDBClient()
	db := mysql.NewDB(dbClient)
	job := Job{
		db:   db,
		Shut: make(chan int, 1),
	}

	// clear data at every 6 am
	go StartJob("* 6 * * *", job)
	select {}
}

func StartJob(spec string, job Job) {
	c := cron.New()
	c.AddJob(spec, &job)

	// 啓動執行任務
	c.Start()
	// 退出時關閉計劃任務
	defer c.Stop()

	// 如果使用 select{} 那麼就一直會循環
	select {
	case <-job.Shut:
		return
	}
}

func StopJob(shut chan int) {
	shut <- 0
}

type Job struct {
	db   mysql.DB
	Shut chan int
}

// Run implement Run() interface to start rsync job
func (j Job) Run() {
	err := j.db.DeleteCustomers()
	if err != nil {
		log.Println("Error deleting customers" + err.Error())
	}
	accounts, err := j.db.DeleteAccounts()
	if err != nil {
		log.Println("Error deleting accounts" + err.Error())
	}
	for _, account := range accounts {
		err = j.db.DeleteTransactions(account.AccountId)
		if err != nil {
			log.Println("Error deleting transactions" + err.Error())
		}
	}
	log.Printf("%s start clear data\n", time.Now())
}
