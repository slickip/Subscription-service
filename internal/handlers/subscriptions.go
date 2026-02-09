package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofrs/uuid"
	"github.com/slickip/Subscription-service/internal/models"
	"gorm.io/gorm"
)

type SubscriptionHandler struct {
	DB *gorm.DB
}

func (h *SubscriptionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	switch {
	case path == "/subscriptions" && r.Method == http.MethodPost:
		h.CreateSubscription(w, r)
	case path == "/subscriptions" && r.Method == http.MethodGet:
		h.ListOfSubscriptions(w, r)
	case path == "/subscriptions/total-cost" && r.Method == http.MethodPost:
		h.TotalCost(w, r)
	case strings.HasPrefix(path, "/subscriptions/") && r.Method == http.MethodGet:
		h.GetSubscriptionByID(w, r)
	case strings.HasPrefix(path, "/subscriptions/") && r.Method == http.MethodPut:
		h.UpdateSubscription(w, r)
	case strings.HasPrefix(path, "/subscriptions/") && r.Method == http.MethodDelete:
		h.DeleteSubscription(w, r)
	default:
		http.Error(w, "Not found", http.StatusNotFound)
	}
}

func (h *SubscriptionHandler) CreateSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		ServiceName string    `json:"service_name"`
		Price       int       `json:"price"`
		UserID      uuid.UUID `json:"user_id"`
		StartDate   string    `json:"start_date"`
		EndDate     *string   `json:"end_date,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	startMonth, startYear, err := parseMonthYear(req.StartDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sub := models.Subscription{
		UserID:      req.UserID,
		ServiceName: req.ServiceName,
		Price:       req.Price,
		StartMonth:  startMonth,
		StartYear:   startYear,
	}

	if req.EndDate != nil {
		endMonth, endYear, err := parseMonthYear(*req.EndDate)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sub.EndMonth = &endMonth
		sub.EndYear = &endYear
	}

	if err := h.DB.Create(&sub).Error; err != nil {
		log.Printf("failed to create subscription: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(sub)
}

func (h *SubscriptionHandler) GetSubscriptionByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/subscriptions/")
	id, err := uuid.FromString(path)
	if err != nil {
		log.Printf("invalid subscription id %q: %v", path, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var sub models.Subscription
	if err := h.DB.First(&sub, "id = ?", id).Error; err != nil {
		log.Printf("subscription not found: %v", err)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(sub)
}

func (h *SubscriptionHandler) ListOfSubscriptions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := h.DB.Model(&models.Subscription{})

	if userID := r.URL.Query().Get("user_id"); userID != "" {
		if id, err := uuid.FromString(userID); err == nil {
			query = query.Where("user_id = ?", id)
		}
	}

	if serviceName := r.URL.Query().Get("service_name"); serviceName != "" {
		query = query.Where("service_name LIKE ?", "%"+serviceName+"%")
	}

	var subs []models.Subscription
	query.Find(&subs)

	json.NewEncoder(w).Encode(subs)
}

func (h *SubscriptionHandler) UpdateSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/subscriptions/")
	id, err := uuid.FromString(path)
	if err != nil {
		log.Printf("invalid subscription id %q: %v", path, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var req struct {
		ServiceName *string `json:"service_name,omitempty"`
		Price       *int    `json:"price,omitempty"`
		StartDate   *string `json:"start_date,omitempty"`
		EndDate     *string `json:"end_date,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("failed to decode update request: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	var sub models.Subscription
	if err := h.DB.First(&sub, "id = ?", id).Error; err != nil {
		log.Printf("subscription not found for update: %v", err)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	if req.ServiceName != nil {
		sub.ServiceName = *req.ServiceName
	}
	if req.Price != nil {
		sub.Price = *req.Price
	}
	if req.StartDate != nil {
		month, year, err := parseMonthYear(*req.StartDate)
		if err != nil {
			log.Printf("invalid start date %q: %v", *req.StartDate, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sub.StartMonth = month
		sub.StartYear = year
	}
	if req.EndDate != nil {
		if *req.EndDate == "" {
			sub.EndMonth = nil
			sub.EndYear = nil
		} else {
			parts := strings.Split(*req.EndDate, "-")
			if len(parts) == 2 {
				month, _ := strconv.Atoi(parts[0])
				year, _ := strconv.Atoi(parts[1])
				sub.EndMonth = &month
				sub.EndYear = &year
			}
		}
	}

	h.DB.Save(&sub)
	json.NewEncoder(w).Encode(sub)
}

func (h *SubscriptionHandler) DeleteSubscription(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/subscriptions/")
	id, err := uuid.FromString(path)
	if err != nil {
		log.Printf("invalid subscription id %q: %v", path, err)
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	result := h.DB.Delete(&models.Subscription{}, "id = ?", id)
	if result.RowsAffected == 0 {
		log.Printf("subscription not found for delete: %v", id)
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SubscriptionHandler) TotalCost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		StartMonth  int        `json:"start_month"`
		StartYear   int        `json:"start_year"`
		EndMonth    int        `json:"end_month"`
		EndYear     int        `json:"end_year"`
		UserID      *uuid.UUID `json:"user_id,omitempty"`
		ServiceName *string    `json:"service_name,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("failed to decode total-cost request: %v", err)
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if req.StartMonth < 1 || req.StartMonth > 12 || req.EndMonth < 1 || req.EndMonth > 12 {
		log.Printf("invalid month range: start_month=%d end_month=%d", req.StartMonth, req.EndMonth)
		http.Error(w, "Invalid month", http.StatusBadRequest)
		return
	}

	query := h.DB.Model(&models.Subscription{})

	if req.UserID != nil {
		query = query.Where("user_id = ?", req.UserID)
	}
	if req.ServiceName != nil {
		query = query.Where("service_name = ?", *req.ServiceName)
	}

	var subs []models.Subscription
	if err := query.Find(&subs).Error; err != nil {
		log.Printf("failed to query subscriptions for total cost: %v", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	total := 0

	for _, sub := range subs {
		subStartY, subStartM := sub.StartYear, sub.StartMonth
		subEndY, subEndM := req.EndYear, req.EndMonth

		if sub.EndYear != nil && sub.EndMonth != nil {
			subEndY = *sub.EndYear
			subEndM = *sub.EndMonth
		}
		startY, startM := maxYearMonth(
			subStartY, subStartM,
			req.StartYear, req.StartMonth,
		)
		endY, endM := minYearMonth(
			subEndY, subEndM,
			req.EndYear, req.EndMonth,
		)
		if startY > endY || (startY == endY && startM > endM) {
			continue
		}

		months := monthsBetween(startM, startY, endM, endY)
		total += months * sub.Price
	}

	json.NewEncoder(w).Encode(map[string]int{
		"total_cost": total,
	})
}
