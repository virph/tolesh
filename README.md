# Tole SH
Tole SH is a program to enable easy remote to the server.

Usage: `./tolesh <command> [arguments]`

The commands are:
- fetch: fetch the server nodes information with `tsh ls` command
- list: print the server nodes list
- copy: copy the server nodes list to the desired path
- iterm: open the server sessions using iTerm2 terminal

To use the `iterm` command, iTerm2 application is required.

List of arguments:
- Global arguments:
	- v: print verbose logs (default: false)
	- p: data path (default: current directory / ".")
- Fetch command:
	- s: use sample node list (default: false)
- List command:
	- g: hostgroup name (default: all)
	- t: node type (default: all)
- Copy command:
	- d: destination path (required)
- Iterm command:
	- g: hostgroup name (required)
	- u: username (default: root)

Example of commands:
- `tolesh fetch -p "~/.tolesh"`
- `tolesh list -g archiveapp`
- `tolesh list -t elasticsearch`
- `tolesh copy -d "~/some/new/place.txt"`
- `tolesh iterm -g notifapp -u rut`
