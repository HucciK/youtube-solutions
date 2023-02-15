package resources

import "yt-solutions-telegram-dashboard/internal/core"

type AvailabilityResponse struct {
	Free     int `json:"free"`
	Lifetime int `json:"lifetime"`
	Price    int `json:"price"`
}

type ProcessOrderResponse struct {
	Success bool     `json:"success"`
	Key     core.Key `json:"key"`
}

type RenewalInfoResponse struct {
	Renewal int `json:"renewal"`
}

type ProcessRenewalResponse struct {
	Success bool `json:"success"`
}

type CreateInvoiceResponse struct {
	ID  string `json:"id"`
	Url string `json:"url"`
}

type InvoiceStatusResponse struct {
	Amount int    `json:"amount"`
	State  string `json:"state"`
	Extra  string `json:"extra"`
}
