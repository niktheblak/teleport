function tp
    if contains $argv[1] "add" "remove" "rm" "list" "ls"
        teleport $argv
    else
        set dir (teleport $argv)
        if test $status -eq 0
            cd $dir
        end
    end
end