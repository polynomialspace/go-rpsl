# go-rpsl
RPSL parser for go

loose parser for [RPSL](http://www.irr.net/docs/rpsl.html) format files.

because several irr db sources use slightly different specifications, parsed output is a `map[string][]string`, which may or may not work well for other usecases.
additionally this fork does some [extra nonsense](https://github.com/polynomialspace/go-rpsl/commit/a113fd539a7e06d2dc53e03d04db3e0e0faa13c4) for our specific use.

see [this example](https://github.com/polynomialspace/go-rpsl/blob/master/cmd/rpsl-lookup/main.go) for basic usage.
