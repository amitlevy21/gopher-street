package main

import (
	"context"
	"testing"
	"time"
)

func TestBadURI(t *testing.T) {
	defer func() { _ = recover() }()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Microsecond)
	defer cancel()
	openDBWithURI(ctx, "")
	t.Errorf("did not panic")
}

func TestBadPing(t *testing.T) {
	defer func() { _ = recover() }()
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Microsecond)
	defer cancel()
	openDBWithURI(ctx, "mongodb://hangup")
	t.Errorf("did not panic")
}

func TestClosedTooEarly(t *testing.T) {
	defer func() { _ = recover() }()
	ctx, cancel := context.WithCancel(context.Background())
	client := openDBWithURI(ctx, "mongodb://localhost")
	cancel()
	closeDB(ctx, client)
	closeDB(ctx, client)
	t.Errorf("did not panic")
}
