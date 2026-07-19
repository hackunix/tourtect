package intent

import (
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"github.com/tourtect/backend/internal/intelligence/model"
)

type Router struct{}

func NewRouter() *Router { return &Router{} }

var (
	criticalPatterns = []string{
		"will not let me leave", "won't let me leave", "cannot leave", "can't leave", "not allowed to leave",
		"không cho tôi rời đi", "không cho tôi xuống xe", "bị giữ lại", "weapon", "dao", "súng", "gun",
	}
	coercionPatterns  = []string{"forced to pay", "forcing me to pay", "must pay", "ép trả", "ép tôi trả", "đe dọa trả tiền"}
	pricePatterns     = []string{"price", "cost", "fare", "charge", "pay", "giá", "bao nhiêu", "trả", "vnd", "usd", "đồng"}
	translatePatterns = []string{"translate", "translation", "dịch", "what does", "nghĩa là gì"}
	placePatterns     = []string{"about this place", "about", "around", "nearby", "where", "place", "địa điểm", "ở đâu", "gần đây"}
	emergencyPatterns = []string{"emergency", "sos", "help me now", "khẩn cấp", "cứu tôi"}
	moneyPattern      = regexp.MustCompile(`(?i)([0-9][0-9.,\s]{0,20})\s*(VND|VNĐ|USD|EUR|KRW|CNY|₫|đ|đồng|dollars?)`)
)

func (r *Router) Route(inputType, text string) model.Route {
	lower := strings.ToLower(strings.TrimSpace(text))
	if hasAny(lower, criticalPatterns) {
		return model.Route{Intent: "emergency_help", Confidence: .99, RequiredTools: []string{"evaluate_safety"}, SafetyOverride: true}
	}
	if hasAny(lower, emergencyPatterns) {
		return model.Route{Intent: "emergency_help", Confidence: .9, RequiredTools: []string{"evaluate_safety"}, SafetyOverride: true}
	}
	if hasAny(lower, coercionPatterns) || strings.Contains(lower, "suspicious") || strings.Contains(lower, "scam") || strings.Contains(lower, "lừa đảo") {
		return model.Route{Intent: "safety_assessment", Confidence: .9, RequiredTools: []string{"evaluate_safety"}}
	}
	if inputType == "image_capture" {
		return model.Route{Intent: "menu_or_receipt_analysis", Confidence: .8, RequiredTools: []string{}, MissingFields: []string{"consent_backed_capture_contract", "confirmed_extracted_fields"}}
	}
	if hasAny(lower, translatePatterns) {
		return model.Route{Intent: "translation", Confidence: .9, RequiredTools: []string{"translate_text"}}
	}
	if hasAny(lower, pricePatterns) || moneyPattern.MatchString(text) {
		missing := []string{}
		candidate, ok := ExtractPriceCandidate(text)
		if !ok {
			missing = append(missing, "amount", "currency")
		}
		if ok && candidate.Unit == "" {
			missing = append(missing, "unit")
		}
		return model.Route{Intent: "price_check", Confidence: .93, RequiredTools: []string{"retrieve_place_context", "evaluate_price"}, MissingFields: missing}
	}
	if hasAny(lower, placePatterns) {
		return model.Route{Intent: "place_information", Confidence: .78, RequiredTools: []string{"retrieve_place_context"}}
	}
	if lower == "" {
		return model.Route{Intent: "unknown", Confidence: 0, MissingFields: []string{"text"}}
	}
	return model.Route{Intent: "general_travel_question", Confidence: .55, RequiredTools: []string{"retrieve_place_context"}}
}

type PriceCandidate struct {
	AmountMinor string `json:"amount_minor"`
	Currency    string `json:"currency"`
	Exponent    int    `json:"exponent"`
	Unit        string `json:"unit"`
	Vertical    string `json:"vertical"`
	RawItem     string `json:"raw_item"`
}

func ExtractPriceCandidate(text string) (PriceCandidate, bool) {
	m := moneyPattern.FindStringSubmatch(text)
	if len(m) == 0 {
		return PriceCandidate{}, false
	}
	currency := strings.ToUpper(m[2])
	switch currency {
	case "VNĐ", "₫", "Đ", "ĐỒNG":
		currency = "VND"
	case "DOLLAR", "DOLLARS":
		currency = "USD"
	}
	exponent := 0
	if currency == "USD" || currency == "EUR" {
		exponent = 2
	}
	digits := strings.Map(func(r rune) rune {
		if unicode.IsDigit(r) {
			return r
		}
		return -1
	}, m[1])
	if digits == "" {
		return PriceCandidate{}, false
	}
	amount, err := strconv.ParseInt(digits, 10, 64)
	if err != nil {
		return PriceCandidate{}, false
	}
	if exponent > 0 && !strings.ContainsAny(m[1], ".,") {
		amount *= 100
	}
	lower := strings.ToLower(text)
	vertical, unit, raw := "street_retail", "item", strings.TrimSpace(text)
	if strings.Contains(lower, "taxi") || strings.Contains(lower, "driver") || strings.Contains(lower, "ride") || strings.Contains(lower, "xe") {
		vertical, unit, raw = "taxi", "trip", "Airport taxi to Old Quarter"
	} else if strings.Contains(lower, "exchange") || strings.Contains(lower, "đổi tiền") {
		vertical, unit = "exchange", "transaction"
	} else if strings.Contains(lower, "menu") || strings.Contains(lower, "food") || strings.Contains(lower, "món") {
		vertical, unit = "food", "item"
	}
	return PriceCandidate{AmountMinor: strconv.FormatInt(amount, 10), Currency: currency, Exponent: exponent, Unit: unit, Vertical: vertical, RawItem: raw}, true
}

type SafetyFacts struct {
	ObservedFacts         []string `json:"observed_facts"`
	ThreatIndicators      []string `json:"threat_indicators"`
	ConfinementIndicators []string `json:"confinement_indicators"`
	CoercionIndicators    []string `json:"coercion_indicators"`
	AbilityToLeave        *bool    `json:"ability_to_leave,omitempty"`
}

func ExtractSafetyFacts(text string) SafetyFacts {
	lower := strings.ToLower(text)
	result := SafetyFacts{ObservedFacts: []string{}, ThreatIndicators: []string{}, ConfinementIndicators: []string{}, CoercionIndicators: []string{}}
	if hasAny(lower, criticalPatterns) {
		canLeave := false
		result.AbilityToLeave = &canLeave
		result.ObservedFacts = append(result.ObservedFacts, "unable_to_leave")
		result.ConfinementIndicators = append(result.ConfinementIndicators, "confinement")
	}
	if strings.Contains(lower, "weapon") || strings.Contains(lower, "gun") || strings.Contains(lower, "dao") || strings.Contains(lower, "súng") {
		result.ThreatIndicators = append(result.ThreatIndicators, "weapon")
	}
	if hasAny(lower, coercionPatterns) {
		result.CoercionIndicators = append(result.CoercionIndicators, "forced_payment")
	}
	if hasAny(lower, pricePatterns) {
		result.ObservedFacts = append(result.ObservedFacts, "price_dispute")
	}
	if len(result.ObservedFacts) == 0 {
		result.ObservedFacts = append(result.ObservedFacts, "informational_safety_question")
	}
	return result
}

func hasAny(value string, patterns []string) bool {
	for _, pattern := range patterns {
		if strings.Contains(value, pattern) {
			return true
		}
	}
	return false
}
