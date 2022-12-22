package containerd

import (
	"fmt"
	"github.com/Interstellarss/lighthouse-plugin/pkg/plugin"
	"github.com/Interstellarss/lighthouse-plugin/pkg/plugin/util"
	"github.com/containerd/containerd"
	"github.com/docker/docker/api/types/container"
	"github.com/gorilla/mux"
	"golang.org/x/net/context"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/klog"
	"net/http"
	"strings"
	"sync"
)

func init() {
	plugin.AllPlugins[ContainerdPreHookUpdatePluginName] = PreHookUpdateContainer
}

const (
	ContainerdPreHookUpdatePluginName = "ContainerdPreUpdate"
)

type preHookContainerdUpdatePlugin struct {
	plugin.BasePlugin
	startOnce sync.Once
	handler   InterceptUpdateFunc
	toolkits  *ToolKits
}

type InterceptUpdateFunc func(
	toolKits *ToolKits,
	updateConfig *container.UpdateConfig,
	metadata *plugin.PodMetadata,
	containerStatus string) error

var PreHookUpdateContainer = &preHookContainerdUpdatePlugin{
	BasePlugin: plugin.BasePlugin{
		Method: http.MethodPost,
		Path:   "/containerd/{name:.*}/update",
	},
}

//SetIgnored set ignored value
func (p *preHookContainerdUpdatePlugin) SetIgnored(ignored plugin.IgnoreNamespacesFunc) {
	p.Ignored = ignored
}

//Path return URL path
func (p *preHookContainerdUpdatePlugin) Path() string {
	return fmt.Sprintf("/prehook%s", p.BasePlugin.Path)
}

//Method return method name
func (p *preHookContainerdUpdatePlugin) Method() string {
	return p.BasePlugin.Method
}

//Handler return handler function, which accept hook request and send to all handlers
func (p *preHookContainerdUpdatePlugin) Handler(
	k8sclient kubernetes.Interface, containerdEndpoint, containerdVersion string) http.HandlerFunc {
	//initialzation just do once
	p.startOnce.Do(func() {
		client, err := containerd.New(address)

		if err != nil {
			klog.Exitf("can not create containerd client, %v", err)
			return
		}

		//initializa tool kit
		p.toolkits = &ToolKits{
			K8sClient:        k8sclient,
			ContainerdClient: client,
		}
	})

	return func(w http.ResponseWriter, r *http.Request) {
		//accepting hook request
		klog.V(5).Infof("Plugin handle request %s", r.URL.String())

		vars := mux.Vars(r)

		if vars == nil {
			klog.Errorf("can not get vars in request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		bodyBytes, err := ioutil.ReadAll(r.Body)

		if err != nil {
			klog.Errorf("can not read request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		updateConfig := &container.UpdateConfig{}
		if err := util.Json.Unmarshal(bodyBytes, updateConfig); err != nil {
			klog.Errorf("hook error, %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		name := vars["name"]
		klog.V(5).Infof("Update %s container", name)

		//filter invalide request
		metadata, status, err := p.generateContainerData(name, w)
		if err != nil {

		}

		//send containerd(CRI) update request to all handlers
		if err := p.handler(p.toolkits, updateConfig, metadata, status); err != nil {
			klog.Errorf("hook error, %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		//group final response data
		patchBytes, err := groupPatchData(updateConfig, bodyBytes)

		if err != nil {
			klog.Errorf("merge patch err: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		klog.V(5).Infof("Containerd patch: %s", string(patchBytes))

		util.Json.NewEncoder(w).Encode(&PatchData{
			PatchType: string(types.MergePatchType),
			PatchData: patchBytes,
		})

	}
}

func (p *preHookContainerdUpdatePlugin) generateContainerData(conName string, w http.ResponseWriter) (
	metadata *plugin.PodMetadata, status string, err error) {

	//maybe this is the correct use of load container infor from containerd?
	containerInfo, err := p.toolkits.ContainerdClient.LoadContainer(context.Background(), conName)
	if err != nil {
		if !strings.Contains(err.Error(), "No such container") {
			klog.Errorf("can not container info from containerd %s, %v", conName, err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write(noChangeResponse)
		return
	}
	var Clabels map[string]string
	if Clabels, err = containerInfo.Labels(context.Background()); Clabels == nil {
		if err != nil {
			return
		}
		klog.Warningf("Empty labels")
		w.WriteHeader(http.StatusOK)
		w.Write(noChangeResponse)
		err = fmt.Errorf("empty lables")
		return
	}

	//containerInfo.

	klog.V(5).Infof("Get pod from container labels")

	metadata = plugin.GetPodMetadata(Clabels)

	if p.Ignored(metadata.Namespace) {
		klog.V(5).Infof("Ignored namespace %s", metadata.Namespace)

		w.WriteHeader(http.StatusOK)
		w.Write(noChangeResponse)
		err = fmt.Errorf("ignred namespace %s", metadata.Namespace)
		return
	}

	//status = containerInfo.Task()
	//TODO: get status of a running container?

	//TODO: should return the actual status of the container
	return metadata, "", err

}

//RegisterSubPlugin export register interface
func (p *preHookContainerdUpdatePlugin) RegisterSubPlugin(h InterceptUpdateFunc) {
	p.handler = h
}
