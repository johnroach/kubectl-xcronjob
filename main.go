package main

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	//
	// Uncomment to load all auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth"
	//
	// Or uncomment to load specific auth plugins
	// _ "k8s.io/client-go/plugin/pkg/client/auth/azure"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	// _ "k8s.io/client-go/plugin/pkg/client/auth/openstack"
)

func main() {
	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	for {
		pods, err := clientset.CoreV1().Pods("cronjob").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Printf("There are %d pods in the cronjob namespace\n", len(pods.Items))

		jobs, err := clientset.BatchV1().Jobs("cronjob").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		cronjobs, err := clientset.BatchV1beta1().CronJobs("cronjob").List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			panic(err.Error())
		}
		fmt.Print("-- CRONJOBS --\n")
		for _, cronjob := range cronjobs.Items {
			fmt.Printf("%v\n", cronjob.Name)
			fmt.Printf("\tCronjob UID:\t%v\n", cronjob.ObjectMeta.UID)
		}

		fmt.Print("-- JOBS --\n")
		for _, job := range jobs.Items {
			fmt.Printf("%v\n", job.Name)
			fmt.Printf("\tJob UID:\t%v\n", job.UID)
			fmt.Printf("\tOwner-ref:\t%v\n", job.OwnerReferences[0].UID)
		}

		fmt.Print("-- PODS --\n")
		for _, pod := range pods.Items {
			fmt.Printf("%v\n", pod.ObjectMeta.Name)
			fmt.Printf("\tStatus:\t%v\n", pod.Status.Phase)
			fmt.Printf("\tOwner-ref:\t%v\n", pod.OwnerReferences[0].UID)
		}

		time.Sleep(10 * time.Second)
	}
}
