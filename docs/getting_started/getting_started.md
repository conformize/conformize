# Getting started

## Installation

**Conformize** can be installed either by compiling the source code or downloading a pre-built binary.

### Installing via pre-built binary

1. Download the latest pre-built binary from the [releases page](https://github.com/conformize/conformize/releases).
2. Extract the binary to a directory of your choice.
3. Add the directory to your system's **PATH** environment variable.

### Compiling from Source

#### 1. Install prerequisites:

- [Go](https://golang.org/dl/)
- [GNU Make](https://www.gnu.org/software/make/)

Alternatively, GNU Make could be installed on macOS via Homebrew:

```
$ brew install make
```

#### 2. Clone the repository:

```
$ git clone https://github.com/conformize/conformize.git
```

#### 3. Change the working directory to the one where the repository is cloned:

```
$ cd conformize
```

#### 4. Build the project using GNU Make:

```
$ make build
```

Shortly after, you should see an output like the one below:

> Binary is available at ./build/dev/conformize

#### 5. Copy the binary from the designated output path to a location of your choice.
#### 6. Add the path to the chosen location to your system's PATH environment variable.

After installing via either of the methods above, you can verify the installation by running:

```
$ conformize version
```

If all is well, depending on the version and the system that **Conforimze** is running, you should see similar output:

```                                                                    
Conformize v0.1.0
running on darwin arm64
```

Or you can run the `help` command like this, to see list of available commands.
```
$ conformize help
```

Truth is, we can't do anything without a [Blueprint](../blueprint/what_is_a_blueprint.md), so let's learn more about what it is and create one.
