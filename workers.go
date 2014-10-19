package main

import (
  "github.com/jrallison/go-workers"
)

func ConfigureWorkerServer()  {
  // Workers
   workers.Configure(map[string]string{
    "server": "localhost:6379",
    "database": "0",
    "pool": "30",
    "process": "1",
  })

  workers.Process("dummyQueue", dummyJob, 5)
}

func dummyJob(message *workers.Msg) {
  _, err := message.Args().Int64()

  if (err != nil) {
    panic(err)
  }
}