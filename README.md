# Winget auto update

## Building

```
go build -o wingetau.exe
```

## Todo

Need to add silent mode somehow (to turn off notifications). This can also be done with an environment variable.
Alternative is to add a json config file.

An option can be to to add a command to create the config file. This way the user can decide if they want to use an env variable or a .json config.