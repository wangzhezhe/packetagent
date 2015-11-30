package app

import (
	"errors"
	"github.com/emicklei/go-restful"
	"github.com/golang/glog"
	"github.com/shirou/gopsutil/process"
	"packetagent/lib"
	"packetagent/model"
	"strconv"
	"time"
)

// This example shows the minimal code needed to get a restful.WebService working.
//
// GET http://localhost:8080/hello

var (
	device string = "eth0"
)

func Register(container *restful.Container) {
	if glog.V(1) {
		glog.Info("start regist")
	}
	ws := new(restful.WebService)
	ws.Path("/packet").
		Doc("control packet listening").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Route(ws.GET("/{port-id}").
		To(StartListen).
		Doc("start collect the packet from the specified port").
		Operation("findUser").
		Param(ws.PathParameter("port-id", "identifier of the listening port").DataType("string")))

	container.Add(ws)
}

//when get the pid of that port
//start collecting the packet for 60s
func StartListen(request *restful.Request, response *restful.Response) {
	portstring := request.PathParameter("port-id")
	glog.Info("get the port number", portstring)
	portint, err := strconv.Atoi(portstring)
	if err != nil {
		response.WriteError(500, err)
		return
	}
	pid, pname, err := lib.Getinfofromport(portint)

	if pid == -1 {
		response.WriteError(500, errors.New("the port is not be listend in this machine ( /proc/net/tcp and /proc/net/tcp6)"))
		return
	}

	if err != nil {
		response.WriteError(500, err)
		return

	}
	glog.Info(pname, pid)

	//create the process instance and get the detail info of specified pid
	Pdetail := &model.ProcessDetail{
		Process: &process.Process{Pid: 22637},
	}
	cmd, err := Pdetail.Cmdinfo()
	if err != nil {
		glog.Info(err)
	}
	glog.Info(cmd)
	//TODO get more info of this instance

	//start listen to specific ip:port for 60s and send the data to es
	timesignal := time.After(time.Second * 30)

	//start collect and check the timesignal every one minutes
	go lib.Startcollect(portint, device, timesignal)

	response.Write([]byte("activated"))

}
