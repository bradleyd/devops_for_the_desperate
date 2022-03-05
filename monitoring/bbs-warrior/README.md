# bbs-warrior

This docker image is hosted on GitHub's Docker registry.

If the image fails to pull from GitHub, you can build it locally in your minikube setup.

## To Run

> Make sure you have Go installed and in your path.

Enter the following command to make sure the application runs on your host:

```bash
go run main.go
```

## To Build

> Make sure your `DOCKER_HOST` variable is pointing at your minikube Docker environment. See `minikube docker env` for more information.

Next, enter the following command while in the `bbs-warrior/` directory to
build the bbs-warrior docker image.

```bash
docker build -t ghcr.io/bradleyd/devops_for_the_desperate/bbs-warrior:latest .
```

After this, when you install the monitoring stack, it should find the image. 

* If you decide to change anything in `main.go` don't forget to rebuild and redeploy using `kubectl apply -f cron_job.yaml` command.
