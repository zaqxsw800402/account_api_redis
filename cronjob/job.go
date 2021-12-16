package main

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"log"
	"red/mysql"
	"red/redis"
	"strconv"
	"time"
)

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
	redis redis.Database
	db    mysql.DB
	Shut  chan int
}

// Run implement Run() interface to start rsync job
func (j Job) Run() {
	fmt.Println("Start clearing deleted data")

	ctx := context.Background()

	customers, err := j.db.DeleteCustomers()
	if err != nil {
		log.Println("Error deleting customers in redis failed with cronjob" + err.Error())
	}

	for _, customer := range customers {
		err = j.redis.DeleteCustomer(ctx, strconv.Itoa(int(customer.Id)))
		if err != nil {
			log.Println("delete customer in redis failed with cronjob" + err.Error())
		}
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

		j.redis.DeleteAccount(ctx, strconv.Itoa(int(account.AccountId)))
	}
	log.Printf("%s start clear data\n", time.Now())
}
