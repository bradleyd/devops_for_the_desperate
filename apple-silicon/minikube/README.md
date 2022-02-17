# Setting up minikube

Currently, minikube only works on Silicon Mac using the docker driver.

To take advantage of this, we'll use `vagrant`, `ansible`, and `parallels` to create and provisioin a VM with Docker running.

We'll use this VM to point minikube at.

## Getting Started

The first thing you need to do is make sure you have `vagrant`, `ansible`, and `minikube` installed. Please see their respective websites for instructions.

While inside the `minikube/` directory, run the following:

```
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
    default: SSH username: vagrant
    default: SSH auth method: private key
    default: Warning: Connection refused. Retrying...
    default:
    default: Vagrant insecure key detected. Vagrant will automatically replace
    default: this with a newly generated keypair for better security.
    default:
    default: Inserting generated public key within guest...
    default: Removing insecure key from the guest if it's present...
    default: Key inserted! Disconnecting and reconnecting using new SSH key...
==> default: Machine booted and ready!
==> default: Checking for Parallels Tools installed on the VM...
==> default: Setting hostname...
==> default: Configuring and enabling network interfaces...
==> default: Mounting shared folders...
    default: /vagrant => /Users/bradleydsmith/Projects/devops_for_the_desperate/apple-silicon/kubernetes
==> default: Running provisioner: ansible...
    default: Running ansible-playbook...

PLAY [dftd] ********************************************************************

TASK [Gathering Facts] *********************************************************
ok: [default]

TASK [Install aptitude using apt] **********************************************
changed: [default]

TASK [Install required system packages] ****************************************
changed: [default] => (item=apt-transport-https)
ok: [default] => (item=ca-certificates)
ok: [default] => (item=curl)
ok: [default] => (item=software-properties-common)
changed: [default] => (item=python3-pip)
changed: [default] => (item=virtualenv)
ok: [default] => (item=python3-setuptools)

TASK [Add Docker GPG apt Key] **************************************************
changed: [default]

TASK [Add Docker Repository] ***************************************************
changed: [default]

TASK [Update apt and install docker-ce] ****************************************
changed: [default]

TASK [Create a directory named 'engineering'] **********************************
changed: [default]

TASK [Copy over docker.service override] ***************************************
changed: [default]

TASK [Restart docker.service and make sure systemd reloads] ********************
changed: [default]

TASK [debug] *******************************************************************
ok: [default] => {
    "hostvars[inventory_hostname]['ansible_default_ipv4']['address']": "10.211.55.6"
}

TASK [Prints two lines of messages, but only if there is an environment value set] ***
ok: [default] => {
    "msg": [
        "To use Docker, you'll need to export DOCKER_HOST variable locally to connect to the VM",
        "Enter the following in your termninal where you are going to run your docker commands",
        "export DOCKER_HOST=tcp://10.211.55.6:2375"
    ]
}

PLAY RECAP *********************************************************************
default                    : ok=11   changed=8    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0
```

If everything went well, you should see output similar to above. Make sure you export the `DOCKER_HOST` before moving on to the next step.

## Minikube

Next, you need to start `minikube`.

Enter the following in the same terminal:

```bash
minikube start --driver=docker
😄  minikube v1.25.1 on Darwin 12.2.1 (arm64)
✨  Using the docker driver based on user configuration
👍  Starting control plane node minikube in cluster minikube
🚜  Pulling base image ...
💾  Downloading Kubernetes v1.23.1 preload ...
    > preloaded-images-k8s-v16-v1...: 417.88 MiB / 417.88 MiB  100.00% 22.02 Mi
    > gcr.io/k8s-minikube/kicbase: 343.02 MiB / 343.02 MiB  100.00% 12.86 MiB p
🔥  Creating docker container (CPUs=2, Memory=3876MB) ...
❗  Listening to 0.0.0.0 on external docker host 10.211.55.6. Please be advised
🐳  Preparing Kubernetes v1.23.1 on Docker 20.10.12 ...
    ▪ kubelet.housekeeping-interval=5m
    ▪ Generating certificates and keys ...
    ▪ Booting up control plane ...
    ▪ Configuring RBAC rules ...
🔎  Verifying Kubernetes components...
    ▪ Using image gcr.io/k8s-minikube/storage-provisioner:v5
🌟  Enabled addons: storage-provisioner, default-storageclass
🏄  Done! kubectl is now configured to use "minikube" cluster and "default" namespace by default
```

To test it out, enter the following `kubectl` command:

```bash
$ minikube kubectl -- get pods -A
NAMESPACE     NAME                               READY   STATUS    RESTARTS   AGE
kube-system   coredns-64897985d-mv9rj            1/1     Running   0          5s
kube-system   etcd-minikube                      1/1     Running   0          18s
kube-system   kube-apiserver-minikube            1/1     Running   0          18s
kube-system   kube-controller-manager-minikube   1/1     Running   0          19s
kube-system   kube-proxy-t2br6                   1/1     Running   0          6s
kube-system   kube-scheduler-minikube            1/1     Running   0          18s
kube-system   storage-provisioner                1/1     Running   0          17s
```

## Troubleshooting


### Docker build Chapter 6

In Chapter 6, we build the telnet-server docker image. We use `minikube docker-env` command to find the docker host. To test that it's working, you do a docker run command on it and use `telnet` against the minikube IP address. The minikube IP is not exposed in this setup, so you'll have to `vagrant ssh` from within the `minikube/` directory and then run the telnet command.

For example, from on the M1 Mac:

```
minikube ip
192.168.49.2
```

Then, `vagrant ssh` into the minikube VM.

```
vagrant ssh
vagrant@dftd:~$
```

Now, from within the minikube VM, you can test the running docker container:

```
telnet 192.168.49.2 2323
Trying 192.168.49.2...
Connected to 192.168.49.2.
Escape character is '^]'.

____________ ___________
|  _  \  ___|_   _|  _  \
| | | | |_    | | | | | |
| | | |  _|   | | | | | |
| |/ /| |     | | | |/ /
|___/ \_|     \_/ |___/

>q
Good Bye!
Connection closed by foreign host.
```

### Kuberenetes/Skaffold Chapter 8

In Chapter 8, we use `skaffold` to test a CI/CD pipeline. We also need to test telnet-server is running and exposed. Because of this setup, `minikube tunnel` did not work correctly for me. The easiest work around is ignore the tunnel command and using the IP address of minikube, but instead use `port-forward` in `kubectl`. 

To test telnet-server running in k8s, enter the following on your M1 mac:

```
minikube kubectl -- port-forward svc/telnet-server --address=0.0.0.0 2323:2323
```
Then in another terminal on your M1 Mac, enter the following to telnet into telnet-server:

```
telnet 127.0.0.1 2323
Trying 192.168.49.2...
Connected to 192.168.49.2.
Escape character is '^]'.

____________ ___________
|  _  \  ___|_   _|  _  \
| | | | |_    | | | | | |
| | | |  _|   | | | | | |
| |/ /| |     | | | |/ /
|___/ \_|     \_/ |___/

>q
Good Bye!
Connection closed by foreign host.
```
