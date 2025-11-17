package sync

import (
	"fmt"
	"log"
	"time"

	"github.com/mcdougaj/go-dispatch/internal/database"
	"github.com/mcdougaj/go-dispatch/internal/motive"
)

// Service handles periodic synchronization of data from Motive API
type Service struct {
	motiveClient *motive.Client
	db           *database.DB
	interval     time.Duration
	stopChan     chan struct{}
}

// NewService creates a new sync service
func NewService(motiveClient *motive.Client, db *database.DB, intervalSeconds int) *Service {
	return &Service{
		motiveClient: motiveClient,
		db:           db,
		interval:     time.Duration(intervalSeconds) * time.Second,
		stopChan:     make(chan struct{}),
	}
}

// Start begins the periodic sync process
func (s *Service) Start() {
	log.Printf("Starting sync service with %v interval", s.interval)

	// Run initial sync immediately
	s.syncAll()

	// Start periodic sync
	ticker := time.NewTicker(s.interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				s.syncAll()
			case <-s.stopChan:
				ticker.Stop()
				log.Println("Sync service stopped")
				return
			}
		}
	}()
}

// Stop stops the sync service
func (s *Service) Stop() {
	close(s.stopChan)
}

// syncAll synchronizes all data from Motive API
func (s *Service) syncAll() {
	log.Println("Starting full sync from Motive API")
	startTime := time.Now()

	// Sync in parallel
	errChan := make(chan error, 7)

	go func() { errChan <- s.syncVehicleLocations() }()
	go func() { errChan <- s.syncVehicles() }()
	go func() { errChan <- s.syncAssets() }()
	go func() { errChan <- s.syncDrivers() }()
	go func() { errChan <- s.syncDriverLocations() }()
	go func() { errChan <- s.syncHOS() }()
	go func() { errChan <- s.syncGeofences() }()

	// Collect errors
	for i := 0; i < 7; i++ {
		if err := <-errChan; err != nil {
			log.Printf("Sync error: %v", err)
		}
	}

	duration := time.Since(startTime)
	log.Printf("Full sync completed in %v", duration)
}

// syncVehicleLocations syncs vehicle locations
func (s *Service) syncVehicleLocations() error {
	log.Println("Syncing vehicle locations...")

	resp, err := s.motiveClient.GetVehicleLocations()
	if err != nil {
		s.db.RecordSyncStatus("vehicle_locations", 0, "failed", &err.Error())
		return fmt.Errorf("failed to get vehicle locations: %w", err)
	}

	count := 0
	for _, vehicle := range resp.Vehicles {
		if err := s.db.UpsertVehicleLocation(&vehicle); err != nil {
			log.Printf("Failed to upsert vehicle %d: %v", vehicle.ID, err)
			continue
		}
		count++
	}

	s.db.RecordSyncStatus("vehicle_locations", count, "success", nil)
	log.Printf("Synced %d vehicle locations", count)
	return nil
}

// syncVehicles syncs vehicle details
func (s *Service) syncVehicles() error {
	log.Println("Syncing vehicles...")

	resp, err := s.motiveClient.GetVehicles()
	if err != nil {
		errMsg := err.Error()
		s.db.RecordSyncStatus("vehicles", 0, "failed", &errMsg)
		return fmt.Errorf("failed to get vehicles: %w", err)
	}

	count := 0
	for _, vehicle := range resp.Vehicles {
		if err := s.db.UpsertVehicle(&vehicle); err != nil {
			log.Printf("Failed to upsert vehicle %d: %v", vehicle.ID, err)
			continue
		}
		count++
	}

	s.db.RecordSyncStatus("vehicles", count, "success", nil)
	log.Printf("Synced %d vehicles", count)
	return nil
}

// syncAssets syncs asset locations
func (s *Service) syncAssets() error {
	log.Println("Syncing assets...")

	resp, err := s.motiveClient.GetAssetLocations()
	if err != nil {
		errMsg := err.Error()
		s.db.RecordSyncStatus("assets", 0, "failed", &errMsg)
		return fmt.Errorf("failed to get asset locations: %w", err)
	}

	count := 0
	for _, asset := range resp.Assets {
		if err := s.db.UpsertAssetLocation(&asset); err != nil {
			log.Printf("Failed to upsert asset %d: %v", asset.ID, err)
			continue
		}
		count++
	}

	s.db.RecordSyncStatus("assets", count, "success", nil)
	log.Printf("Synced %d assets", count)
	return nil
}

// syncDrivers syncs driver details
func (s *Service) syncDrivers() error {
	log.Println("Syncing drivers...")

	resp, err := s.motiveClient.GetDrivers()
	if err != nil {
		errMsg := err.Error()
		s.db.RecordSyncStatus("drivers", 0, "failed", &errMsg)
		return fmt.Errorf("failed to get drivers: %w", err)
	}

	count := 0
	for _, driver := range resp.Users {
		if err := s.db.UpsertDriver(&driver); err != nil {
			log.Printf("Failed to upsert driver %d: %v", driver.ID, err)
			continue
		}
		count++
	}

	s.db.RecordSyncStatus("drivers", count, "success", nil)
	log.Printf("Synced %d drivers", count)
	return nil
}

// syncDriverLocations syncs driver locations
func (s *Service) syncDriverLocations() error {
	log.Println("Syncing driver locations...")

	resp, err := s.motiveClient.GetDriverLocations()
	if err != nil {
		errMsg := err.Error()
		s.db.RecordSyncStatus("driver_locations", 0, "failed", &errMsg)
		return fmt.Errorf("failed to get driver locations: %w", err)
	}

	count := 0
	for _, location := range resp.Drivers {
		if err := s.db.UpsertDriverLocation(&location); err != nil {
			log.Printf("Failed to upsert driver location %d: %v", location.DriverID, err)
			continue
		}
		count++
	}

	s.db.RecordSyncStatus("driver_locations", count, "success", nil)
	log.Printf("Synced %d driver locations", count)
	return nil
}

// syncHOS syncs driver HOS data
func (s *Service) syncHOS() error {
	log.Println("Syncing HOS data...")

	resp, err := s.motiveClient.GetHOSAvailableTime()
	if err != nil {
		errMsg := err.Error()
		s.db.RecordSyncStatus("hos", 0, "failed", &errMsg)
		return fmt.Errorf("failed to get HOS data: %w", err)
	}

	count := 0
	for _, hos := range resp.HOSData {
		if err := s.db.UpsertHOSData(&hos); err != nil {
			log.Printf("Failed to upsert HOS for driver %d: %v", hos.DriverID, err)
			continue
		}
		count++
	}

	s.db.RecordSyncStatus("hos", count, "success", nil)
	log.Printf("Synced %d HOS records", count)
	return nil
}

// syncGeofences syncs geofences
func (s *Service) syncGeofences() error {
	log.Println("Syncing geofences...")

	resp, err := s.motiveClient.GetGeofences()
	if err != nil {
		errMsg := err.Error()
		s.db.RecordSyncStatus("geofences", 0, "failed", &errMsg)
		return fmt.Errorf("failed to get geofences: %w", err)
	}

	count := 0
	for _, geofence := range resp.Geofences {
		if err := s.db.UpsertGeofence(&geofence); err != nil {
			log.Printf("Failed to upsert geofence %d: %v", geofence.ID, err)
			continue
		}
		count++
	}

	s.db.RecordSyncStatus("geofences", count, "success", nil)
	log.Printf("Synced %d geofences", count)
	return nil
}
