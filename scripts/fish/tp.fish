function tp
    if test -z $argv
        teleport
    else if contains $argv[1] "add" "remove" "rm" "list" "ls" "help"
        teleport $argv
    else
        set dir (teleport $argv)
        if test $status -eq 0
            cd $dir
        end
    end
end