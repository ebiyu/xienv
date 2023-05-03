# xienv

Vivado/Vitis version management tool for Linux like `pyenv`.

Vivado/Vitis must be installed in `/tools/Xilinx/Vivado` and `/tools/Xilinx/Vitis` respectively.

## Installation

Please install via `go install`:

```sh
go install github.com/ebiyuu1121/xienv@latest
```

and open your `.zshrc`, `.bashrc`, or `.bash_profile` and add the following line:

```sh
eval "$(xienv init -)"
```

## Usage

```sh
xienv global 2020.2 # set global version
xienv local 2020.2  # set project-local version
xienv versions      # list installed versions

# run Vivado/Vitis with global/local version
# If both global and local versions are set, the local version is used.
# If both global and local versions are not set, the latest version is used.
vivado
vitis

# If some commands are not found, please contact me.
# Workaround:
xienv_run vivado
```
