# Service Overview

This is the telnet-server and provides nostalgic vibes for its users.

## How to Build

The application is written in the Go programming language.

You can use the `skaffold` local development setup, or you can run:

```bash
go build main.go
```

## Deploying

This application runs in a Kubernetes cluster and can be deployed
using the normal CI/CD pipeline.

## Common Tasks

None...yet

## On Call

This section describes how to manage any on-call alerts that come up.

## Alerts

### High CPU Throttling

If you are reaching this page via your email alert, __Congratulations!__

This alert occurs when the Pod is using too much CPU time and K8s starts
throttling the container in the Pod.

One of the reasons this can occur is because of too many connections.
Please check the connections rate to correlate. If connections are low, please
check the host and Pod metrics. Also, check the logs for any issues as well.

### High Error Rate

If you are reaching this page via your email alert, __Congratulations!__

This alert occurs when the application has a high error count.
Errors can be anything from unexpected dropped network connections
to invalid commands from users. Please view the logs to identify any
errors.

### High Connection Rate

If you are reaching this page via your email alert, __Congratulations!__

This alert occurs when there are many connections to the service.
Please make sure there are enough Pods in the Deployment to accommodate the spike.
Auto-scaling should be enabled, so there is a chance that is broken too.

## Disaster Recovery

See Deploy section for redeploying application.

## Service Level Objectives

> telnet-server will be available 99% of the time in a 14-day window.
> telnet-server will have an error rate less than 2 95% of the time for 14-day window.
