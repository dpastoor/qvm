# Quarto Version Manager

manage and switch between versions of Quarto.

## setup


```shell
qvm use $QUARTO_VERSION 
```

```shell
qvm use $QUARTO_VERSION 
```

```shell
qvm install # alias for latest
qvm install 0.9.234
```

what is available on the system

```
qvm ls
```

```
qvm ls --remote
```

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

