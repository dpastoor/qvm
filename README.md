# Quarto Version Manager

manage and switch between versions of Quarto.

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

### understanding what is available

which versions of quarto do you have installed?

```
qvm ls
```

which ones are available remotely? 

```
qvm ls --remote
qvm ls --remote --since 2022-05-01 # releases since 2022-05-01, follows YYYY-MM-DD
qvm ls --remote -n 10 # latest 10 releases
```

### installing a version

```shell
qvm install $QUARTO_VERSION
```

can also allow interactive selection and filtering by just
running `qvm install`

TODO: example of what this looks like:

```
qvm install
```

### updating which global version to use

```
qvm use $QUARTO_VERSION
```

### using a version in a particular shell session

At any point can use a particular version in your shell session:

```shell
qvm use $QUARTO_VERSION --local
```


### programmatic support utilities

TODO: explain - goal is to allow qvm to be used in other shell scripts/etc
more easily

```
qvm path
```


## internals

If you want to poke around yourself, the config file location follows the XDG base directory spec. 


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

