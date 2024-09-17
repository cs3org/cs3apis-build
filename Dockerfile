# Docker file to build the cs3apis
# To push locally:
# docker build .
# docker tag xxxx cs3org/cs3apis:latest
# docker push cs3org/cs3apis
FROM golang:1.18-bullseye
RUN apt-get update
RUN apt-get install build-essential curl unzip sudo -y
RUN apt-get install python3-pip python3-full -y

# deps for node.js
RUN sudo apt-get update && sudo apt-get install -y ca-certificates curl gnupg
RUN curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
RUN echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_20.x nodistro main" | sudo tee /etc/apt/sources.list.d/nodesource.list
RUN apt-get update
RUN sudo apt-get install nodejs -y 
RUN node -v
RUN npm install -g @bufbuild/buf

# compile build tool and put it into path
ADD . /root/cs3apis-build
RUN cd /root/cs3apis-build/ && go build . &&  sudo cp cs3apis-build /usr/local/bin && sudo chmod u+x cs3apis-build

WORKDIR /root/cs3apis
