package htttp

import "htttp/repositoriy"

func main() {
	pol := repositoriy.NewPolzovately()
	httpp := NewHTTPHandler(pol)
	serv := NewHTTPServer(httpp)

	serv.Start()
}
