# Docker file to build the cs3apis
# To push locally:
# docker build .
# docker tag xxxx cs3org/cs3apis:latest
# docker push cs3org/cs3apis
FROM golang
RUN apt-get update
RUN apt-get install build-essential curl unzip sudo -y
RUN apt-get install python3-pip python3-full -y

# deps for protoc
RUN cd /tmp && curl -sSL https://github.com/protocolbuffers/protobuf/releases/download/v25.0/protoc-25.0-linux-x86_64.zip -o protoc.zip && unzip -o protoc.zip && sudo cp bin/protoc /usr/local/bin/protoc
RUN cd /tmp && curl -sSL https://github.com/uber/prototool/releases/download/v1.10.0/prototool-Linux-x86_64 -o prototool && sudo cp prototool /usr/local/bin/ && sudo chmod u+x /usr/local/bin/prototool
RUN cd /tmp && curl -sSL https://github.com/nilslice/protolock/releases/download/v0.16.0/protolock.20220302T184110Z.linux-amd64.tgz -o protolock.tgz && tar -xzf protolock.tgz && sudo cp protolock /usr/local/bin/
RUN cd /tmp && curl -sSL https://github.com/pseudomuto/protoc-gen-doc/releases/download/v1.5.1/protoc-gen-doc_1.5.1_linux_amd64.tar.gz -o protoc-gen-doc.tar.gz && tar xzfv protoc-gen-doc.tar.gz && sudo cp protoc-gen-doc /usr/local/bin/
RUN go install github.com/golang/protobuf/protoc-gen-go@v1.5.3


# deps for python
RUN pip install grpcio grpcio-tools --ignore-installed --break-system-packages

# deps for js
RUN curl -sSL https://github.com/grpc/grpc-web/releases/download/1.5.0/protoc-gen-grpc-web-1.5.0-linux-x86_64 -o /tmp/protoc-gen-grpc-web
RUN sudo mv /tmp/protoc-gen-grpc-web /usr/local/bin/ && sudo chmod u+x /usr/local/bin/protoc-gen-grpc-web

# deps for node.js
RUN sudo apt-get update && sudo apt-get install -y ca-certificates curl gnupg
RUN curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
RUN echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_20.x nodistro main" | sudo tee /etc/apt/sources.list.d/nodesource.list
RUN apt-get update
RUN sudo apt-get install nodejs -y 
RUN node -v
RUN npm install protoc-gen-grpc -g -unsafe-perm

# compile build tool and put it into path
ADD . /root/cs3apis-build
RUN cd /root/cs3apis-build/ && go build . &&  sudo cp cs3apis-build /usr/local/bin && sudo chmod u+x cs3apis-build

WORKDIR /root/cs3apis
