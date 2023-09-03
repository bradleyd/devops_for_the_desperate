# Apple Silicon

The easiest way to follow along with the book is to use `Parallels Pro`. This is a paid software solution, but it will work with `vagrant` and `minikube`.

> If you do not want to pay for Parallels, you can do something similar with QEMU but will require more setup

## Getting Started

The first thing you need to do is make sure you have `vagrant`, `ansible` installed. From a terminal, enter the following:

```bash
brew install ansible vagrant
```

Then, you'll need to install the `vagrant-parallels` plugin.

```bash
vagrant plugin install vagrant-parallels
```

Once that is installed, make sure you are in `apple-silicon/vagrant/` directory before running the `vagrant up` command like we do in Chapter 1:

> Make sure you add `--provider=parallels` to the command

```vagrant
$ vagrant up --provider=parallels
Bringing machine 'default' up with 'parallels' provider...
==> default: Registering VM image from the base box 'bento/ubuntu-20.04-arm64'...
==> default: Creating new virtual machine as a linked clone of the box image...
==> default: Unregistering the box VM image...
==> default: Setting the default configuration for VM...
==> default: Checking if box 'bento/ubuntu-20.04-arm64' version '202112.19.0' is up to date...
==> default: Setting the name of the VM: apple_silicon_default_1645040335644_35727
==> default: Preparing network interfaces based on configuration...
    default: Adapter 0: shared
    default: Adapter 1: hostonly
==> default: Clearing any previously set network interfaces...
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
    default: /vagrant => /Users/bradleydsmith/Projects/devops_for_the_desperate/vagrant/apple_silicon
==> default: Running provisioner: ansible...
    default: Running ansible-playbook...

PLAY [Provision VM] ************************************************************

TASK [Gathering Facts] *********************************************************
ok: [default]

PLAY RECAP *********************************************************************
default                    : ok=1    changed=0    unreachable=0    failed=0    skipped=0    rescued=0    ignored=0
```

__Please note, Vagrant and Parallels use the default SSH port (22) and the IP address of the VM instead of `localhost` and port 2222.__

To get the IP of the VM you can enter the following in a terminal:

```
vagrant ssh-config
Host default
  HostName 10.211.55.20
  User vagrant
  Port 22
  UserKnownHostsFile /dev/null
  StrictHostKeyChecking no
  PasswordAuthentication no
  IdentityFile /Users/bradleydsmith/Projects/devops_for_the_desperate/apple-silicon/vagrant/.vagrant/machines/default/parallels/private_key
  IdentitiesOnly yes
  LogLevel FATAL
```

Grab `HostName` IP, 10.211.55.20 in my case, and use that in any of the books text where it ask you to ssh in.
You can also use this IP when the book references accessing the public IP of the VM too.

```
ssh -i ~/.ssh/dftd bender@<VM IP> 

```

You are now ready to follow along on the book as you were on a Linux or Intel-based Mac.
