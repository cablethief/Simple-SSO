# Simple-SSO

Start of a service that allows for a simple SSO login for Traefik forward-auth.

(Very much beta :D)

Uses bcrypt to keep pass, requires that this is set in the environment for the moment. 

Required environment variables and examples:

```
"DOMAIN=local.host"
"USERNAME=testing"
"PASSWORD=$$2a$$14$$9hNpsQYk9Y4iP4en62e.UuqW3lIlpvj6MU9ejiT44ELDTmLqA.Zha"
```

A script has been provided to easily bcrypt your password for you. 

```
go run gen-bcrypt.go -p "P@ssw0rd"
```

An example traefik config has been provided. 

## Why

I would prefer to login to a single place and have a cookie for all my services, rather than many individual logins. 

## What it does

The service provides: 

1. A page to provide a cookie on successful login for the configured domain.
2. A page which returns a `200 OK` if the cookie is valid.

The service will also run on `0.0.0.0:8000` 

## How 2

check the provided compose files for more details. set your forward auth to the auth container for the container you are protecting:

```
        - "traefik.http.middlewares.test-auth.forwardauth.address=http://auth:8000/"
```

On the auth container setup traefik to handle the auth request. 

```
        - "traefik.http.routers.auth.rule=Path(`/simple-sso-signin`)"
        # Required to be the primary rule to handle the creds coming through
        - "traefik.http.routers.auth.priority=100"
```


## TODO

 - Fancyness so that you dont have to manually visit page.
 - Read Users from an htaccess file.
 - Read Users from a traefik basic auth format label.
 - Add options for interface etc. 
 - Add basic auth option.
