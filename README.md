Packagist API Client
=============================

Packagist is the main repository used by composer for PHP packages: https://packagist.org/

They have undocumented API for accessing information about registered packages. This client provides access to that
API.

### Usage

__packagist.NewAPIClient()__

Constructor to create a new client instance.


__ListPackages(filters map[string]string)__

Returns a list of all packages in packagist. These can be filtered using the first argument with:
* type (e.g. library)
* vendor (e.g. warmans)
* tags (e.g. cache)

__GetPackage(name string)__

Returns the details for a single package by name.

### Example

```Go
package main

import (
    "fmt"
    "github.com/warmans/golang-packagist-api"
)

func main() {
    client := packagist.NewAPIClient()

    //show a vendor's packages
    fmt.Println(client.ListPackages(map[string]string{"vendor":"warmans"}))

    //show everything for a specific package
    fmt.Println(client.GetPackage("guzzle/guzzle"))
}

```
