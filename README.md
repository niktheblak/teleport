# teleport

Directory navigation tool for *nix shells. Contributions for integrations to different shells welcome!

## Usage

Start by building the `teleport` binary and copying it to your executable path.
You also need to copy the shell script to Fish shell's (TODO: support for other shells) `functions` directory.

```shell
$ go build cmd/teleport/teleport.go
$ cp teleport ~/bin
```

If you're using Bash shell, load the file `scripts/bash/tp.sh` into your current Bash session
or do it in your `.profile` file:

```shell
$ source scripts/bash/tp.sh
```

If you're using Fish shell, copy the file `scripts/fish/tp.fish` under your Fish functions:

```shell
$ cp scripts/fish/tp.fish ~/.config/fish/functions/
```

Now you can see the available commands with the `tp` command.

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
