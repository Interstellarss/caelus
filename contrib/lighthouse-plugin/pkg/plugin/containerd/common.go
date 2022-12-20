package containerd

import (
	"github.com/Interstellarss/lighthouse-plugin/pkg/plugin/util"
	"github.com/containerd/containerd"

	//"github.com/containerd/containerd/client"
	jsonpatch "github.com/evanphx/json-patch"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/events"

	"encoding/json"
)

var (
	noChangeResponse = []byte(`{}`)
)

type ToolKits struct {
	K8sClient     kubernetes.Interface
	EventRecorder events.EventRecorder

	//using grpc for CRI instead
	//ContainerdClient *containerd.Client
	ContainerdClient *containerd.Client
}

type PatchData struct {
	PatchType string `json:"patchType,omitempty"`
	PatchData []byte
}

type PostHookData struct {
	StatusCode int             `json:"statusCode,omitempty"`
	Body       json.RawMessage `json:"body,omitempty"`
}

func groupPatchData(config interface{}, bodyBytes []byte) ([]byte, error) {
	newBodyBytes, err := util.Json.Marshal(config)
	if err != nil {

	}
	//containerd.

	patchBytes, err := jsonpatch.CreateMergePatch(bodyBytes, newBodyBytes)

	if err != nil {

	}
	return patchBytes, nil
}
