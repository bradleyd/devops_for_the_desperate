# QEMU

QEMU https://www.qemu.org/ is an open-source emulator. It can mimic different architectures, like the Apple Silicon processor. QEMU is free and can be installed on your M1 Mac by using the arm version of brew.

```bash
brew install qemu
```

Once it is installed, you can use a Debian based machine that I have created to get you going.

It is located [here](https://dftd-qemu.sfo3.digitaloceanspaces.com/dftd-aarch64-qcow2.tar.gz)

> __The file is around 1 GB, so it may take some time to download.__

Go ahead and download the tar file to your host. Once there go ahead and unzip it, and you should have three files:

```bash
tar -xzf dftd-aarch64-qcow2.tar.gz
```

The archive contains three files:

```bash
debian-dftd-aarch64.qcow2
initrd.img-5.10.0-11-arm64
vmlinuz-5.10.0-11-arm64
```

We'll use these three files to run a Debian host in QEMU.

> Please note, this prebuilt image has the Docker daemon already installed.
> Please use this host to follow along with Chapter 6-9 in the book.

Before you launch the image, here are some important details you'll need to log in and connect.

* username `dftd`
* password `dftd`
* root password `dftd`
* Docker port `20375`
* SSH port `10022`
* Greeting Service port `5001`

> There is also `vagrant` user and `ubuntu` user to match the book's examples. The passwords for these users are `vagrant` and `ubuntu` respectively.
> Please adjust the ports if they conflict with existing services.

## Starting the Host

To start the QEMU host, enter the following command in a terminal where you extracted the files from above:

```bash
qemu-system-aarch64 -cpu cortex-a72 -smp cpus=4,sockets=1,cores=4,threads=1 -machine virt,highmem=off -accel hvf -accel tcg,tb-size=2048 -m 4G -initrd initrd.img-5.10.0-11-arm64 -kernel vmlinuz-5.10.0-11-arm64 -append "root=/dev/vda2 console=ttyAMA0" -drive if=virtio,file=debian-dftd-aarch64.qcow2,format=qcow2,id=hd -net user,hostfwd=tcp::10022-:22,hostfwd=tcp::20375-:2375,hostfwd=tcp::5001-:5000 -net nic -device intel-hda -device hda-duplex -nographic
```

> Note, the `cpus` are set to 4 and the `memory` is set to 4G. Please adjust to your taste.

In a few seconds you should be able to log into the Debian host using the credentials and port above like the following:

```bash
ssh dftd@127.0.0.1 -p 10022
```

## Pre-Provision

You need to configure the QEMU host for the book. I have already included stuff like Docker in the prebuilt image, but things like Golang and users need to be added.

In Chapter 1 of the book we use Vagrant to create our infrastructure which in turn launches Ansible to provision our host.

Since we have already created our infrastructure, we'll launch the Ansible part from the terminal. Please make sure Ansible is installed on your local host and run the following command to set up the QEMU host for the books examples:

```bash
ansible-playbook ./ansible/site.yml -i hosts -c paramiko --ask-pass --ask-become-pass -u dftd
SSH password: <dftd>
BECOME password[defaults to SSH password]: <dftd>
...
PLAY RECAP ********************************************************************************************************************************************************************************************************************************************************************
127.0.0.1                  : ok=8    changed=5    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0
```

Now that is done, you are ready to follow along!

## Chapter 1-5

In Chapters 1-5, every time you see `vagrant provision`, you are going to run `ansible-playbook` command instead.

In Chapter 1 and Chapter 2, you can run the following command in lieu of `vagrant provision`

```bash
ansible-playbook ../../ansible/site.yml -i hosts -c paramiko --ask-pass --ask-become-pass -u ubuntu
SSH password:
BECOME password[defaults to SSH password]:

PLAY [Provision VM] **************************************************************************************************************************************
```

Enter `ubuntu` when prompted for the SSH password and BECOME password.

> NOTE: This is still using the `ansible/` playbook in the main directory of this repository. The local `ansible/` relative to this README, is only to set up the host before Chapter 1.

Once you get to Chapters 3,4, and 5, you'll need to adjust `ansible-playbook` command by appending the private-key flag. You'll create the public key pair in Chapter 3 and will need it to connect once the Chapter 3 tasks are run as they disable passwords for SSH login.

```bash
ansible-playbook ../../ansible/site.yml -i hosts -c paramiko --ask-pass --ask-become-pass -u ubuntu --private-key ~/.ssh/dftd
```

In Chapter 4 and Chapter 5, the book's examples have you install the Greeting web service which listens on port 5000. On my M1 Mac, port 5000 is being used by the Airplay receiver and because of this I set the QEMU host to listen on port 5001. Any request to 5001 on the Silicon Mac will be forwarded to port 5000 inside the QEMU host.

In chapter 5, if the example calls for testing against port 5000, substitute port 5001 instead.

In Chapter 5, the book instructs you to use the additional interface the Vagrant and VirtualBox creates for `nmap`.
We do not have this for our QEMU host, so just run `nmap` against 127.0.0.1.

## Chapters 6-9

These chapters will have you install `minikube` to provide a Kubernetes cluster and docker registry.

### Minikube

Do not install minikube on your Silicon Mac, instead you'll install it in the VM you just created via vagrant.

To do this, you need to ssh into the vm. From a terminal, enter the following:

`ssh -i ~/.ssh/dftd ubuntu@127.0.0.1 -p 10022`

Now that you are in the QEMU host, enter the following commands to install the ARM version of minikube.

```bash
curl -LO https://storage.googleapis.com/minikube/releases/latest/minikube-linux-arm64
sudo install minikube-linux-arm64 /usr/local/bin/minikube
```

Next, you can follow the book's example and start minikube. The only difference is you need to choose the docker driver instead of virtualbox like below:

```bash
minikube start --driver=docker
😄  minikube v1.25.2 on Ubuntu 20.04 (arm64)
✨  Using the docker driver based on user configuration
👍  Starting control plane node minikube in cluster minikube
🚜  Pulling base image ...
...
```


__You should be able to build the `telnet-server` image as instructed in Chapter 6 now!__

### Setting Up Your Pipeline

In Chapter 8, you'll install `skaffold` for setting up our CI/CD pipeline.
Please make sure you choose the Linux arm64 version since you are running in arm64 architecture.

```bash
curl -Lo skaffold https://storage.googleapis.com/skaffold/releases/latest/skaffold-linux-arm64 && \
sudo install skaffold /usr/local/bin/
```

You'll also need to install the `container-structure-test` application. Please choose the arm64 version as well.

```bash
curl -LO https://storage.googleapis.com/container-structure-test/v1.11.0/container-structure-test-linux-arm64 && \
mv container-structure-test-linux-arm64 container-structure-test && chmod +x container-structure-test && sudo mv container-structure-test /usr/local/bin/
```

You should be all set on your Silicon Mac!

## Troubleshooting

> If something does not work, please open an issue here in this repository and I will do my best to respond and find a fix!
