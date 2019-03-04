package httpserver

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func RunAndWaitGracefulShutdown(server *http.Server) {
	chBack := make(chan error)
	go func() {
		_ = server.ListenAndServe()
		fmt.Println("server shutdown")
		time.Sleep(time.Second * 1)
		fmt.Println("server shutdown delay ended")
		chBack <- nil
	}()

	// Handle SIGINT and SIGTERM.
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch

	// Stop the service gracefully.
	fmt.Println("start shutdown")
	fmt.Println("shutdown result :", server.Shutdown(nil))
	fmt.Println("shutdown called")

	fmt.Println("error :", <-chBack)
	fmt.Println("exit")
}
