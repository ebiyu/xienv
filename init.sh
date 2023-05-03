# Xilinx
function xienv_run() {
    xienv check && 
    zsh -c ". /tools/Xilinx/Vivado/$(xienv version)/settings64.sh &&. /tools/Xilinx/Vitis/$(xienv version)/settings64.sh && $*"
}

alias vivado="xienv_run vivado"
alias vitis="xienv_run vitis"
alias xsct="xienv_run xsct"
