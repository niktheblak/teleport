# teleport

Directory navigation tool for Fish shell.

## Usage

Start by building the `teleport` binary and copying it to your executable path.
You also need to copy the shell script to Fish shell's (TODO: support for other shells) `functions` directory.

```shell
$ go build
$ cp teleport ~/bin
$ cp tp.fish ~/.config/fish/functions/
```

Now you can see the available commands with the `teleport` command.

After this you can create teleport points to your favorite directories with `tp add`.

## Adding a teleport point

tp add {name} [directory]

Omitting [directory] will add a teleport point to the current directory.

For example:
```shell
$ cd ~/projects/go/my_fav_proj
$ tp add fav
$ cd
$ tp fav
# pwd is now at ~/projects/go/my_fav_proj
```

## Listing teleport points

tp list

For example:
```shell
$ tp list
```
