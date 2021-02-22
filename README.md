# getip-hit-counter

## Getting Started Using Docker Compose

Build and start the hit-counter and redis containers:

    $ docker-compose up -d --build

### Using Hit-Counter

    $ curl http://<IP-Address>
    {"User Address ":"<IP-Address>","Hit ":1}