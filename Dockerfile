# Docker file to build the cs3apis
# To push locally:
# docker build .
# docker tag xxxx cs3org/cs3apis:latest
# docker push cs3org/cs3apis
FROM golang
RUN apt-get update
RUN apt-get install build-essential curl unzip sudo ca-certificates gnupg -y
RUN apt-get install python3-pip python3-full -y

# deps for node.js
RUN curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg
RUN echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_20.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list

# deps for GH CLI
RUN curl -fsSL https://cli.github.com/packages/githubcli-archive-keyring.gpg | gpg --dearmor -o /usr/share/keyrings/githubcli-archive-keyring.gpg
RUN echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/githubcli-archive-keyring.gpg] https://cli.github.com/packages stable main" | tee /etc/apt/sources.list.d/github-cli.list

RUN apt-get update
RUN sudo apt-get install gh nodejs -y 
RUN node -v
RUN npm install -g @bufbuild/buf

# compile build tool and put it into path
ADD . /root/cs3apis-build
RUN cd /root/cs3apis-build/ && go build . &&  sudo cp cs3apis-build /usr/local/bin && sudo chmod u+x cs3apis-build

WORKDIR /root/cs3apis
