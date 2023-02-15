package handler

import "yt-solutions-server/internal/core"

type AvailabilityResponse struct {
	Free     int `json:"free"`
	Lifetime int `json:"lifetime"`
	Price    int `json:"price"`
}

type ProcessOrderResponse struct {
	Success bool      `json:"success"`
	Key     *core.Key `json:"key,omitempty"`
}

type RenewalInfoResponse struct {
	Renewal int `json:"renewal"`
}

type ProcessRenewalResponse struct {
	Success bool `json:"success"`
}

type KeyAuthResponse struct {
	Valid         bool      `json:"valid"`
	LatestVersion string    `json:"latest_version"`
	Key           *core.Key `json:"key_info,omitempty"`
}
