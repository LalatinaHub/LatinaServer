package config

import (
	"context"
	"fmt"
	"sync"
)

// GenerateConfigsParallel generates both Caddy and Sing-box configurations concurrently.
// It returns an error if any of the generators fail.
func GenerateConfigsParallel(ctx context.Context) error {
	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := GenerateCaddyConfig(); err != nil {
			errChan <- fmt.Errorf("GenerateCaddyConfig failed: %w", err)
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := GenerateSingConfig(); err != nil {
			errChan <- fmt.Errorf("GenerateSingConfig failed: %w", err)
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}
