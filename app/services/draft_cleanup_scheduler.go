package services

import (
	"log"
	"time"
)

// StartDraftCleanupScheduler starts a background goroutine that periodically cleans up old drafts
// It runs the cleanup once per day at the specified hour (0-23)
// The cleanup deletes all drafts older than the specified number of days
func StartDraftCleanupScheduler(cleanupHour int, daysOld int) {
	if cleanupHour < 0 || cleanupHour > 23 {
		log.Printf("Invalid cleanup hour %d, defaulting to 2 AM", cleanupHour)
		cleanupHour = 2
	}

	if daysOld < 0 {
		log.Printf("Invalid daysOld %d, defaulting to 10 days", daysOld)
		daysOld = 10
	}

	log.Printf("Starting draft cleanup scheduler: runs daily at %02d:00, deletes drafts older than %d days", cleanupHour, daysOld)

	go func() {
		// Calculate the next run time
		now := time.Now()
		nextRun := time.Date(now.Year(), now.Month(), now.Day(), cleanupHour, 0, 0, 0, now.Location())

		// If the scheduled time has already passed today, schedule for tomorrow
		if nextRun.Before(now) || nextRun.Equal(now) {
			nextRun = nextRun.AddDate(0, 0, 1)
		}

		// Wait until the first scheduled time
		durationUntilNextRun := time.Until(nextRun)
		log.Printf("Draft cleanup scheduler: first run scheduled for %s (in %v)", nextRun.Format("2006-01-02 15:04:05"), durationUntilNextRun)
		time.Sleep(durationUntilNextRun)

		// Create a ticker that runs once per day
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		// Run cleanup at the scheduled time
		runCleanup(daysOld)

		// Then run cleanup every 24 hours
		for range ticker.C {
			runCleanup(daysOld)
		}
	}()
}

// runCleanup executes the cleanup and logs the results
func runCleanup(daysOld int) {
	log.Printf("Starting draft cleanup: deleting drafts older than %d days...", daysOld)
	
	startTime := time.Now()
	deletedCount, err := CleanupOldDrafts(daysOld)
	duration := time.Since(startTime)

	if err != nil {
		log.Printf("ERROR: Draft cleanup failed after %v: %v", duration, err)
	} else {
		if deletedCount > 0 {
			log.Printf("✓ Draft cleanup completed in %v: deleted %d draft(s) older than %d days", duration, deletedCount, daysOld)
		} else {
			log.Printf("✓ Draft cleanup completed in %v: no drafts to delete (all drafts are newer than %d days)", duration, daysOld)
		}
	}
}

