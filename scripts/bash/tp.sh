function tp () {
    local args=($*)
    local tp_cmd=0
    case ${args[1]} in
        add)
            tp_cmd=1
            ;;
        remove|rm)
            tp_cmd=1
            ;;
        list|ls)
            tp_cmd=1
            ;;
        help)
            tp_cmd=1
            ;;
    esac
    if [[ -z "$args" ]]; then
        teleport
        return 0
    fi
    if [[ "$tp_cmd" -eq 1 ]]; then
        teleport "$args"
        return 0
    fi
    local dir=$(teleport "$args")
    if [[ "$?" -eq 0 ]]; then
        cd "$dir"
    fi
}
