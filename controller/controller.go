package controller

import (
	"fmt"
	"log"
	"os"
	"time"

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
	imgEditada    []string
}

type savedImages struct {
	image_file_name string
	image_ID        string
	image_type      string
}

var images []savedImages
var imgIDs []string
var workloadCtrl []workloads

func SaveWorkload(name string, token string, status bool, filtro string) {

	worker := workloads{
		workload_id:   token,
		workload_name: name,
		wl_status:     status,
		wl_filter:     filtro,
	}

	workloadCtrl = append(workloadCtrl, worker)
}

func GetWorkers(name string) bool {
	for i := range workloadCtrl {
		if workloadCtrl[i].workload_name == name {
			return true
		}
	}
	return false
}

func GetImgIDs(wkName string) []string {

	for i := range images {
		if images[i].image_type == "Original" {
			images[i].image_type = "Editada"

			for x := range workloadCtrl {
				if workloadCtrl[x].workload_name == wkName {
					workloadCtrl[x].imgEditada = append(workloadCtrl[x].imgEditada, images[i].image_ID)
					workloadCtrl[x].wl_status = !workloadCtrl[x].wl_status
					return workloadCtrl[x].imgEditada
				}
			}

		}
	}
	return nil
}
func GetStatus(wkName string) bool {

	for i := range workloadCtrl {
		if workloadCtrl[i].workload_name == wkName {
			return workloadCtrl[i].wl_status
		}
	}

	return false
}

func SaveImage(name string, token string, imgType string) {

	image := savedImages{
		image_file_name: name,
		image_ID:        token,
		image_type:      imgType,
	}

	images = append(images, image)
	imgIDs = append(imgIDs, image.image_ID)

}

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

func date() string {
	return time.Now().Format(time.ANSIC)
}

func Active_workloads() string {

	var workerInActive string

	for x := range workloadCtrl {
		if workloadCtrl[x].wl_status == true {
			workerInActive += workloadCtrl[x].workload_id
			workerInActive += "/"
		}
	}

	return workerInActive
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
