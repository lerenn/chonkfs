# Backends

This directory contains the backends for the chonkfs package. A backend is a
module that implements the low-level operations that the filesystem needs to
perform. The backends are responsible for reading and writing data to the
underlying storage, as well as managing the metadata that the filesystem needs
to keep track of.

The backends are designed to be pluggable, so that different backends can be
used depending on the requirements of the user. For example, a user might want
to use a backend that stores data in a cloud storage service, or a backend that
stores data on a local disk.

Here are the backends that are currently implemented:
* `mem`: A backend that stores data in memory. This backend is useful for
  testing, development and caching purposes, as it does not persist data across
  restarts.

Each backend is implemented as a separate module in this directory. The module
should export a struct that implements the `Backend` interface defined in
`backend.go`.

Each backend should also being tested with tests from the package `pkg/backends/test`.