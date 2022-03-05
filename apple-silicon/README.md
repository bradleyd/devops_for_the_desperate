# Apple Silicon

> __Please note, these instructions are a work in progress and may not work on your host.__

If you are using an Apple Silicon computer as your host machine, VirtualBox is not an option. Apple Silicon’s CPU is based off the ARM architecture, and VirtualBox only works on x86.

Instead, you'll need to use a virtualization technology like Parallels (https://parallels.com), VMware Fusion (https://vmware.com), or Qemu (https://www.qemu.org) to create an ARM-based virtual machine.

The first two options are paid software and may provide a better user experience. The last option, Qemu, is open source and will require some extra configuration.

Please see the `vagrant/` directory and `minikube/` directory for more information.
