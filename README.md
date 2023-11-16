# cs3apis-build
Docker image to build the CS3APIS

## Development
See the [cs3apis local compiliation instructions](https://github.com/cs3org/cs3apis) and the
[cs3apis Makefile](https://github.com/cs3org/cs3apis/blob/main/Makefile)
for pointers on how to run and test the code in this repo on your local machine.

So for instance:
```sh
git clone https://github.com/cs3org/cs3apis-build
cd cs3apis-build
docker build -t cs3apis-build .
cd ..
git clone https://github.com/cs3org/cs3apis
cd cs3apis
docker run -v `pwd`:/root/cs3apis cs3apis-build cs3apis-build
cd build/go-cs3apis
git status
// see the result
```
