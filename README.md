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

    // Create VM from a runner label. The label selects the image and the
    // machine type; set Type only to override the machine type.
    req := &getmac.CreateVirtualMachineRequest{
        Name:   "test",
        Image:  "getmac-tahoe",
        Region: "eu-central-ltu-1",
    }
    _, vm, err := vmsService.Create(ctx, projectID, req)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created VM: %+v\n", vm)
}
```

## Images and Machine Types

`CreateVirtualMachineRequest.Image` accepts either:

- a **GetMac runner label**, the same names you use in GitHub Actions `runs-on`, such as `getmac`, `getmac-tahoe` or `getmac-sequoia`. GetMac maps each label to an image and a machine type, so `Type` can be left empty.
- an **image slug**, such as `macos-tahoe` or `macos-sequoia`. An image slug doesn't imply a machine type, so `Type` is required (for example `mac-m4-c4-m8`).

A `Type` you set always takes precedence over the label's machine type.

> **Note:** Labels in `Image` require a GetMac API version with label support. Older API versions accept only image slugs and always require `Type`.

## Errors

When the API responds with an unexpected status code, methods return a `*getmac.APIError` with the status code and the reason the API gave:

```go
_, _, err := vmsService.Create(ctx, projectID, req)

var apiErr *getmac.APIError
if errors.As(err, &apiErr) {
    fmt.Println(apiErr.StatusCode, apiErr.Message) // 400 Image or label getmac-nope not found
}
```

Errors for resources that don't exist wrap `getmac.ErrNotFound`. That covers a `404` response and `GetByName` finding no virtual machine with the name:

```go
_, vm, err := vmsService.GetByName(ctx, projectID, "my-vm")
switch {
case errors.Is(err, getmac.ErrNotFound):
    // no virtual machine named "my-vm"
case err != nil:
    log.Fatal(err)
default:
    fmt.Println(vm.ID)
}
```

## API Reference

See [API Reference](https://pkg.go.dev/github.com/getmac-io/getmac-sdk-golang) for detailed documentation.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
