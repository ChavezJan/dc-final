package controller

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ChavezJan/dc-final/scheduler"
	"go.nanomsg.org/mangos"
	"go.nanomsg.org/mangos/protocol/pub"

	// register transports
	_ "go.nanomsg.org/mangos/transport/all"
)

var controllerAddress = "tcp://localhost:40899"

type workloads struct {
	workload_id   string
	workload_name string
	wl_status     bool
	wl_filter     string
}

type savedImages struct {
	image_file_name string
	image_ID        string
	image_type      string
}

var images []savedImages
var workloadCtrl []workloads

func SaveWorkload(name string, token string, status bool, filtro string) {

	worker := workloads{
		workload_id:   token,
		workload_name: name,
		wl_status:     status,
		wl_filter:     filtro,
	}
	/*
		if len(workloadCtrl) > 0 {
			worker := append(workloadCtrl, len(workloadCtrl)-1)
		}
		if len(workloadCtrl) == nil {

		}*/

	fmt.Println(worker)

}

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func date() string {
	return time.Now().Format(time.ANSIC)
}

func Active_workloads() string {

	var disponibles string

	disponibles = scheduler.Active_workloads()

	return disponibles
}

func Start() {
	var sock mangos.Socket
	var err error
	if sock, err = pub.NewSocket(); err != nil {
		die("can't get new pub socket: %s", err)
	}
	if err = sock.Listen(controllerAddress); err != nil {
		die("can't listen on pub socket: %s", err.Error())
	}
	for {
		// Could also use sock.RecvMsg to get header
		d := date()
		log.Printf("Controller: Publishing Date %s\n", d)
		if err = sock.Send([]byte(d)); err != nil {
			die("Failed publishing: %s", err.Error())
		}
		time.Sleep(time.Second * 3)
	}
}
