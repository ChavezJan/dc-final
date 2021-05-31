package scheduler

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "github.com/ChavezJan/dc-final/proto"
	"google.golang.org/grpc"
)

//const (
//	address     = "localhost:50051"
//	defaultName = "world"
//)

func Active_workloads() string {
	var disponibles string

	fmt.Println("aqui estoy ")

	for i := 0; i < 5; i++ {
		/*
			checar si esta disponible
			concatenar los trabajadores acticos con su ID
		*/
		disponibles += "listas"
		disponibles += "/"
	}

	return disponibles
}

type Job struct {
	Address string
	RPCName string
}

func schedule(job Job) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(job.Address, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	// reparticion de trabajo a los workers
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: job.RPCName})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Scheduler: RPC respose from %s : %s", job.Address, r.GetMessage())
}

func Start(jobs chan Job) error {
	for {
		job := <-jobs
		schedule(job)
	}
	return nil
}
