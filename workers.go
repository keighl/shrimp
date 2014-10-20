package main

import (
  "github.com/jrallison/go-workers"
)

// TODO Make this read from a Configuration object
func ConfigureWorkerServer(forTest bool)  {

  database := "0"
  if (forTest) { database = "15" }

   workers.Configure(map[string]string{
    "server": "localhost:6379",
    "database": database,
    "pool": "30",
    "process": "1",
  })

  // // Set a queue
  // workers.Process("dummyQueue", dummyJob, 5)
  //
  // // Trigger a background job on the queue
  // workers.Enqueue("dummyQueue", "Add", user.Id.Int64)
  //
  // // Handle the job when it comes off the queue
  // func dummyJob(message *workers.Msg) {
  //   _, err := message.Args().Int64()
  //   if (err != nil) { panic(err) }
  // }
}



