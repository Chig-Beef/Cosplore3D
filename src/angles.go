package main

import "math"

func to_radians(angle float64) float64 {
	return angle * math.Pi / 180.0
}

func to_degrees(angle float64) float64 {
	return angle * 180.0 / math.Pi
}

func bound_angle(angle float64) float64 {
	for angle < 0 {
		//fmt.Println(angle)
		angle += 360
	}
	for angle > 360 {
		angle -= 360
	}
	return angle
}
