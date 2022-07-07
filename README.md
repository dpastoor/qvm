# Quarto Version Manager

manage and switch between versions of Quarto.

[![asciicast](https://asciinema.org/a/nB2VzKCeuW0iyBANuRGCVPywp.svg)](https://asciinema.org/a/nB2VzKCeuW0iyBANuRGCVPywp)

## how it works

taking inspiration from pyenv/conda/homebrew/nvm/choco/other env managers, qvm supports
two mechanisms for managing your quarto version:
*  a user-level global version
   *  scoop
   *  nvm w/ global settings
   *  homebrew (though not as easy to switch versions quickly in homebrew)
*  a version for the particular shell session
   *  examples: pyenv/nvm, conda

The global version can be helpful for managing your overall system configuration
where you want discover how your current versions compare to whats available,
upgrade, and ultimately be able to swap versions if you need to. 

The quarto team's working hypothesis is that the rate of breaking changes in quarto
by v1.0 should be sufficiently low that this should be adequate for the majority
of users. That said, there are particular development scenarios that can be
nice to explore for the duration of a single shell session. In addition, this
will make it easier for users to experiment with the notion of project level
quarto versions more easily for the quarto team to understand how much that
workflow should be supported.

## setup

Run `qvm init` and follow the instructions!

The instructions approximately will recommend: 

add the following to your `~/.profile`:

```
export PATH="$(qvm path add)"
```

This will make sure the appropriate directories are prepended to your path to detect
the managed quarto version.

NOTE: adding this to your profile instead of `~/.bashrc` or `~/.zshrc` will
enable both the shell sessions *and* R sessions launched via Rstudio Workbench
to pick up the change. If you only add to the `rc` file, it will only be present
in the shell.

An example complete `~/.profile` file might look like:

```
# ~/.profile: executed by the command interpreter for login shells.
# This file is not read by bash(1), if ~/.bash_profile or ~/.bash_login
# exists.
# see /usr/share/doc/bash/examples/startup-files for examples.
# the files are located in the bash-doc package.

# the default umask is set in /etc/profile; for setting the umask
# for ssh logins, install and configure the libpam-umask package.
#umask 022

# if running bash
if [ -n "$BASH_VERSION" ]; then
    # include .bashrc if it exists
    if [ -f "$HOME/.bashrc" ]; then
        . "$HOME/.bashrc"
    fi
fi

# set PATH so it includes user's private bin if it exists
if [ -d "$HOME/bin" ] ; then
    PATH="$HOME/bin:$PATH"
fi

# set PATH so it includes user's private bin if it exists
if [ -d "$HOME/.local/bin" ] ; then
    PATH="$HOME/.local/bin:$PATH"
fi

export PATH="$(qvm path add)"

```

### understanding what is available

which versions of quarto do you have installed?

```
qvm ls
```

which ones are available remotely? 

```
qvm ls --remote
```

### installing versions

`install` will install quarto versions given the release name(s). It will
attempt to install all versions in parallel, and will coalesce and return
any failures at the end. 

You can also allow interactive selection and filtering by just
running `qvm install`

```
qvm install
```

```
qvm install
? Which version do you want to install?  [Use arrows to move, type to filter]
> v0.9.637
  v0.9.636 - **installed**
  v0.9.634
  v0.9.633
  v0.9.632
  v0.9.629 - **installed**
  v0.9.628
  v0.9.626
  v0.9.624
  v0.9.622
```

```shell
qvm install $QUARTO_VERSION $ANOTHER_QUARTO_VERSION
```

```shell
qvm install v0.9.466 v0.9.432
```

the `latest` keyword will dynamically install the latest version

```
qvm install latest
```

### updating which global version to use

```
qvm use $QUARTO_VERSION
```

To interactively select a version, run `qvm use`

```
qvm use
? Which version do you want to use?  [Use arrows to move, type to filter]
> v0.9.636
  v0.9.629 - **active**
  v0.9.587
  v0.9.583
  v0.9.565
  v0.9.563
  v0.9.562
  v0.9.561
  v0.9.559
  v0.9.550
```
### programmatic support utilities

```shell
qvm path roots # various roots
qvm path versions # path to versions
qvm path active # path to active bin dir
```


## internals

If you want to poke around yourself, the config file location follows the XDG base directory spec, with the caveat that on OSX, we normalize to follow the linux XDG
convention of using `.local/share` and `.config` instead of `Application Library`
as quarto has some trouble with paths with strings. Though this will be fixed in 
subsequent quarto versions, we don't want to make qvm incompatible with all 
quarto versions on OSX prior to that date.


<details open>
    <summary><strong>Unix-like operating systems</strong></summary>
    <br/>

| <a href="#xdg-base-directory"><img width="400" height="0"></a> | <a href="#xdg-base-directory"><img width="500" height="0"></a><p>Unix</p> | <a href="#xdg-base-directory"><img width="600" height="0"></a><p>macOS</p>                                            | <a href="#xdg-base-directory"><img width="500" height="0"></a><p>Plan 9</p> |
| :------------------------------------------------------------: | :-----------------------------------------------------------------------: | :-------------------------------------------------------------------------------------------------------------------: | :-------------------------------------------------------------------------: |
| <kbd><b>XDG_CONFIG_HOME</b></kbd>                              | <kbd>~/.config</kbd>                                                      | <kbd>~/Library/Application&nbsp;Support</kbd>                                                                         | <kbd>$home/lib</kbd>                                                        |
| <kbd><b>XDG_CONFIG_DIRS</b></kbd>                              | <kbd>/etc/xdg</kbd>                                                       | <kbd>~/Library/Preferences</kbd><br/><kbd>/Library/Application&nbsp;Support</kbd><br/><kbd>/Library/Preferences</kbd> | <kbd>/lib</kbd>                                                             |

</details>

<details open>
    <summary><strong>Microsoft Windows</strong></summary>
    <br/>

| <a href="#xdg-base-directory"><img width="400" height="0"></a> | <a href="#xdg-base-directory"><img width="700" height="0"></a><p>Known&nbsp;Folder(s)</p> | <a href="#xdg-base-directory"><img width="900" height="0"></a><p>Fallback(s)</p> |
| :------------------------------------------------------------: | :---------------------------------------------------------------------------------------: | :------------------------------------------------------------------------------: |
| <kbd><b>XDG_CONFIG_HOME</b></kbd>                              | <kbd>LocalAppData</kbd>                                                                   | <kbd>%LOCALAPPDATA%</kbd>                                                        |
| <kbd><b>XDG_CONFIG_DIRS</b></kbd>                              | <kbd>ProgramData</kbd><br/><kbd>RoamingAppData</kbd>                                      | <kbd>%ProgramData%</kbd><br/><kbd>%APPDATA%</kbd>                                |

</details>

