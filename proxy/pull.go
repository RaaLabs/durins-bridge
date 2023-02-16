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

	if !r.URL.Query().Has("fromImage") {
		return false
	}

	if !r.URL.Query().Has("tag") {
		slog.Warn("Not allowing pull request without tag, it would pull all the tags for the givene repository", "image", r.URL.Query().Get("fromImage"))
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
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(
		"{\"status\":\"Will not pull %s - it is already present locally\"}",
		imageRef,
	)))

	return true
}
