package model

import "container/list"

type Trip struct {
	ID          string     `json:"id"`
	RouteNumber string     `json:"route_number"`
	Direction   string     `json:"direction"`
	Stops       *list.List `json:"-"`
}
