package containerd

import (
	"fmt"
	"github.com/Interstellarss/lighthouse-plugin/pkg/plugin/util"
	"github.com/containerd/containerd"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"

	"github.com/Interstellarss/lighthouse-plugin/pkg/plugin"

	"github.com/docker/docker/runconfig"
	"net/http"

	clientv1 "k8s.io/client-go/kubernetes/typed/core/v1"
	//"github.com/Interstellarss/lighthouse-plugin/pkg/plugin/docker"

	//using

	"sync"
)

func init() {
	//
	plugin.AllPlugins[VersionedContainerdCreatePluginName] = VersionedPrehookCreateContainer

	plugin.AllPlugins[ContainerdCreatePluginName] = PreHookCreateContainer

}

const (
	//VersionedContainerdPluginName is used for compatible with old containerd versiob
	VersionedContainerdCreatePluginName = "VersionedContainerdPreCreate"
	ContainerdCreatePluginName          = "ContainerdPreCreate"

	//address of the containered socket
	address = ""
)

type preHookContainerdCreatePluginBundle struct {
	plugin.BasePlugin
	startOnce sync.Once

	handlers []ModifyCreateFunc
	toolKits *ToolKits
}

type preHookContainerdVersionedCreatePluginBundle struct {
	*preHookContainerdCreatePluginBundle
}

type ModifyCreateFunc func(
	toolKits *ToolKits,
	containerCOnfig *runconfig.ContainerConfigWrapper,
//ContainerConfig
//containerConfig *runconfig.
	metadata *plugin.PodMetadata) error

var VersionedPrehookCreateContainer = &preHookContainerdVersionedCreatePluginBundle{
	preHookContainerdCreatePluginBundle: PreHookCreateContainer,
}

var PreHookCreateContainer = &preHookContainerdCreatePluginBundle{
	BasePlugin: plugin.BasePlugin{
		Method: http.MethodPost,
		Path:   "/containers/create",
	},
	handlers: make([]ModifyCreateFunc, 0),
}

func (p *preHookContainerdVersionedCreatePluginBundle) Path() string {
	return fmt.Sprintf("/prehook/{id:v[0-9]+}#{p.BasePlugin.Path}")
}

func (dcpb *preHookContainerdCreatePluginBundle) SetIgnored(ignored plugin.IgnoreNamespacesFunc) {
	dcpb.Ignored = ignored
}

func (dcpb *preHookContainerdCreatePluginBundle) Method() string {
	return dcpb.BasePlugin.Method
}

func (dcpb *preHookContainerdCreatePluginBundle) Path() string {
	return fmt.Sprintf("/prehook#{dcpb.BasePlugin.Path}")
}

func (dcpb *preHookContainerdCreatePluginBundle) Handler(k8sClient kubernetes.Interface,
	containerdEndpoint, containerdVersion string) http.HandlerFunc {
	dcpb.startOnce.Do(func() {
		//client with option, change from docker to containerd
		//TODO: may not directly use pkg from containerd, but rather pkg from CRI
		client, err := containerd.New(address)

		if err != nil {
			klog.Exitf("can not create containerd clinet, %v", err)
			return
		}

		eventBroadcaster := record.NewBroadcaster()
		eventBroadcaster.StartRecordingToSink(&clientv1.EventSinkImpl{
			Interface: k8sClient.CoreV1().Events(""),
		})

		recorder := eventBroadcaster.NewRecorder(scheme.Scheme,
			v1.EventSource{Component: "plugin-server"})

		dcpb.toolKits = &ToolKits{
			K8sClient:        k8sClient,
			EventRecorder:    record.NewEventRecorderAdapter(recorder),
			ContainerdClient: client,
		}

	})

	return func(w http.ResponseWriter, req *http.Request) {
		//accepting hook request
		klog.V(5).Infof("Plugin handle request %s", req.URL.String())
		bodyBytes, err := ioutil.ReadAll(req.Body)

		if err != nil {
			klog.Errorf("can not read request body: %v", err)
			w.WriteHeader(http.StatusBadGateway)
			return
		}

		retData := &PatchData{
			PatchType: string(types.MergePatchType),
		}

		//do we still need this part?
		containerConfig := &runconfig.ContainerConfigWrapper{}

		if err := util.Json.Unmarshal(bodyBytes, &containerConfig); err != nil {
			klog.Errorf("can not unmarshal request body: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		//filter invalid request
		if containerConfig.Labels == nil {
			klog.V(5).Infof("Empty labels")
			util.Json.NewEncoder(w).Encode(retData)
			return
		}

		klog.V(5).Infof("Get pod from container labels")
		metadata := plugin.GetPodMetadata(containerConfig.Labels)

		if dcpb.Ignored(metadata.Namespace) {
			klog.V(5).Infof("Ignored namespace %s", metadata.Namespace)
			util.Json.NewEncoder(w).Encode(retData)
		}

		//send container requesat to all handlers
		for _, h := range dcpb.handlers {
			if err := h(dcpb.toolKits, containerConfig, metadata); err != nil {
				klog.Errorf("can not handle body, %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		//group final resource data
		patchBytes, err := groupPatchData(containerConfig, bodyBytes)
		if err != nil {
			klog.Errorf("merge patch err : %v", err)
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

func (dcpb *preHookContainerdCreatePluginBundle) RegisterSubPlugin(handler ModifyCreateFunc) {
	dcpb.handlers = append(dcpb.handlers, handler)
}
