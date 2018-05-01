// +build !windows

package main

func QPC() int64          { return 0 }
func QPCFrequency() int64 { return 1 }
