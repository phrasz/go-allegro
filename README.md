NAG Bindings - New Allegro Go Bindings
======================================

This library provides bindings for Allegro 5.2.

The bindings are dependent upon a pre-installed Allegro development environment.
This work extends the original work of [go-allegro](https://github.com/dradtke/go-allegro).

The following functions are not implemented, due to coverage in Go.

-
-
-
-

The original bindiging had 100% coverage of Allegro 5.0.10.

Test coverage
-------------
Run: `go test coverage_test.go`

Linux Installation
------------------
Assuming Allegro 5.2 has been built, tun the following commands:
1. cd && `pkg-config`
2. go get -d github.com/dradtke/go-allegro`.
3. `go install github.com/dradtke/go-allegro/allegro`

<!-- Windows Installation
-------
1. Set the `ALLEGRO_HOME` to this github.com/phrasz/go-allegro
2. Set `ALLEGRO_VERSION` to the version of Allegro downloaded ( e.g., 5.0.10)
3. (Optional) Set `ALLEGRO_LIB` to match which allegro version (default value is `monolith-static-mt-debug`).
4. Run `setenv.bat`; if it successfully runs, then build and install the library
-->
