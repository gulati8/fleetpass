import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../services/api';

function VehicleProfile() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [vehicle, setVehicle] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [currentImageIndex, setCurrentImageIndex] = useState(0);

  useEffect(() => {
    fetchVehicle();
  }, [id]);

  const fetchVehicle = async () => {
    try {
      const response = await api.get(`/api/vehicles/${id}`);
      setVehicle(response.data);
    } catch (err) {
      setError('Failed to fetch vehicle details');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async () => {
    if (!window.confirm('Are you sure you want to delete this vehicle?')) {
      return;
    }

    try {
      await api.delete(`/api/vehicles/${id}`);
      navigate('/vehicles');
    } catch (err) {
      setError('Failed to delete vehicle');
    }
  };

  const getStatusBadge = (status) => {
    const badges = {
      available: 'bg-success',
      rented: 'bg-warning text-dark',
      maintenance: 'bg-danger',
      inactive: 'bg-secondary'
    };
    return badges[status] || 'bg-secondary';
  };

  if (loading) {
    return (
      <div className="d-flex justify-content-center align-items-center" style={{ minHeight: '50vh' }}>
        <div className="spinner-border text-primary" role="status">
          <span className="visually-hidden">Loading...</span>
        </div>
      </div>
    );
  }

  if (error || !vehicle) {
    return (
      <div className="container px-5 py-5">
        <div className="alert alert-danger" role="alert">
          {error || 'Vehicle not found'}
        </div>
        <button className="btn btn-primary" onClick={() => navigate('/vehicles')}>
          Back to Vehicles
        </button>
      </div>
    );
  }

  const images = vehicle.images && vehicle.images.length > 0
    ? vehicle.images
    : ['https://via.placeholder.com/800x600?text=No+Image+Available'];

  return (
    <main className="flex-shrink-0">
      <section className="py-5">
        <div className="container px-5">
          {/* Header */}
          <div className="row mb-4">
            <div className="col-lg-8">
              <h1 className="fw-bolder mb-2">
                {vehicle.year} {vehicle.make} {vehicle.model}
                {vehicle.trim && <span className="text-muted"> {vehicle.trim}</span>}
              </h1>
              <div className="mb-3">
                <span className={`badge ${getStatusBadge(vehicle.status)} me-2`}>
                  {vehicle.status}
                </span>
                {vehicle.condition && (
                  <span className="badge bg-info text-dark">
                    {vehicle.condition.replace('_', ' ')}
                  </span>
                )}
              </div>
            </div>
            <div className="col-lg-4 text-lg-end">
              <button
                className="btn btn-outline-primary me-2"
                onClick={() => navigate(`/vehicles/${id}/edit`)}
              >
                Edit Vehicle
              </button>
              <button
                className="btn btn-outline-danger"
                onClick={handleDelete}
              >
                Delete
              </button>
            </div>
          </div>

          <div className="row gx-5">
            {/* Left Column - Images */}
            <div className="col-lg-8">
              {/* Image Gallery */}
              <div className="card mb-4">
                <div className="card-body p-0">
                  <div id="vehicleCarousel" className="carousel slide" data-bs-ride="carousel">
                    <div className="carousel-inner">
                      {images.map((image, index) => (
                        <div
                          key={index}
                          className={`carousel-item ${index === currentImageIndex ? 'active' : ''}`}
                        >
                          <img
                            src={image}
                            className="d-block w-100"
                            alt={`${vehicle.make} ${vehicle.model}`}
                            style={{ height: '500px', objectFit: 'cover' }}
                          />
                        </div>
                      ))}
                    </div>
                    {images.length > 1 && (
                      <>
                        <button
                          className="carousel-control-prev"
                          type="button"
                          data-bs-target="#vehicleCarousel"
                          data-bs-slide="prev"
                          onClick={() => setCurrentImageIndex((currentImageIndex - 1 + images.length) % images.length)}
                        >
                          <span className="carousel-control-prev-icon" aria-hidden="true"></span>
                          <span className="visually-hidden">Previous</span>
                        </button>
                        <button
                          className="carousel-control-next"
                          type="button"
                          data-bs-target="#vehicleCarousel"
                          data-bs-slide="next"
                          onClick={() => setCurrentImageIndex((currentImageIndex + 1) % images.length)}
                        >
                          <span className="carousel-control-next-icon" aria-hidden="true"></span>
                          <span className="visually-hidden">Next</span>
                        </button>
                      </>
                    )}
                  </div>
                </div>
              </div>

              {/* Vehicle Description */}
              {vehicle.description && (
                <div className="card mb-4">
                  <div className="card-body">
                    <h4 className="card-title mb-3">Description</h4>
                    <p className="card-text">{vehicle.description}</p>
                  </div>
                </div>
              )}

              {/* Vehicle Specifications */}
              <div className="card mb-4">
                <div className="card-body">
                  <h4 className="card-title mb-3">Vehicle Specifications</h4>
                  <div className="row">
                    <div className="col-md-6 mb-3">
                      <h6 className="text-muted">VIN</h6>
                      <p><code>{vehicle.vin}</code></p>
                    </div>
                    <div className="col-md-6 mb-3">
                      <h6 className="text-muted">Stock Number</h6>
                      <p>{vehicle.stock_number || 'N/A'}</p>
                    </div>
                    <div className="col-md-6 mb-3">
                      <h6 className="text-muted">Body Style</h6>
                      <p>{vehicle.body_style || 'N/A'}</p>
                    </div>
                    <div className="col-md-6 mb-3">
                      <h6 className="text-muted">Transmission</h6>
                      <p>{vehicle.transmission || 'N/A'}</p>
                    </div>
                    <div className="col-md-6 mb-3">
                      <h6 className="text-muted">Drivetrain</h6>
                      <p>{vehicle.drivetrain || 'N/A'}</p>
                    </div>
                    <div className="col-md-6 mb-3">
                      <h6 className="text-muted">Fuel Type</h6>
                      <p>{vehicle.fuel_type || 'N/A'}</p>
                    </div>
                    <div className="col-md-6 mb-3">
                      <h6 className="text-muted">Engine</h6>
                      <p>{vehicle.engine || 'N/A'}</p>
                    </div>
                    <div className="col-md-6 mb-3">
                      <h6 className="text-muted">MPG</h6>
                      <p>
                        {vehicle.mpg_city && vehicle.mpg_highway
                          ? `${vehicle.mpg_city} city / ${vehicle.mpg_highway} highway`
                          : 'N/A'}
                      </p>
                    </div>
                    <div className="col-md-6 mb-3">
                      <h6 className="text-muted">Exterior Color</h6>
                      <p>{vehicle.color_exterior || 'N/A'}</p>
                    </div>
                    <div className="col-md-6 mb-3">
                      <h6 className="text-muted">Interior Color</h6>
                      <p>{vehicle.color_interior || 'N/A'}</p>
                    </div>
                    <div className="col-md-6 mb-3">
                      <h6 className="text-muted">Seats</h6>
                      <p>{vehicle.seats || 'N/A'}</p>
                    </div>
                    <div className="col-md-6 mb-3">
                      <h6 className="text-muted">Doors</h6>
                      <p>{vehicle.doors || 'N/A'}</p>
                    </div>
                    {vehicle.license_plate && (
                      <div className="col-md-6 mb-3">
                        <h6 className="text-muted">License Plate</h6>
                        <p>{vehicle.license_plate}</p>
                      </div>
                    )}
                  </div>
                </div>
              </div>

              {/* Features */}
              {vehicle.features && vehicle.features.length > 0 && (
                <div className="card mb-4">
                  <div className="card-body">
                    <h4 className="card-title mb-3">Features & Options</h4>
                    <div className="row">
                      {vehicle.features.map((feature, index) => (
                        <div key={index} className="col-md-6 mb-2">
                          <i className="bi bi-check-circle-fill text-success me-2"></i>
                          {feature}
                        </div>
                      ))}
                    </div>
                  </div>
                </div>
              )}
            </div>

            {/* Right Column - Pricing & Key Info */}
            <div className="col-lg-4">
              {/* Pricing Card */}
              <div className="card mb-4 border-primary">
                <div className="card-body">
                  <h4 className="card-title mb-3">Rental Pricing</h4>
                  {vehicle.daily_rate > 0 && (
                    <div className="mb-3">
                      <h6 className="text-muted mb-1">Daily Rate</h6>
                      <h3 className="text-primary mb-0">${vehicle.daily_rate.toFixed(2)}</h3>
                      <small className="text-muted">per day</small>
                    </div>
                  )}
                  {vehicle.weekly_rate > 0 && (
                    <div className="mb-3">
                      <h6 className="text-muted mb-1">Weekly Rate</h6>
                      <h4 className="mb-0">${vehicle.weekly_rate.toFixed(2)}</h4>
                      <small className="text-muted">per week</small>
                    </div>
                  )}
                  {vehicle.monthly_rate > 0 && (
                    <div className="mb-3">
                      <h6 className="text-muted mb-1">Monthly Rate</h6>
                      <h4 className="mb-0">${vehicle.monthly_rate.toFixed(2)}</h4>
                      <small className="text-muted">per month</small>
                    </div>
                  )}
                  {vehicle.status === 'available' && (
                    <button className="btn btn-primary w-100 mt-3">
                      Reserve Now
                    </button>
                  )}
                </div>
              </div>

              {/* Key Information Card */}
              <div className="card mb-4">
                <div className="card-body">
                  <h5 className="card-title mb-3">Key Information</h5>
                  <div className="mb-3">
                    <h6 className="text-muted mb-1">Mileage</h6>
                    <p className="mb-0 fw-bold">
                      {vehicle.mileage ? vehicle.mileage.toLocaleString() : 'N/A'} miles
                    </p>
                  </div>
                  <div className="mb-3">
                    <h6 className="text-muted mb-1">Condition</h6>
                    <p className="mb-0 fw-bold text-capitalize">
                      {vehicle.condition ? vehicle.condition.replace('_', ' ') : 'N/A'}
                    </p>
                  </div>
                  <div className="mb-3">
                    <h6 className="text-muted mb-1">Eligible for Service</h6>
                    <p className="mb-0">
                      {vehicle.is_eligible_for_service ? (
                        <span className="badge bg-success">Yes</span>
                      ) : (
                        <span className="badge bg-warning text-dark">No</span>
                      )}
                    </p>
                  </div>
                  {vehicle.warranty_type && (
                    <div className="mb-3">
                      <h6 className="text-muted mb-1">Warranty</h6>
                      <p className="mb-0">{vehicle.warranty_type}</p>
                      {vehicle.warranty_details && (
                        <small className="text-muted">{vehicle.warranty_details}</small>
                      )}
                    </div>
                  )}
                </div>
              </div>

              {/* Back Button */}
              <button
                className="btn btn-outline-secondary w-100"
                onClick={() => navigate('/vehicles')}
              >
                Back to Vehicles
              </button>
            </div>
          </div>
        </div>
      </section>
    </main>
  );
}

export default VehicleProfile;
