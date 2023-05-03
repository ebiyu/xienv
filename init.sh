# Xilinx
function xienv_run() {
    xienv check && 
    zsh -c ". /tools/Xilinx/Vivado/$(xienv version)/settings64.sh &&. /tools/Xilinx/Vitis/$(xienv version)/settings64.sh && $*"
}

alias vivado="xienv_run vivado"
alias vitis="xienv_run vitis"
alias xsct="xienv_run xsct"

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
