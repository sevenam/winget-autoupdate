# Winget auto update

## Building

```
go build -o wingetau.exe
```

## Todo

Need to add silent mode somehow (to turn off notifications). This can also be done with an environment variable.
Alternative is to add a json config file.

An option can be to to add a command to create the config file. This way the user can decide if they want to use an env variable or a .json config.

There's a bug with notifications when running as a service. They're not getting through

Another bug - the service running as local system will not see all apps. the service needs to run as user and also in admin context. seems fairly impossible to set it up to run as a service using a microsoft account. seemingly services are no longer the way to go for a prive account on a machine doing something in the background, unless you set up another account for it and run the service as that one, but then it's not your account anymore and maybe winget won't list your apps, so generally pretty crappy. tried setting up credentials using credential manager, but that didn't work either.

for the above problem, maybe task scheduler is the way to go now? task scheduler seems to work, but gives annoying terminal popup unless build like this: go build -ldflags="-H windowsgui" -o wingetau.exe

other options could be: setting the cli to run on startup "with runasservice" argument and then handle it's own timer/schedule using cron or something. Actually - this won't work very well as a startup app like this wont easily be able to run with admin privileges. Anyway - example here:

package main

import (
    "os"
    "os/user"
    "path/filepath"
    "fmt"
)

// ...existing code...

func addToStartup(appName, exePath string) error {
    u, err := user.Current()
    if err != nil {
        return err
    }
    startupDir := filepath.Join(u.HomeDir, "AppData", "Roaming", "Microsoft", "Windows", "Start Menu", "Programs", "Startup")
    shortcutPath := filepath.Join(startupDir, appName+".bat")
    batContent := fmt.Sprintf(`start "" "%s"`, exePath)
    return os.WriteFile(shortcutPath, []byte(batContent), 0644)
}

//usage: 

err := addToStartup("WingetAutoUpdate", "C:\\Path\\To\\wingetau.exe")
if err != nil {
    logMessage(fmt.Sprintf("Failed to add to startup: %v", err))
}
