package main

import (
	"context"
	"fmt"
	"log"

	leap0 "github.com/leap0dev/leap0-go"
)

func main() {
	client, err := leap0.NewClient()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	ctx := context.Background()

	sb, err := client.Sandboxes.Create(ctx, &leap0.CreateSandboxParams{
		TemplateName: leap0.DefaultTemplate,
		VCPU:         2,
		Memory:       2048,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer client.Sandboxes.Delete(ctx, sb.ID)

	fmt.Printf("Sandbox: %s (%s)\n", sb.ID, sb.State)

	result, err := client.Process.Execute(ctx, sb.ID, &leap0.ExecParams{
		Command: "echo 'Hello from Leap0!'",
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Output: %s", result.Stdout)

	_ = client.Filesystem.WriteFile(ctx, sb.ID, "/tmp/hello.txt", "Hello!", "")
	content, _ := client.Filesystem.ReadFile(ctx, sb.ID, "/tmp/hello.txt", nil)
	fmt.Printf("File: %s\n", content)
}
