package proxy

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type PullOverride struct {
	Client *client.Client
}

func (po *PullOverride) ServeHTTP(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost || !strings.HasPrefix(r.URL.Path, "/images/create") {
		return false
	}

	if !r.URL.Query().Has("fromImage") {
		return false
	}

	if !r.URL.Query().Has("tag") {
		// TODO: If pulling without tag, Docker will pull all tags - we probably don't want that?
		http.Error(w, "pulling all tags for a repository is not allowed", http.StatusForbidden)
		return true
	}

	imageRef := fmt.Sprintf(
		"%s:%s",
		r.URL.Query().Get("fromImage"),
		r.URL.Query().Get("tag"),
	)

	images, err := po.Client.ImageList(
		r.Context(),
		types.ImageListOptions{
			All:     true,
			Filters: filters.NewArgs(filters.Arg("reference", imageRef)),
		},
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to list images: %s", err), http.StatusBadGateway)
		return true
	}

	for _, image := range images {
		for _, digest := range image.RepoTags {
			if imageRef == digest {
				goto ignore
			}
		}
	}

	log.Println("Image not found locally, letting it through")
	return false

ignore:
	log.Println("Image", imageRef, "was already found locally, will ignore pull request.")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(
		"{\"status\":\"Will not pull %s - it is already present locally\"}",
		imageRef,
	)))

	return true
}
