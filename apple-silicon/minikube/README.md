# Setting up minikube

Currently, `minikube` only works on Silicon Mac using the [docker](https://github.com/kubernetes/minikube/issues/11219) driver.

To take advantage of this, we'll use `vagrant`, `ansible`, and `parallels` to create and provision a VM with Docker running.

> __We'll use this VM to follow the examples in Chapter 6,7,8,9.__

## Getting Started

The first thing you need to do is make sure you have `vagrant`, `ansible`, and `parallels` installed. Please see their respective websites for instructions.

While inside this `minikube/` directory, enter the following command to create and provision our VM:

```bash
vagrant up --provider=parallels
Bringing machine 'default' up with 'parallels' provider...
==> default: Registering VM image from the base box 'bento/ubuntu-20.04-arm64'...
==> default: Creating new virtual machine as a linked clone of the box image...
==> default: Unregistering the box VM image...
==> default: Setting the default configuration for VM...
==> default: Checking if box 'bento/ubuntu-20.04-arm64' version '202112.19.0' is up to date...
==> default: Setting the name of the VM: kubernetes_default_1645049630447_62496
==> default: Preparing network interfaces based on configuration...
    default: Adapter 0: shared
    default: Adapter 1: hostonly
==> default: Clearing any previously set network interfaces...
==> default: Running 'pre-boot' VM customizations...
==> default: Booting VM...
==> default: Waiting for machine to boot. This may take a few minutes...
    default: SSH address: :22
...

PLAY RECAP *********************************************************************
default                    : ok=11   changed=8    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0
```

If everything went well, you should see output similar to above. I leveraged Ansible to provision this VM with the following to make it easier for you:

* docker
* Golang

Feel free to check out the playbook in `ansible/site.yml` for more details.

### Installing minikube

Do not install minikube on your Silicon Mac, instead you'll install it in the VM you just created via vagrant.

> You can install minikube on your Silicon Mac, but you will end up having to proxy all Kubernetes service requests and that can get cumbersome and confusing. Feel free to install it on your local host if you are comfortable with the differences.

To install minikube, you need to ssh into the vm. From a terminal, enter the following:

`vagrant ssh`

Now that you are in the VM, enter the following commands to install the ARM version of minikube.

```bash
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-arm64
sudo install minikube-linux-arm64 /usr/local/bin/minikube
```

Next, you should be able to follow the book as written and start minikube. The only difference is you need to choose the docker driver instead of virtaulbox like below:

```bash
minikube start --driver=docker
😄  minikube v1.25.2 on Ubuntu 20.04 (arm64)
✨  Using the docker driver based on user configuration
👍  Starting control plane node minikube in cluster minikube
🚜  Pulling base image ...
...
```

With minikube installed, you should be able to perform the examples as written! However, there are some caveats where the book differs from the actual implementation on your Silicon Mac. You should read the [troubleshooting](#troubleshooting) section below before following Chapters 6-9.

## Troubleshooting

> If something does not work, please open an issue here in this repository and I will do my best to respond and find a fix!

Chapters 6 and 7 should work as written, just make sure you are performing the steps inside this VM.

### Chapter 8

In Chapter 8, you'll install build out a simple CI/CD pipeline. Please make sure you choose the Linux arm64 versions for all the software required since you are running them in an arm64 host. Here is a quick reference to the versions you want to install for both `skaffold` and `container-structure-test`.

#### skaffold

```bash
curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-arm64 && \
sudo install skaffold /usr/local/bin/
```

#### container-structure-test

```bash
curl -LO https://storage.googleapis.com/container-structure-test/v1.11.0/container-structure-test-linux-arm64 && \
mv container-structure-test-linux-arm64 container-structure-test && chmod +x container-structure-test && sudo mv container-structure-test /usr/local/bin/
```

### Chapter 9

In this Chapter, you'll set up the monitoring stack. Before you follow those instructions to install it, you need to build the connection simulator `bbs-warrior` docker image in the VM. The one provided in the manifest was built for non-arm CPU's and hosted on GitHub.

#### bbs-warrior

Navigate to the `monitoring/bbs-warrior/` directory and follow the README instructions on how to build the docker image. Don't forget to make sure you have run `eval $(minikube -p minikube docker-env)` command in the same terminal as you are working in from inside the VM. Here is a quick command reference:

```bash
docker build -t ghcr.io/bradleyd/devops_for_the_desperate/bbs-warrior:latest .
```

Once that is successfully built, you can follow along with Chapter 9 to install the stack as written.

#### Services in the Browser

The next hurdle you will encounter is you will be asked to view Grafana and Prometheus dashboards in your browser via the `minikube service` command. Since we are doing this in a VM without a desktop, you'll have to port forward these services via SSH on your local host outside the VM.

To do this you'll need to find out the IP address that Vagrant uses to SSH into. You'll use this in the SSH port forward command later on. Enter the following command to get the IP address:

```bash
vagrant ssh-config | grep HostName
HostName 10.211.55.12
```

The IP address assigned to HostName is the one we want. Remember this address as we'll need it next.

Next, we need to locate the NodePort port that minikube assigned to Prometheus, Grafana, and Alertmanager. These ports will allow us to forward traffic from our local host's browser into the VM for those services.

Use the following command inside the VM to view the dynamic ports for each service.

```bash
minikube -n monitoring service --all
|----------------------|---------------------------|
|         NAME         |            URL            |
|----------------------|---------------------------|
| alertmanager-service | http://192.168.49.2:30717 |
| grafana-service      | http://192.168.49.2:31797 |
| prometheus-service   | http://192.168.49.2:31234 |
|----------------------|---------------------------|
```

I have shortened the output for readability. But the data you want to pay attention to is the URL column. In particular the IP address and port.

For each of these 3 ports you'll need to run the following command in a terminal on your local host. For example, to see the Grafana dashboard on my local host's browser enter the following:

```bash
ssh -L 31797:192.168.49.2:31797 vagrant@10.211.55.12
```

It will prompt you for vagrant's password which is just `vagrant`.

This command states, forward any traffic on your local machine (127.0.0.1) destined for port 31797, to IP address 192.168.49.2 (minikube IP inside the MV) over the SSH connection--vagrant@10.211.55.12.

Once you are logged in, leave the SSH terminal open and visit `http://127.0.0.1:31797` in your browser. You should now see the Grafana dashboard!

Repeat the above SSH forwarding step for each of the other two services in a different terminal on your local host. Please note, leave these services open in your browser and the SSH terminal for the entirety of Chapter 9 as you'll be referencing them more than once.
