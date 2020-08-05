# BP 
## Pipe Dockerfile through `stdin`

bash```
docker build -t {imageName}:{TAG} -<<EOF
FROM busybox
RUN echo "hello world"
EOF
```  
> Note: Attempting to build a Dockerfile that uses COPY or ADD will fail if this syntax is used  

