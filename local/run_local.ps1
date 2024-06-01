$buildDir="$PSScriptRoot\.build"
Invoke-Expression "go build -o $buildDir\go-wspubsub.exe $PSScriptRoot\..\main.go"