import React from 'react';

function AssetPanel({ assets, selectedAsset, onSelectAsset }) {
  return (
    <div className="panel">
      <h2>Assets ({assets.length})</h2>
      <div className="panel-list">
        {assets.length === 0 ? (
          <p className="empty-message">No assets found</p>
        ) : (
          assets.map((asset) => (
            <div
              key={asset.AssetID}
              className={`panel-item ${selectedAsset?.AssetID === asset.AssetID ? 'selected' : ''}`}
              onClick={() => onSelectAsset(asset)}
            >
              <div className="panel-item-header">
                <span className="panel-item-title">📦 {asset.AssetName}</span>
                <span className={`panel-item-badge badge-${asset.Status?.toLowerCase() || 'unknown'}`}>
                  {asset.Status || 'Unknown'}
                </span>
              </div>
              <div className="panel-item-details">
                {asset.AssetType && (
                  <div className="panel-item-detail">
                    <span className="detail-label">Type:</span>
                    <span className="detail-value">{asset.AssetType}</span>
                  </div>
                )}
                {asset.AttachedToVehicleID && (
                  <div className="panel-item-detail">
                    <span className="detail-label">Attached To:</span>
                    <span className="detail-value">Vehicle {asset.AttachedToVehicleID}</span>
                  </div>
                )}
                {asset.LocationDescription && (
                  <div className="panel-item-detail">
                    <span className="detail-label">Location:</span>
                    <span className="detail-value">{asset.LocationDescription}</span>
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

export default AssetPanel;
