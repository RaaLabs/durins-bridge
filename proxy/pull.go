package proxy

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"golang.org/x/exp/slog"
)

type PullOverride struct {
	Client *client.Client
}

func (po *PullOverride) ServeHTTP(w http.ResponseWriter, r *http.Request) bool {
	if r.Method != http.MethodPost || !strings.HasPrefix(r.URL.Path, "/images/create") {
		return false
	}

	imageRef := r.URL.Query().Get("fromImage")
	tag := r.URL.Query().Get("tag")

	if imageRef == "" {
		return false
	}

	if tag != "" {
		imageRef = fmt.Sprintf("%s:%s", imageRef, tag)
	}

	slog.Debug("Got image pull request for", "tag", imageRef)
	images, err := po.Client.ImageList(
		r.Context(),
		types.ImageListOptions{
			All: true,
			Filters: filters.NewArgs(filters.Arg(
				"reference",
				imageRef,
			)),
		},
	)
	if err != nil {
		slog.Error("Failed to list Docker images", err)
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

	slog.Info("Image not found locally, not handling request", "tag", imageRef)
	return false

ignore:
	slog.Info("Image was found locally, will complete the request without actually pulling", "tag", imageRef)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(
		"{\"status\":\"Will not pull %s - it is already present locally\"}",
		imageRef,
	)))

	return true
}
