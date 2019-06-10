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
  config := &aws.Config{
    Region: aws.String("us-west-2"),
  }

  session, err := session.NewSession(config)
  if err != nil {
    log.Fatal(err)
  }

  service := sqs.New(session)

  reader := bufio.NewReader(os.Stdin)
  queueName := "ExampleQueue"
  messageCount := 0

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
    messageCount += 1

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

  for ; messageCount != 0; messageCount-- {
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
