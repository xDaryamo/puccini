Puccini Quickstart
==================

[Download and install Puccini](https://github.com/tliron/puccini/releases).

If you don't have access to a running Kubernetes cluster, an easy way to get one is via
[Minikube](https://github.com/kubernetes/minikube). You're also going to need **kubectl**. Here's
a quick script to get them both:

    cd /tmp
    wget -O kubectl https://storage.googleapis.com/kubernetes-release/release/v1.13.0/bin/linux/amd64/kubectl
    wget -O minikube https://storage.googleapis.com/minikube/releases/v0.30.0/minikube-linux-amd64
    chmod +x kubectl minikube
    sudo mv kubectl minikube /usr/bin/

Start a Minikube virtual machine with enough memory for our demo application:

    minikube start --memory=4096

Give it some time to start, and then access the dashboard from your web browser:

    minikube dashboard

Now compile and apply the demo application's TOSCA:

    puccini-tosca compile examples/kubernetes/bookinfo/bookinfo-simple.yaml | puccini-js exec kubernetes.generate | kubectl apply -f -

On the dashboard you'll see the pods coming up. When they're finally up, forward a port from the
frontend pod so that we can access it via your browser.

    POD=$(kubectl get pods -l service=productpage -o jsonpath='{.items[0].metadata.name}')
    kubectl port-forward $POD 9080:9080 &

Now you can see the application at [http://localhost:9080](http://localhost:9080).

When you're done, you can stop the port forwarding and destroy the Minikube:

    killall kubectl
    minikube delete

The next step would to be look at the [examples](examples/) and learn more about what
Puccini can do.
