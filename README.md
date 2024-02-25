# Conformize

## Introduction

**Conformize** is a tool that enables writing and validating test cases for application service configurations.\
It aims to easily integrate within various CI/CD pipeline setups, supporting static sources like JSON, YAML, XML, and other files,\
as well as stores such as HashiCorp Consul, Etcd, AWS SSM, and more.

## Getting started

### Installation

**Conformize** can be installed either by compiling the source code or downloading a pre-built binary.

#### Installing via pre-built binary

1. Download the latest pre-built binary from the [releases page](https://github.com/conformize/conformize/releases).
2. Extract the binary to a directory of your choice.
3. Add the directory to your system's PATH environment variable.

#### Compiling from Source

1. Install prerequisites:

    - [Go](https://golang.org/dl/)
    - [GNU Make](https://www.gnu.org/software/make/)

	Alternatively, GNU Make could be installed on macOS via Homebrew:

    ```
	brew install make
    ```

2. Clone the repository:

    ```
    git clone https://github.com/conformize/conformize.git
    ```

3. Navigate to the project directory:

    ```
	cd conformize
    ```

4. Build the project using GNU Make:

    ```
    make build
    ```

After installing via either of these methods, you can verify the installation by running:

```
conformize version
```

## License

[Mozilla Public License, version 2.0](./LICENSE.md)
