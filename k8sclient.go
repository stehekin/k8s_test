package main

// import "k8s.io/client-go/kubernetes"
import (
	"context"
	"fmt"
	"path/filepath"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/local"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	v1 "k8s.io/cri-api/pkg/apis/runtime/v1"
)

const containerdSockPath = "unix:///var/run/containerd/containerd.sock"

func main() {
	// cri_containerd()
	cri_containerd_events()
}

func kube() {
	config, err := config(true)
	if err != nil {
		panic(err.Error())
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	pods, err := client.CoreV1().Pods("default").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, pod := range pods.Items {
		fmt.Printf("\n\n\n%v \n", pod.Name)
		for _, container := range pod.Status.ContainerStatuses {
			fmt.Printf("\t\t %v\n", container.ContainerID)
		}
		for _, container := range pod.Status.InitContainerStatuses {
			fmt.Printf("\t\t init %v\n", container.ContainerID)
		}

	}
}

func config(inCluster bool) (*rest.Config, error) {
	if inCluster {
		return rest.InClusterConfig()
	} else {
		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
}

func cri_containerd() {
	conn, err := grpc.NewClient(containerdSockPath, grpc.WithTransportCredentials(local.NewCredentials()))
	if err != nil {
		panic(err.Error())
	}

	client := v1.NewRuntimeServiceClient(conn)

	request := v1.ListContainersRequest{}

	response, err := client.ListContainers(context.TODO(), &request)
	if err != nil {
		panic(err.Error())
	}

	for _, c := range response.Containers {
		fmt.Printf("%v\n", c)
	}
}

func cri_containerd_events() {
	conn, err := grpc.NewClient(containerdSockPath, grpc.WithTransportCredentials(local.NewCredentials()))
	if err != nil {
		panic(err.Error())
	}

	client := v1.NewRuntimeServiceClient(conn)
	request := v1.GetEventsRequest{}

	eventClient, err := client.GetContainerEvents(context.TODO(), &request)
	if err != nil {
		panic(err.Error())
	}

	for {
		response, err := eventClient.Recv()
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("%v\n", response)
	}

}

func cri_containerd_pid() {
	conn, err := grpc.NewClient(containerdSockPath, grpc.WithTransportCredentials(local.NewCredentials()))
	if err != nil {
		panic(err.Error())
	}

	client := v1.NewRuntimeServiceClient(conn)

	request := v1.ContainerStatusRequest{
		ContainerId: "720cc29a8cc98",
		Verbose:     true,
	}

	stats, err := client.ContainerStatus(context.TODO(), &request)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("%v\n", stats)
}

func cri_containerd_image() {
	conn, err := grpc.NewClient(containerdSockPath, grpc.WithTransportCredentials(local.NewCredentials()))
	if err != nil {
		panic(err.Error())
	}
	iClient := v1.NewImageServiceClient(conn)
	request := v1.ListImagesRequest{}
	response, err := iClient.ListImages(context.TODO(), &request)
	if err != nil {
		panic(err.Error())
	}

	for _, i := range response.Images {
		// You can find image URI here.
		fmt.Printf("%v\n", i)
	}
}
