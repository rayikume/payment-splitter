package services

import "sync"

type SplitService struct {
	mux sync.RWMutex
}
