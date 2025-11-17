import React, { useEffect, useRef, useState } from 'react';
import { Loader } from '@googlemaps/js-api-loader';
import './FleetMap.css';

const GOOGLE_MAPS_API_KEY = process.env.REACT_APP_GOOGLE_MAPS_API_KEY;

function FleetMap({
  vehicles,
  drivers,
  assets,
  selectedVehicle,
  selectedDriver,
  selectedAsset,
  onSelectVehicle,
  onSelectDriver,
  onSelectAsset
}) {
  const mapRef = useRef(null);
  const [map, setMap] = useState(null);
  const [markers, setMarkers] = useState([]);

  // Initialize map
  useEffect(() => {
    if (!GOOGLE_MAPS_API_KEY) {
      console.error('Google Maps API key not found');
      return;
    }

    const loader = new Loader({
      apiKey: GOOGLE_MAPS_API_KEY,
      version: 'weekly',
    });

    loader.load().then(() => {
      const mapInstance = new window.google.maps.Map(mapRef.current, {
        center: { lat: 32.7767, lng: -96.7970 }, // Dallas, TX
        zoom: 6,
        mapTypeControl: true,
        streetViewControl: false,
        fullscreenControl: true,
      });

      setMap(mapInstance);
    });
  }, []);

  // Update markers when fleet data changes
  useEffect(() => {
    if (!map) return;

    // Clear existing markers
    markers.forEach(marker => marker.setMap(null));
    const newMarkers = [];

    // Add vehicle markers
    vehicles.forEach(vehicle => {
      if (vehicle.CurrentLat && vehicle.CurrentLon) {
        const marker = new window.google.maps.Marker({
          position: { lat: vehicle.CurrentLat, lng: vehicle.CurrentLon },
          map: map,
          title: `Vehicle ${vehicle.VehicleNumber}`,
          icon: {
            path: window.google.maps.SymbolPath.FORWARD_CLOSED_ARROW,
            scale: 5,
            fillColor: getVehicleColor(vehicle.Status),
            fillOpacity: 1,
            strokeColor: '#000',
            strokeWeight: 1,
            rotation: vehicle.Heading || 0,
          },
        });

        const infoWindow = new window.google.maps.InfoWindow({
          content: `
            <div class="marker-info">
              <h3>🚛 ${vehicle.VehicleNumber}</h3>
              <p><strong>Status:</strong> ${vehicle.Status || 'Unknown'}</p>
              <p><strong>Speed:</strong> ${vehicle.SpeedMPH || 0} mph</p>
              <p><strong>Fuel:</strong> ${vehicle.FuelGallons ? vehicle.FuelGallons.toFixed(1) : 'N/A'} gal</p>
              ${vehicle.AssignedDriverID ? `<p><strong>Driver:</strong> ${vehicle.AssignedDriverID}</p>` : ''}
            </div>
          `,
        });

        marker.addListener('click', () => {
          infoWindow.open(map, marker);
          onSelectVehicle(vehicle);
        });

        newMarkers.push(marker);
      }
    });

    // Add driver markers (if not in vehicle)
    drivers.forEach(driver => {
      if (driver.CurrentLat && driver.CurrentLon && !driver.CurrentVehicleID) {
        const marker = new window.google.maps.Marker({
          position: { lat: driver.CurrentLat, lng: driver.CurrentLon },
          map: map,
          title: `Driver ${driver.DriverID}`,
          icon: {
            url: 'data:image/svg+xml;charset=UTF-8,' + encodeURIComponent(`
              <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24">
                <circle cx="12" cy="12" r="10" fill="${getDriverColor(driver.HOSStatus)}" stroke="#000" stroke-width="2"/>
              </svg>
            `),
            scaledSize: new window.google.maps.Size(24, 24),
          },
        });

        const infoWindow = new window.google.maps.InfoWindow({
          content: `
            <div class="marker-info">
              <h3>👤 Driver ${driver.DriverID}</h3>
              <p><strong>Status:</strong> ${driver.HOSStatus || 'Unknown'}</p>
              <p><strong>HOS Remaining:</strong> ${driver.HOSRemainingMins ? Math.floor(driver.HOSRemainingMins / 60) + 'h ' + (driver.HOSRemainingMins % 60) + 'm' : 'N/A'}</p>
            </div>
          `,
        });

        marker.addListener('click', () => {
          infoWindow.open(map, marker);
          onSelectDriver(driver);
        });

        newMarkers.push(marker);
      }
    });

    // Add asset markers (trailers/containers)
    assets.forEach(asset => {
      if (asset.CurrentLat && asset.CurrentLon) {
        const marker = new window.google.maps.Marker({
          position: { lat: asset.CurrentLat, lng: asset.CurrentLon },
          map: map,
          title: asset.AssetName,
          icon: {
            url: 'data:image/svg+xml;charset=UTF-8,' + encodeURIComponent(`
              <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 20 20">
                <rect x="2" y="4" width="16" height="12" fill="${getAssetColor(asset.Status)}" stroke="#000" stroke-width="2"/>
              </svg>
            `),
            scaledSize: new window.google.maps.Size(20, 20),
          },
        });

        const infoWindow = new window.google.maps.InfoWindow({
          content: `
            <div class="marker-info">
              <h3>📦 ${asset.AssetName}</h3>
              <p><strong>Type:</strong> ${asset.AssetType || 'Unknown'}</p>
              <p><strong>Status:</strong> ${asset.Status || 'Unknown'}</p>
              ${asset.AttachedToVehicleID ? `<p><strong>Attached to:</strong> Vehicle ${asset.AttachedToVehicleID}</p>` : ''}
            </div>
          `,
        });

        marker.addListener('click', () => {
          infoWindow.open(map, marker);
          onSelectAsset(asset);
        });

        newMarkers.push(marker);
      }
    });

    setMarkers(newMarkers);

    // Fit bounds to show all markers
    if (newMarkers.length > 0) {
      const bounds = new window.google.maps.LatLngBounds();
      newMarkers.forEach(marker => bounds.extend(marker.getPosition()));
      map.fitBounds(bounds);
    }
  }, [map, vehicles, drivers, assets, onSelectVehicle, onSelectDriver, onSelectAsset]);

  // Center on selected item
  useEffect(() => {
    if (!map) return;

    if (selectedVehicle && selectedVehicle.CurrentLat && selectedVehicle.CurrentLon) {
      map.panTo({ lat: selectedVehicle.CurrentLat, lng: selectedVehicle.CurrentLon });
      map.setZoom(12);
    } else if (selectedDriver && selectedDriver.CurrentLat && selectedDriver.CurrentLon) {
      map.panTo({ lat: selectedDriver.CurrentLat, lng: selectedDriver.CurrentLon });
      map.setZoom(12);
    } else if (selectedAsset && selectedAsset.CurrentLat && selectedAsset.CurrentLon) {
      map.panTo({ lat: selectedAsset.CurrentLat, lng: selectedAsset.CurrentLon });
      map.setZoom(12);
    }
  }, [map, selectedVehicle, selectedDriver, selectedAsset]);

  return <div ref={mapRef} className="fleet-map" />;
}

// Helper functions for marker colors
function getVehicleColor(status) {
  switch (status?.toLowerCase()) {
    case 'moving':
      return '#4CAF50'; // Green
    case 'stopped':
      return '#FF9800'; // Orange
    case 'idling':
      return '#FFC107'; // Yellow
    default:
      return '#9E9E9E'; // Gray
  }
}

function getDriverColor(hosStatus) {
  switch (hosStatus?.toLowerCase()) {
    case 'available':
    case 'off_duty':
      return '#4CAF50'; // Green
    case 'driving':
    case 'on_duty':
      return '#2196F3'; // Blue
    case 'sleeper':
      return '#9C27B0'; // Purple
    default:
      return '#9E9E9E'; // Gray
  }
}

function getAssetColor(status) {
  switch (status?.toLowerCase()) {
    case 'loaded':
      return '#F44336'; // Red
    case 'empty':
      return '#4CAF50'; // Green
    default:
      return '#9E9E9E'; // Gray
  }
}

export default FleetMap;
