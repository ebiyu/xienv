# =============== start xienv init script ===============

mkdir -p $HOME/.xienv/bin

# install shims script
cat <<'EOF' > $HOME/.xienv/bin/xienv-run
#!/bin/bash

set -eu

path_remove ()  { export PATH=`echo -n $PATH | awk -v RS=: -v ORS=: '$0 != "'$1'"' | sed 's/:$//'`; }
path_remove "$HOME/.xienv/bin"

xienv check
. /tools/Xilinx/Vivado/$(xienv version)/settings64.sh

command=$1
shift 1
$command "$@"
EOF
chmod +x $HOME/.xienv/bin/xienv-run


# install shims script
cat <<'EOF' > $HOME/.xienv/run
#!/bin/bash

set -eu

path_remove ()  { export PATH=`echo -n $PATH | awk -v RS=: -v ORS=: '$0 != "'$1'"' | sed 's/:$//'`; }
path_remove "$HOME/.xienv/bin"

xienv check
. /tools/Xilinx/Vivado/$(xienv version)/settings64.sh

${0##*/} "$@"
EOF
chmod +x $HOME/.xienv/run

for app in \
vivado \
vitis \
xsct \
bootgen \
dtc
do
ls $HOME/.xienv/bin/$app > /dev/null 2>&1 || ln -s $HOME/.xienv/run $HOME/.xienv/bin/$app
done

#path_remove $HOME/.xienv/bin
export PATH="$HOME/.xienv/bin:$PATH"

#alias vivado="xienv_run vivado"
#alias vitis="xienv_run vitis"
#alias xsct="xienv_run xsct"

# shell completion
_xienv_cmd() {
    local -a _c
    _c=(
        'global:Set or show the global vivado/vitis version(s)'
        'local:Set or show the local project-specific vivado/vitis version(s)'
        'version:Show the current vivado/vitis version(s) and its origin'
        'versions:List all vivado/vitis versions available to xienv'
    )

    _describe -t commands Commands _c
}

_xienv_cmp() {
    _arguments \
        '1: :_xienv_cmd' \
        '*:: :->args' \

    case $state in
        (args)
            case $words[1] in
                (global)
                    _values 'versions' $(xienv versions --short)
                    ;;
                (local)
                    _values 'versions' $(xienv versions --short)
                    ;;
            esac
            ;;
    esac
}

compdef _xienv_cmp xienv

# =============== end xienv init script ===============
