package main

import (
	htps "htttp/htps"
	"htttp/htps/repositoriy"
)

func main() {
	pol := repositoriy.NewPolzovately()
	httpp := htps.NewHTTPHandler(pol)
	serv := htps.NewHTTPServer(httpp)

	serv.Start()
}
