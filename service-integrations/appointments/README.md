# Appointments Service

This service allow users to schedule new appointments and list them all. 


## Create a container with `ko`

Running the following command creates a container for all platforms (amd64, arm64) and push it to Docker Hub if your account is configured (with `docker login`)

```shell
ko build appointments.go --platform=all
```

Then you can run this container locally by running: 

```
docker run docker.io/salaboy/appointments.go-6ae2051fa00e6032a383a1bd2a5aa1b4@sha256:ce6e8acd4395196a9f3b554115bc728777540c2cdd2cf77b741a1ceedf4263a7
Unable to find image 'salaboy/appointments.go-6ae2051fa00e6032a383a1bd2a5aa1b4@sha256:ce6e8acd4395196a9f3b554115bc728777540c2cdd2cf77b741a1ceedf4263a7' locally
9816217da073: Download complete 
493349de57e3: Download complete 
250c06f7c38e: Download complete 
b1e5889354e0: Download complete 
2025/03/29 10:14:56 Starting Appointments Service in Port: 8081
2025/03/29 10:14:56 Connecting to Database: postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable.
2025/03/29 10:14:56 dial tcp [::1]:5432: connect: connection refused
```
The container try to starts the service but it fails because there is no PostgreSQL running. 


## Deploy on Kubernetes

Create a Kind cluster with: 

```shell
kind create cluster
```

You should see something like this: 

```shell
Creating cluster "kind" ...
 ✓ Ensuring node image (kindest/node:v1.27.3) 🖼
 ✓ Preparing nodes 📦  
 ✓ Writing configuration 📜 
 ✓ Starting control-plane 🕹️ 
 ✓ Installing CNI 🔌 
 ✓ Installing StorageClass 💾 
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Not sure what to do next? 😅  Check out https://kind.sigs.k8s.io/docs/user/quick-start/
```

Then we can build, package and deploy the appointments service to our cluster by running: 


```shell
ko apply -f kubernetes/
```

This also should fail on Kubernetes, you can check the logs: 

```
k logs -f appointments-788bcd88fd-5nr78
2025/03/29 10:27:02 Starting Appointments Service in Port: 8081
2025/03/29 10:27:02 Connecting to Database: postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable.
2025/03/29 10:27:02 dial tcp [::1]:5432: connect: connection refused
```

Installing PostgreSQL using Helm:

```shell
helm install postgresql oci://registry-1.docker.io/bitnamicharts/postgresql
```