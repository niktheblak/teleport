function tp () {
    local tp_cmd=0
    if [[ $# -eq 0 ]]; then
        teleport
        return 0
    fi
    case $1 in
        add)
            tp_cmd=1
            ;;
        remove|rm)
            tp_cmd=1
            ;;
        list|ls)
            tp_cmd=1
            ;;
        prune)
            tp_cmd=1
            ;;
        help)
            tp_cmd=1
            ;;
    esac
    if [[ "$tp_cmd" -eq 1 ]]; then
        teleport "$@"
        return 0
    fi
    local dir=$(teleport warp "$1")
    if [[ "$?" -eq 0 ]]; then
        cd "$dir"
    fi
}
