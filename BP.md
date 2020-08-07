# BP 
## Pipe Dockerfile through `stdin`

```bash
docker build -t {imageName}:{TAG} -<<EOF
FROM busybox
RUN echo "hello world"
EOF
```  
> Note: Attempting to build a Dockerfile that uses COPY or ADD will fail if this syntax is used  

## Remarks
- Don’t install unnecessary packages  
- Decouple applications into multiple containers makes it easier to scale horizontally and reuse containers  
- Minimize the number of layers  
  - Only the instructions `RUN`, `COPY`, `ADD` create layers. Other instructions create temporary intermediate images, and do not increase the size of the build.  
  - Where possible, use multi-stage builds.  
- Sort multi-line arguments, example :  

```Dockerfile
RUN apt-get update && apt-get install -y \
  curl \
  tree \
  git \
&& rm -rf /var/lib/apt/lists/*
```

- Explicitly choose a shell that does support the `pipefail` option. For example:

```Dockerfile
RUN ["/bin/bash", "-c", "set -o pipefail && wget -O - https://some.site | wc -l > /number"]
```

- Use `ENV PATH /usr/local/nginx/bin:$PATH` ensures that `CMD ["nginx"]` just works.  
- Set commonly used version numbers. For example: `ENV PG_MAJOR 9.3` .
- Each `ENV` line creates a new intermediate layer. If you unset the environment variable in a future layer, it `still persists` in this layer and its value `can’t be dumped`. . For example:

```Dockerfile
FROM alpine
ENV ADMIN_USER="mark"
RUN echo $ADMIN_USER > ./mark
RUN unset ADMIN_USER
```
To prevent this, use a `RUN` command with shell commands, to **set**, **use**, and **unset** the variable all in a single layer.
```Dockerfile 
FROM alpine
RUN export ADMIN_USER="mark" \
    && echo $ADMIN_USER > ./mark \
    && unset ADMIN_USER
CMD sh

```

```bash 
docker run --rm test sh -c 'echo $ADMIN_USER'
```
- The best use for `ADD` is local tar file __auto-extraction__ into the image, as in `ADD rootfs.tar.xz /`.
- The best use for `ENTRYPOINT` is to set the image’s main command, allowing that image to be run as though it was that command (and then use `CMD` as the default flags).
```Dockerfile
ENTRYPOINT ["s3cmd"]
CMD ["--help"]
```  
Now `docker run s3cmd` command show the command’s help

- Use`WORKDIR` instead of `RUN cd ... && do-something`  

- Leverage build cache , otherwise use `--no-cache=true` on the `docker build` command.