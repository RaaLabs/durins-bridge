package listening

import (
	"log"
	"net"
	"os"
)

func CreateListenerFromArgument(path string) (net.Listener, error) {
	info, err := os.Stat(path)

	if err != nil && os.IsNotExist(err) {
		goto create
	} else if err != nil {
		return nil, err
	}

	if info.Mode()&os.ModeSocket != 0 {
		log.Println(path, "already exists and is a socket, will try to recreate")
		if err := os.Remove(path); err != nil {
			return nil, err
		}
	}

create:
	return net.Listen("unix", path)
}
