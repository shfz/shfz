/*
Copyright © 2022 shfz

*/
package run

import (
	"context"
	"errors"
	"log"
	"os/exec"
	"time"

	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/semaphore"
)

func Run(target string, number int, parallel int, timeout int) error {
	// https://zenn.dev/aikizoku/articles/golang-goroutine
	ctx := context.Background()
	eg := errgroup.Group{}
	sem := semaphore.NewWeighted(int64(parallel))
	for i := 0; i < number; i++ {
		if err := sem.Acquire(ctx, 1); err != nil {
			return err
		}
		eg.Go(func() error {
			if err := ExecCommand(target, timeout); err != nil {
				return err
			}
			sem.Release(1)
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func ExecCommand(target string, timeout int) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	defer cancel()

	_, err := exec.CommandContext(ctx, "node", target).Output()

	if ctx.Err() == context.DeadlineExceeded {
		return errors.New("ExecCommand Timeout")
	}

	if err != nil {
		log.Println("[i] scenario error :", err)
	}
	return nil
}
