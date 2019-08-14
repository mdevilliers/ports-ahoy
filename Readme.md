# ports-ahoy

TODO: add a description for your project


## Develop

To list makefile targets

```
make help
```

Golang
-----

To build the application locally -

```
make build
```

The built image is outputted to the /bin folder

Docker
------

To build the Docker image -

```
make image
```

The application is built locally before being published to your configured Docker repository.


Kubernetes
----------

To deploy to Kubernetes -

```
make image
make deploy
```
Note that the image points to the `latest` tag for developing locally.

Please remember that 'latest is not a version' and amend for your production deploy accordingly.