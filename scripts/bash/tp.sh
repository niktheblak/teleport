#!/usr/bin/env bash

function tp () {
    local args=$1
    local tp_cmd=0
    case "${args[0]}" in
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
        return
    fi
    if [[ "$tp_cmd" -eq 1 ]]; then
        teleport "$args"
        return
    fi
    local dir=$(teleport "$args")
    if [[ "$?" -eq 0 ]]; then
        cd "$dir"
    fi
}
