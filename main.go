package main

import (
  "os"
  "fmt"
  "log"
  "bufio"
  "strings"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/sqs"
)

func main() {
  session, err := session.NewSession(&aws.Config{
    Region: aws.String("us-west-2"),
  })
  if err != nil {
    log.Fatal(err)
  }

  service := sqs.New(session)

  reader := bufio.NewReader(os.Stdin)
  queueName := "ExampleQueue"
  count := 0

  println("Creating queue...")

  queue, err := service.CreateQueue(&sqs.CreateQueueInput{
    QueueName: &queueName,
  })
  if err != nil {
    log.Fatal(err)
  }

  print("Queue URL: ")
  println(*queue.QueueUrl)

  defer (func() {
    _, err = service.DeleteQueue(&sqs.DeleteQueueInput{
      QueueUrl: queue.QueueUrl,
    })
    if err != nil {
      log.Fatal(err)
    }
 })()

  println("type and send messages to SQS by hitting enter (type 'q' to quit):")

  for {
    count += 1

    fmt.Print("> ")
    text, _ := reader.ReadString('\n')
    message := strings.TrimRight(text, "\n")

    if message == "q" {
      break
    }

    _, err = service.SendMessage(&sqs.SendMessageInput{
      MessageBody: &message,
      QueueUrl: queue.QueueUrl,
    })
    if err != nil {
      log.Fatal(err)
    }
  }

  for ; count != 0; count-- {
    messages, err := service.ReceiveMessage(&sqs.ReceiveMessageInput{
      QueueUrl: queue.QueueUrl,
    })
    if err != nil {
      log.Fatal(err)
    }

    for _, msg := range messages.Messages {
      println(*msg.Body)
    }
  }
}
