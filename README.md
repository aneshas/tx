# tx
[![Go Test](https://github.com/aneshas/tx/actions/workflows/test.yml/badge.svg)](https://github.com/aneshas/tx/actions/workflows/test.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/aneshas/tx)](https://goreportcard.com/report/github.com/aneshas/tx)
[![Coverage Status](https://coveralls.io/repos/github/aneshas/tx/badge.svg)](https://coveralls.io/github/aneshas/tx)
[![Go Reference](https://pkg.go.dev/badge/github.com/aneshas/tx.svg)](https://pkg.go.dev/github.com/aneshas/tx)

Package tx provides a simple transaction abstraction in order to enable decoupling/abstraction of persistence from
application/domain logic while still leaving transaction control to the application service.
(Something like @Transactional annotation in Java, without an annotation)

`go get github.com/aneshas/tx`
