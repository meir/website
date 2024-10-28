# website

This is a personal site made using a static site generator entirely written in bash.
I would not suggest using this, unless you have a lot of confidence in bash and in me (I don't).

Unless you want to use a screwdriver as a fork, I would not suggest bash as a static site generator, it works, but it's cursed.
If you still do decide to use this, you're welcome to ask questions, but please do come to me with expectations to fix your scripts.

## Why make this in bash?
Why not?

But seriously, I like to experiment with the avant-garde of programming, a lot of projects don't work out and get deleted, this one so happened to work.
Some extra reasoning:
 - Don't need to install any compiler, bash is already on every system with a Unix shell
 - Dependencies are often already installed if you've used bash frequently
 - It won't break next year because of dependency updates (looking at you NodeJS)
 - It's easy to read, write, extend and modify since its just bash
 - No weird syntax code going on, its just html and bash 
 - No need for a difficult pipeline/setup/dev environment
 - It's incredibly small and quite fast

## Dependencies

- bash (5+)
- jq
- python3 (only for `./dev.sh`)

### Nix/NixOS

If you're using the [Nix](https://nixos.org/) package manager, consider getting [direnv](https://direnv.net/).
Once you have direnv installed you just have to run `direnv allow` to automatically install the dependencies using Nix.

## Developing

It's quite easy, get a file watcher for your editor if you want and just let it run `./dev.sh`

### NVIM
For NVIM I suggest installing [Overseer](https://github.com/stevearc/overseer.nvim).
Using this plugin you can start a job that will run `./dev.sh` and make it automatically watch all files for hot reload.

## Files

The files are divided into the following folders:
- assets
- scripts
- components
- src

### assets
The assets are all the files that are not html or bash scripts
You can decide your own folder system in here but keep in mind, once the generator moves the assets to the output directory, all the files will be in `/assets/`, so always keep all file names unique even across folders.

### scripts
In this directory are all the bash scripts, everything in here gets loaded once the `./generate.sh` script is called.
If you dont want this, consider moving scripts to a subdirectory or putting the script code in a bash function.
All scripts get sources from the root directory, so calling for a script file will require the full path such as `./scripts/your_script.sh`.

You can also add prerender hooks by calling `prerender_hook "function_name"` this can be used to perform tasks right after the prerender.

### components

In here you can define all the components that can be loaded from your html files such as layouts, navigations, code blocks, etc.
Components can be given content by sending an EOF to it, or be given individual variables by setting the variable before the command such as the following:
```
$(title="hello variable" component "component_file" <<EOF
    <p>
        Hello $(echo "world!")
    </p>
EOF
)
```

### src
This is where your site lives, all the individual pages are stored here and will be generated from here to the output directory.
The path is decided by the path and name you give the file in here.
For example `./src/cool_topics/music.htm` will generate `./$OUT/cool_topics/music/index.htm`.
This is so that the url in the browser will be just `/cool_topics/music`.

## Tags

Tags are to register pages in a map and use them later on, such as for navigation, related pages, etc.
This can be done by using `$(tag "topic")` in the page.

If you want to use the pages in the tag you can use it as followed:
```
$(for page in $(get_tagged "topic"); do
    cat <<EOF
    <a href="$(get_url $pags)">$(get_title $page)</a>
EOF
done)
```

## common issues

### EOF Issues

Be sure to have the closing tag in a newline with EOF, this will otherwise cause issues on Linux
So always do this:
```
EOF
)
```
And not `EOF)`

