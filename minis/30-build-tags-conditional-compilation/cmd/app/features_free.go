//go:build !premium

package main

import "fmt"

// IsPremiumEnabled returns whether premium features are enabled
func IsPremiumEnabled() bool {
	return false
}

// EnableAdvancedAnalytics is limited in free tier
func EnableAdvancedAnalytics() {
	fmt.Println("📊 Basic Analytics: ENABLED")
	fmt.Println("   - Simple charts")
	fmt.Println("   - Weekly reports")
	fmt.Println("   (Upgrade to Premium for advanced features)")
}

// EnablePrioritySupport is not available in free tier
func EnablePrioritySupport() {
	fmt.Println("📧 Community Support: ENABLED")
	fmt.Println("   - Email support (48hr response)")
	fmt.Println("   - Community forums")
	fmt.Println("   (Upgrade to Premium for priority support)")
}

// GetMaxConcurrentUsers returns the limit for free tier
func GetMaxConcurrentUsers() int {
	return 100 // Free tier: 100 users
}

// GetAPIRateLimit returns the API rate limit for free tier
func GetAPIRateLimit() int {
	return 100 // Free: 100 requests/hour
}

// GetPremiumFeatures returns a description of available features
func GetPremiumFeatures() string {
	return "FREE TIER: Basic analytics, community support, limited exports, 100 users"
}
