package main

import (
	"Atlan_Collect_Ingest/repository"
	"Atlan_Collect_Ingest/usecase/googleSheet"
	"Atlan_Collect_Ingest/usecase/jobStatus"
	"Atlan_Collect_Ingest/usecase/sendSms"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/streadway/amqp"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		fmt.Println("Failed Initializing Broker Connection")
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println(err)
	}
	defer ch.Close()

	DATABASE_URL := os.Getenv("JOB_DATABASE_URL")
	ctx := context.Background()
	fmt.Printf("\n %v\n", DATABASE_URL)
	cofig, err := pgxpool.ParseConfig(DATABASE_URL)
	if err != nil {
		fmt.Printf("Error : %v", err)
	}

	pool, err := pgxpool.ConnectConfig(ctx, cofig)
	if err != nil {
		fmt.Printf("Error : %v", err)
	}
	defer pool.Close()

	if err != nil {
		fmt.Println(err)
	}

	msgs, err := ch.Consume(
		"PluginQueue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	sheetRepo := repository.NewGoogleSheetDbSql(pool)
	sheetService := googleSheet.NewService(sheetRepo)

	jobRepo := repository.NewJobDbSql(pool)
	jobService := jobStatus.NewService(jobRepo)

	smsRepo := repository.NewSMSDbSql(pool)
	smsService := sendSms.NewService(smsRepo)

	forever := make(chan bool)

	go func() {
		//msg format "formid,jobId,code,pluginCode"
		for d := range msgs {
			fmt.Printf("Recieved Message: %s\n", d.Body)
			body := strings.Split(string(d.Body), ",")

			jobId, err := strconv.Atoi(body[1])

			formId, err := strconv.Atoi(body[0])
			if err != nil {
				log.Printf("%v\n", err.Error())
				jobService.UpdateJob(jobId, 500, err.Error())
				continue
			}
			pluginCode, err := strconv.Atoi(body[3])
			if err != nil {
				log.Printf("%v\n", err.Error())
				jobService.UpdateJob(jobId, 500, err.Error())
				continue
			}

			if pluginCode == 1 {
				_, err = sheetService.AddToSheet(int8(formId), body[2])
				if err != nil {
					log.Printf("%v\n", err.Error())
					jobService.UpdateJob(jobId, 500, err.Error())
					continue
				}
			} else if pluginCode == 2 {

				err := smsService.SendSMS()
				if err != nil {
					log.Printf("%v\n", err.Error())
					continue
				}
			}

			jobService.UpdateJob(jobId, 200, "Completed")

		}
	}()

	fmt.Println("Successfully Connected to our RabbitMQ Instance")
	fmt.Println(" [*] - Waiting for messages")
	<-forever

}

/*func main() {
DATABASE_URL := os.Getenv("JOB_DATABASE_URL")
ctx := context.Background()
cofig, err := pgxpool.ParseConfig(DATABASE_URL)
if err != nil {
	fmt.Printf("Error : %v", err)
}

pool, err := pgxpool.ConnectConfig(ctx, cofig)
if err != nil {
	fmt.Printf("Error : %v", err)
}
defer pool.Close()

if err != nil {
	fmt.Println(err)
}

repo := repository.NewGoogleSheetDbSql(pool)
serv := googleSheet.NewService(repo)
serv.AddToSheet(98)
/*var resp = responses[1]
fmt.Printf("id = %v Q1 = %v Q2 = %v Q3 = %v \n", resp.UserId, resp.Response[14], resp.Response[15], resp.Response[17])
}
*/
