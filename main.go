package main

import (
	"raalabs.tech/durins-bridge/cmd"
)

func main() {
	cmd.Execute()
	// listeners, err := activation.Listeners()
	// if err != nil {
	// 	log.Fatalln("Failed to activate listeners", err)
	// }

	// if len(listeners) != 1 {
	// 	log.Fatalln("Expected exactly one file descriptor, got", len(listeners))
	// }

	// proxy := httputil.ReverseProxy{
	// 	Director: func(r *http.Request) {
	// 		log.Println("HELLO", r.URL.String())
	// 	},
	// }

	// server := http.Server{
	// 	Handler: &proxy,
	// }

	// // listener := noCloseListener{listeners[0]}
	// listener := listeners[0]
	// err = server.Serve(listener)
	// if err != nil {
	// 	log.Fatalln("Listening failed", err)
	// }
}
