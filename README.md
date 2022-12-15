## About

Android Toolbox 



## Building

first install golang、npm、upx(optional)，and wails

```
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

check env, if all conditions are met, just build

```
wails doctor
wails build -upx
```

## debug

and you can use `wails dev` to debug
