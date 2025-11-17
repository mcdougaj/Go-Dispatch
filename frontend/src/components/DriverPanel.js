import React from 'react';

function DriverPanel({ drivers, selectedDriver, onSelectDriver }) {
  const formatHOS = (minutes) => {
    if (!minutes) return 'N/A';
    const hours = Math.floor(minutes / 60);
    const mins = minutes % 60;
    return `${hours}h ${mins}m`;
  };

  return (
    <div className="panel">
      <h2>Drivers ({drivers.length})</h2>
      <div className="panel-list">
        {drivers.length === 0 ? (
          <p className="empty-message">No drivers found</p>
        ) : (
          drivers.map((driver) => (
            <div
              key={driver.DriverID}
              className={`panel-item ${selectedDriver?.DriverID === driver.DriverID ? 'selected' : ''}`}
              onClick={() => onSelectDriver(driver)}
            >
              <div className="panel-item-header">
                <span className="panel-item-title">👤 Driver {driver.DriverID}</span>
                <span className={`panel-item-badge badge-${driver.HOSStatus?.toLowerCase() || 'unknown'}`}>
                  {driver.HOSStatus || 'Unknown'}
                </span>
              </div>
              <div className="panel-item-details">
                {driver.HOSRemainingMins !== null && (
                  <div className="panel-item-detail">
                    <span className="detail-label">HOS Left:</span>
                    <span className="detail-value">{formatHOS(driver.HOSRemainingMins)}</span>
                  </div>
                )}
                {driver.CurrentVehicleID && (
                  <div className="panel-item-detail">
                    <span className="detail-label">Vehicle:</span>
                    <span className="detail-value">ID {driver.CurrentVehicleID}</span>
                  </div>
                )}
                {driver.DriverCompanyID && (
                  <div className="panel-item-detail">
                    <span className="detail-label">Company ID:</span>
                    <span className="detail-value">{driver.DriverCompanyID}</span>
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

export default DriverPanel;
