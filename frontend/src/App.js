import React, { useState, useEffect, useCallback } from 'react';
import './App.css';
import FleetMap from './components/FleetMap';
import QueryBar from './components/QueryBar';
import VehiclePanel from './components/VehiclePanel';
import DriverPanel from './components/DriverPanel';
import AssetPanel from './components/AssetPanel';
import { fetchVehicles, fetchDrivers, fetchAssets, sendQuery } from './services/api';

function App() {
  const [vehicles, setVehicles] = useState([]);
  const [drivers, setDrivers] = useState([]);
  const [assets, setAssets] = useState([]);
  const [selectedVehicle, setSelectedVehicle] = useState(null);
  const [selectedDriver, setSelectedDriver] = useState(null);
  const [selectedAsset, setSelectedAsset] = useState(null);
  const [queryResponse, setQueryResponse] = useState(null);
  const [loading, setLoading] = useState(true);

  // Load fleet data
  const loadFleetData = useCallback(async () => {
    try {
      const [vehiclesData, driversData, assetsData] = await Promise.all([
        fetchVehicles(),
        fetchDrivers(),
        fetchAssets()
      ]);

      setVehicles(vehiclesData);
      setDrivers(driversData);
      setAssets(assetsData);
      setLoading(false);
    } catch (error) {
      console.error('Error loading fleet data:', error);
      setLoading(false);
    }
  }, []);

  // Initial load
  useEffect(() => {
    loadFleetData();

    // Refresh data every 30 seconds
    const interval = setInterval(loadFleetData, 30000);
    return () => clearInterval(interval);
  }, [loadFleetData]);

  // Handle natural language query
  const handleQuery = async (query) => {
    try {
      const response = await sendQuery(query);
      setQueryResponse(response);
    } catch (error) {
      console.error('Error processing query:', error);
      setQueryResponse({
        response: 'Sorry, I encountered an error processing your query. Please try again.'
      });
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <h1>🚛 Motive Dispatch - Fleet Tracking</h1>
        <div className="fleet-stats">
          <span className="stat">
            <strong>{vehicles.length}</strong> Vehicles
          </span>
          <span className="stat">
            <strong>{drivers.length}</strong> Drivers
          </span>
          <span className="stat">
            <strong>{assets.length}</strong> Assets
          </span>
        </div>
      </header>

      <QueryBar onQuery={handleQuery} queryResponse={queryResponse} />

      <div className="dashboard">
        <aside className="sidebar left">
          <VehiclePanel
            vehicles={vehicles}
            selectedVehicle={selectedVehicle}
            onSelectVehicle={setSelectedVehicle}
          />
        </aside>

        <main className="map-container">
          {loading ? (
            <div className="loading">Loading fleet data...</div>
          ) : (
            <FleetMap
              vehicles={vehicles}
              drivers={drivers}
              assets={assets}
              selectedVehicle={selectedVehicle}
              selectedDriver={selectedDriver}
              selectedAsset={selectedAsset}
              onSelectVehicle={setSelectedVehicle}
              onSelectDriver={setSelectedDriver}
              onSelectAsset={setSelectedAsset}
            />
          )}
        </main>

        <aside className="sidebar right">
          <DriverPanel
            drivers={drivers}
            selectedDriver={selectedDriver}
            onSelectDriver={setSelectedDriver}
          />
          <AssetPanel
            assets={assets}
            selectedAsset={selectedAsset}
            onSelectAsset={setSelectedAsset}
          />
        </aside>
      </div>
    </div>
  );
}

export default App;
