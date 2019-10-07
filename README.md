# Simple-SSO

Start of a service that allows for a simple SSO login for Traefik forward-auth.

Uses bcrypt to keep files, requires that this is set in the environment for the moment. 

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


## What it does

The service provides 2 functions. 

1. A login page to provide a cookie on successful login for the configured domain. (`/signin`)
2. A check page which returns a `200 OK` if the cookie is valid. (`/check`)

The service will also run on `0.0.0.0:8000` 

## TODO

 - Fancyness so that you dont have to manually visit page.
 - Read Users from an htaccess file.
 - Read Users from a traefik basic auth format label.
 - Add options for interface etc. 
