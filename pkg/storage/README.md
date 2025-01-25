# Backends

This directory contains the storages for the chonkfs package. A storage is a
module that implements the low-level operations that the filesystem needs to
perform. The storages are responsible for reading and writing data to the
underlying storage, as well as managing the metadata that the filesystem needs
to keep track of.

The storages are designed to be pluggable, so that different storages can be
used depending on the requirements of the user. For example, a user might want
to use a storage that stores data in a cloud storage service, or a storage that
stores data on a local disk.

Here are the storages that are currently implemented:
* `mem`: A storage that stores data in memory. This storage is useful for
  testing, development and caching purposes, as it does not persist data across
  restarts.

Each storage is implemented as a separate module in this directory. The module
should export a struct that implements the `Backend` interface defined in
`storage.go`.

Each storage should also being tested with tests from the package `pkg/storages/test`.