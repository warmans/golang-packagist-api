package packagist

import (
    "fmt"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/suite"
)

type PackagistTestSuite struct {
    suite.Suite
    packagist Client
}

func (suite *PackagistTestSuite) SetupTest() {
    suite.packagist = NewAPIClient()
}

func (suite *PackagistTestSuite) TestMakeUriCreatesValidURIs() {

    //uri without querystring
    suite.Equal(suite.packagist.MakeURI("/packages", make(map[string]string)), "https://packagist.org/packages")

    //uri with querystring
    suite.Equal(suite.packagist.MakeURI("/packages", map[string]string{"foo": "bar"}), "https://packagist.org/packages?foo=bar")
}

func (suite *PackagistTestSuite) TestListPackages() {

    //setup server
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "{\"packageNames\":[\"vendor1/package1\",\"vendor2/package2\"]}")
    }))
    defer ts.Close()

    //force client to use server
    suite.packagist.host = ts.URL

    res, _ := suite.packagist.ListPackages(make(map[string]string))

    //check for expected contents
    suite.Equal(res.PackageNames[0], "vendor1/package1")
    suite.Equal(res.PackageNames[1], "vendor2/package2")
}

func (suite *PackagistTestSuite) TestListPackagesDecodeFailure() {

    //setup server
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, "{\"packageNames\":[\"vendor1/package1\",\"vendor2/package2\"")
    }))
    defer ts.Close()

    //force client to use server
    suite.packagist.host = ts.URL

    res, err := suite.packagist.ListPackages(make(map[string]string))

    //empty package list result
    suite.Equal(res, PackageListResult{})

    //check for error
    suite.Equal(err.Error(), "unexpected EOF")
}

func (suite *PackagistTestSuite) TestGetPackage() {

    //setup server
    ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, `{
            "package": {
                "name": "warmans/dlock",
                "description": "Distributed locking library",
                "time": "2014-05-11T15:01:42+00:00",
                "maintainers": [{
                    "name": "warmans"
                }],
                "versions": {
                    "dev-master": {
                        "name": "warmans/dlock",
                        "description": "Distributed locking library",
                        "keywords": [
                            "locking",
                            "lock",
                            "distributed"
                        ],
                        "homepage": "https://github.com/warmans/dlock",
                        "version": "dev-master",
                        "version_normalized": "9999999-dev",
                        "license": [],
                        "authors": [{
                            "name": "Stefan Warman",
                            "email": "stefan.warman@gmail.com"
                        }],
                        "source": {
                            "type": "git",
                            "url": "https://github.com/warmans/dlock.git",
                            "reference": "b9d48c7aba1bcafec4d2a34583b95f5cf519021b"
                        },
                        "dist": {
                            "type": "zip",
                            "url": "https://api.github.com/repos/warmans/dlock/zipball/b9d48c7aba1bcafec4d2a34583b95f5cf519021b",
                            "reference": "b9d48c7aba1bcafec4d2a34583b95f5cf519021b",
                            "shasum": ""
                        },
                        "type": "library",
                        "time": "2014-05-11T21:47:51+00:00",
                        "autoload": {
                            "psr-4": {
                                "Dlock\\": "src/"
                            }
                        },
                        "require": {
                            "php": ">=5.3.0"
                        },
                        "suggest": {
                            "ext-redis": "For redis datastore",
                            "ext-memcache": "For memcache datastore"
                        }
                    }
                },
                "type": "library",
                "repository": "https://github.com/warmans/dlock",
                "downloads": {
                    "total": 13,
                    "monthly": 0,
                    "daily": 0
                    },
                    "favers": 0
                }
            }
        `)
    }))
    defer ts.Close()

    //force client to use server
    suite.packagist.host = ts.URL

    res, err := suite.packagist.GetPackage("vendor1/package1")

    suite.Equal(err, nil)                                                 //no error
    suite.Equal(res.Package.Name, "warmans/dlock")                        //name looks okay
    suite.Equal(res.Package.Versions["dev-master"].Name, "warmans/dlock") //version name looks okay

}

func TestPackagist(t *testing.T) {
    suite.Run(t, new(PackagistTestSuite))
}
