package main

import (
        "context"
        "fmt"
        "strconv"
        "time"

        v1 "k8s.io/api/core/v1"
        metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

        "k8s.io/client-go/kubernetes"
        "k8s.io/client-go/rest"
)

const annotationKey = "pod.kubernetes.io/lifetime"

func main() {
        config, err := rest.InClusterConfig()
        if err != nil {
                panic(err.Error())
        }

        clientset, err := kubernetes.NewForConfig(config)
        if err != nil {
                panic(err.Error())
        }

        fmt.Println("Pod TTL Controller started...")

        for {
                pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
                if err != nil {
                        fmt.Println("Error listing pods:", err)
                        time.Sleep(10 * time.Second)
                        continue
                }

                for _, pod := range pods.Items {
                        handlePod(clientset, pod)
                }

                time.Sleep(30 * time.Second)
        }
}

func handlePod(clientset *kubernetes.Clientset, pod v1.Pod) {
        annotations := pod.Annotations
        if annotations == nil {
                return
        }

        ttlStr, exists := annotations[annotationKey]
        if !exists {
                return
        }

        ttlSeconds, err := strconv.Atoi(ttlStr)
        if err != nil {
                fmt.Println("Invalid TTL:", ttlStr)
                return
        }

        creationTime := pod.CreationTimestamp.Time
        age := time.Since(creationTime)

        if age.Seconds() > float64(ttlSeconds) {
                fmt.Printf("Deleting pod %s/%s (age: %v)\n",
                        pod.Namespace, pod.Name, age)

                err := clientset.CoreV1().Pods(pod.Namespace).Delete(
                        context.TODO(),
                        pod.Name,
                        metav1.DeleteOptions{},
                )

                if err != nil {
                        fmt.Println("Delete error:", err)
                }
        }
}
