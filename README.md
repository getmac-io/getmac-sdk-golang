# getmac-sdk-golang

Golang SDK for [GetMac](https://getmac.io) API.

This SDK provides a convenient way to interact with the GetMac API, allowing you to manage virtual machines, projects, and more.

> **Note:** This is currently a work in progress. More features will be added soon.

## Installation

```bash
go get -u github.com/getmac-io/getmac-sdk-golang
```

## Usage

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/getmac-io/getmac-sdk-golang"
)

func main() {
    client := getmac.NewClient(
        getmac.WithToken("YOUR_API_TOKEN"),
    )

    vmsService := client.VirtualMachines()
    ctx := context.Background()
    projectID := "your_project_id"

    // List VMs
    _, vms, err := vmsService.List(ctx, projectID)
    if err != nil {
        log.Fatal(err)
    }
    for _, vm := range vms {
        fmt.Println(vm.ID, vm.Name)
    }

    // Create VM
    req := &getmac.CreateVirtualMachineRequest{
        Name:   "test",
        Image:  "macos-sequoia",
        Region: "eu-central-ltu-1",
        Type:   "mac-m4-c4-m8",
    }
    _, vm, err := vmsService.Create(ctx, projectID, req)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created VM: %+v\n", vm)
}
```

## API Reference

See [API Reference](https://pkg.go.dev/github.com/getmac-io/getmac-sdk-golang) for detailed documentation.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
