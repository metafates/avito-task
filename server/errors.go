package server

import "errors"

var (
	errUserNotFound           = errors.New("user not found")
	errUserExists             = errors.New("user exists")
	errSegmentNotFound        = errors.New("segment not found")
	errSegmentExists          = errors.New("segment exists")
	errSegmentAssignedAlready = errors.New("segment is already assigned")
	errSegmentNotAssigned     = errors.New("segment is not assigned")
)
