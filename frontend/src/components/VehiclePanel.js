import React from 'react';

function VehiclePanel({ vehicles, selectedVehicle, onSelectVehicle }) {
  return (
    <div className="panel">
      <h2>Vehicles ({vehicles.length})</h2>
      <div className="panel-list">
        {vehicles.length === 0 ? (
          <p className="empty-message">No vehicles found</p>
        ) : (
          vehicles.map((vehicle) => (
            <div
              key={vehicle.VehicleID}
              className={`panel-item ${selectedVehicle?.VehicleID === vehicle.VehicleID ? 'selected' : ''}`}
              onClick={() => onSelectVehicle(vehicle)}
            >
              <div className="panel-item-header">
                <span className="panel-item-title">🚛 {vehicle.VehicleNumber}</span>
                <span className={`panel-item-badge badge-${vehicle.Status?.toLowerCase() || 'unknown'}`}>
                  {vehicle.Status || 'Unknown'}
                </span>
              </div>
              <div className="panel-item-details">
                {vehicle.SpeedMPH !== null && (
                  <div className="panel-item-detail">
                    <span className="detail-label">Speed:</span>
                    <span className="detail-value">{vehicle.SpeedMPH} mph</span>
                  </div>
                )}
                {vehicle.FuelGallons !== null && (
                  <div className="panel-item-detail">
                    <span className="detail-label">Fuel:</span>
                    <span className="detail-value">{vehicle.FuelGallons.toFixed(1)} gal</span>
                  </div>
                )}
                {vehicle.AssignedDriverID && (
                  <div className="panel-item-detail">
                    <span className="detail-label">Driver ID:</span>
                    <span className="detail-value">{vehicle.AssignedDriverID}</span>
                  </div>
                )}
                {vehicle.Make && vehicle.Model && (
                  <div className="panel-item-detail">
                    <span className="detail-label">Model:</span>
                    <span className="detail-value">{vehicle.Make} {vehicle.Model}</span>
                  </div>
                )}
              </div>
            </div>
          ))
        )}
      </div>
    </div>
  );
}

export default VehiclePanel;
