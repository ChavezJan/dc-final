# User Guide

#User Guide
_Distributed and Parallel Image Processing_

To access Distributed and Parallel Image Processing, you must access the following github link. Which directs to a repository that has to be forked
https://github.com/ChavezJan/dc-final

In a terminal you must create a folder with the following path:

    src/go/github.com/githubUser/

And inside that directory in the terminal you must clone the repository obtaining the httpURL in github

     git clone https://github.com/ChavezJan/dc-final.git

In the dc-final folder there is a file called "main.go" in the dependencies section there are some that come from github, the username must be modified to your own github username.

In order to be able to run the API, it is necessary to install Gin. It is a HTTP web framework, which helps to build applications and microservices in Go (Golang). To install Gin run in a terminal the following command:

    $ go get -u github.com/gin-gonic/gin

Once Gin is installed, you must use the following commands in the console:

    $ export GO111MODULE=off
    $ go run main.go

After running the server in a different terminal client requests may be sent. One of them is access providing username, password and port to be use:

    $ curl -u username:password http://localhost:8080/login

Receiving a similar JSON message with a access token:

> "message": "Hi username, welcome to the DPIP System",

> "token" "OjIE89GzFw"

Other client request is status adding the access token given previously:

    $ curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/status

Receiving JSON message with exact time and date:

> "message": "Hi username, the DPIP System is Up and Running"

> "time": "2015-03-07 11:06:39"

The third possible client request is upload, which allows the client to open a local file using the file path and the access token:

    $ curl -F 'data=@/path/to/local/image.png' -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/images

Receiving JSON message with the file name and size:

> "message": "An image has been successfully uploaded",
> "filename": "image.png",
> "size": "500kb"

An endpoint that connects to the controller is workloads that creates a new workload with filter parameters and the name of the workload.

    $ curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/workloads

Receiving JSON message with workload's information:

> "workload_id": 7e,

    "filter": grayscale,
    "workload_name" : "Trabajador 1",
    "status" : "scheduling"

Another way to access information from the controller is with an endpoint together with an id, to obtain specific information about a workload.

    $ curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/workloads/:workload_id

Receiving JSON message with exact time and date:

> "workload_id": 7e,

    "filter": grayscale,
    "workload_name" : "Trabajador 1",
    "status" : "scheduling"

Finally the last cliente request is logout using the access token:

    $ curl -H "Authorization: Bearer <ACCESS_TOKEN>" http://localhost:8080/logout

Receiving JSON message with the following:

> "message": "Bye username, your token has been revoked"
