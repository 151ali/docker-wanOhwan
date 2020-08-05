# docker-wanOhwan
docker 101

## Build and test your image
```bash
docker build --tag {NAME}:{TAG} .
```

## Run your image as a container
```bash
docker run --publish 8000:8080 --detach --name {containerName} {NAME}:{TAG}
```