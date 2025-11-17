import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import api from '../services/api';

function VehicleForm() {
  const { id } = useParams();
  const navigate = useNavigate();
  const isEditMode = !!id;

  const [locations, setLocations] = useState([]);
  const [loading, setLoading] = useState(isEditMode);
  const [error, setError] = useState('');
  const [formData, setFormData] = useState({
    location_id: '',
    vin: '',
    make: '',
    model: '',
    year: new Date().getFullYear(),
    trim: '',
    color_exterior: '',
    color_interior: '',
    condition: 'used',
    mileage: 0,
    license_plate: '',
    status: 'available',
    is_eligible_for_service: true,
    body_style: '',
    transmission: '',
    drivetrain: '',
    fuel_type: '',
    engine: '',
    mpg_city: 0,
    mpg_highway: 0,
    seats: 5,
    doors: 4,
    stock_number: '',
    description: '',
    daily_rate: 0,
    weekly_rate: 0,
    monthly_rate: 0,
    features: '',
    images: ''
  });

  useEffect(() => {
    fetchLocations();
    if (isEditMode) {
      fetchVehicle();
    }
  }, [id]);

  const fetchLocations = async () => {
    try {
      const response = await api.get('/api/locations');
      setLocations(response.data || []);
    } catch (err) {
      setError('Failed to fetch locations');
    }
  };

  const fetchVehicle = async () => {
    try {
      const response = await api.get(`/api/vehicles/${id}`);
      const vehicle = response.data;
      setFormData({
        location_id: vehicle.location_id || '',
        vin: vehicle.vin || '',
        make: vehicle.make || '',
        model: vehicle.model || '',
        year: vehicle.year || new Date().getFullYear(),
        trim: vehicle.trim || '',
        color_exterior: vehicle.color_exterior || '',
        color_interior: vehicle.color_interior || '',
        condition: vehicle.condition || 'used',
        mileage: vehicle.mileage || 0,
        license_plate: vehicle.license_plate || '',
        status: vehicle.status || 'available',
        is_eligible_for_service: vehicle.is_eligible_for_service !== false,
        body_style: vehicle.body_style || '',
        transmission: vehicle.transmission || '',
        drivetrain: vehicle.drivetrain || '',
        fuel_type: vehicle.fuel_type || '',
        engine: vehicle.engine || '',
        mpg_city: vehicle.mpg_city || 0,
        mpg_highway: vehicle.mpg_highway || 0,
        seats: vehicle.seats || 5,
        doors: vehicle.doors || 4,
        stock_number: vehicle.stock_number || '',
        description: vehicle.description || '',
        daily_rate: vehicle.daily_rate || 0,
        weekly_rate: vehicle.weekly_rate || 0,
        monthly_rate: vehicle.monthly_rate || 0,
        features: vehicle.features ? vehicle.features.join('\n') : '',
        images: vehicle.images ? vehicle.images.join('\n') : ''
      });
    } catch (err) {
      setError('Failed to fetch vehicle');
    } finally {
      setLoading(false);
    }
  };

  const handleChange = (e) => {
    const { name, value, type, checked } = e.target;
    setFormData({
      ...formData,
      [name]: type === 'checkbox' ? checked : value
    });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    // Convert features and images from text to arrays
    const payload = {
      ...formData,
      year: parseInt(formData.year),
      mileage: parseInt(formData.mileage),
      mpg_city: parseInt(formData.mpg_city),
      mpg_highway: parseInt(formData.mpg_highway),
      seats: parseInt(formData.seats),
      doors: parseInt(formData.doors),
      daily_rate: parseFloat(formData.daily_rate),
      weekly_rate: parseFloat(formData.weekly_rate),
      monthly_rate: parseFloat(formData.monthly_rate),
      features: formData.features
        ? formData.features.split('\n').map(f => f.trim()).filter(f => f)
        : [],
      images: formData.images
        ? formData.images.split('\n').map(i => i.trim()).filter(i => i)
        : []
    };

    try {
      if (isEditMode) {
        await api.put(`/api/vehicles/${id}`, payload);
      } else {
        await api.post('/api/vehicles', payload);
      }
      navigate('/vehicles');
    } catch (err) {
      setError(err.response?.data || `Failed to ${isEditMode ? 'update' : 'create'} vehicle`);
    }
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

  return (
    <main className="flex-shrink-0">
      <section className="py-5">
        <div className="container px-5">
          <div className="row justify-content-center">
            <div className="col-lg-10">
              <h1 className="fw-bolder mb-4">
                {isEditMode ? 'Edit Vehicle' : 'Add New Vehicle'}
              </h1>

              {error && (
                <div className="alert alert-danger" role="alert">
                  {error}
                </div>
              )}

              <form onSubmit={handleSubmit}>
                {/* Basic Information */}
                <div className="card mb-4">
                  <div className="card-header">
                    <h5 className="mb-0">Basic Information</h5>
                  </div>
                  <div className="card-body">
                    <div className="row">
                      <div className="col-md-6 mb-3">
                        <label className="form-label">Location *</label>
                        <select
                          className="form-select"
                          name="location_id"
                          value={formData.location_id}
                          onChange={handleChange}
                          required
                        >
                          <option value="">Select a location</option>
                          {locations.map((loc) => (
                            <option key={loc.id} value={loc.id}>
                              {loc.name} - {loc.address}
                            </option>
                          ))}
                        </select>
                      </div>
                      <div className="col-md-6 mb-3">
                        <label className="form-label">VIN *</label>
                        <input
                          type="text"
                          className="form-control"
                          name="vin"
                          value={formData.vin}
                          onChange={handleChange}
                          required
                        />
                      </div>
                      <div className="col-md-4 mb-3">
                        <label className="form-label">Make *</label>
                        <input
                          type="text"
                          className="form-control"
                          name="make"
                          value={formData.make}
                          onChange={handleChange}
                          placeholder="e.g., Toyota"
                          required
                        />
                      </div>
                      <div className="col-md-4 mb-3">
                        <label className="form-label">Model *</label>
                        <input
                          type="text"
                          className="form-control"
                          name="model"
                          value={formData.model}
                          onChange={handleChange}
                          placeholder="e.g., Camry"
                          required
                        />
                      </div>
                      <div className="col-md-4 mb-3">
                        <label className="form-label">Year *</label>
                        <input
                          type="number"
                          className="form-control"
                          name="year"
                          value={formData.year}
                          onChange={handleChange}
                          min="1900"
                          max={new Date().getFullYear() + 1}
                          required
                        />
                      </div>
                      <div className="col-md-4 mb-3">
                        <label className="form-label">Trim</label>
                        <input
                          type="text"
                          className="form-control"
                          name="trim"
                          value={formData.trim}
                          onChange={handleChange}
                          placeholder="e.g., XLE, Sport"
                        />
                      </div>
                      <div className="col-md-4 mb-3">
                        <label className="form-label">Stock Number</label>
                        <input
                          type="text"
                          className="form-control"
                          name="stock_number"
                          value={formData.stock_number}
                          onChange={handleChange}
                        />
                      </div>
                      <div className="col-md-4 mb-3">
                        <label className="form-label">License Plate</label>
                        <input
                          type="text"
                          className="form-control"
                          name="license_plate"
                          value={formData.license_plate}
                          onChange={handleChange}
                        />
                      </div>
                    </div>
                  </div>
                </div>

                {/* Condition & Status */}
                <div className="card mb-4">
                  <div className="card-header">
                    <h5 className="mb-0">Condition & Status</h5>
                  </div>
                  <div className="card-body">
                    <div className="row">
                      <div className="col-md-4 mb-3">
                        <label className="form-label">Condition *</label>
                        <select
                          className="form-select"
                          name="condition"
                          value={formData.condition}
                          onChange={handleChange}
                          required
                        >
                          <option value="new">New</option>
                          <option value="used">Used</option>
                          <option value="certified_pre_owned">Certified Pre-Owned</option>
                        </select>
                      </div>
                      <div className="col-md-4 mb-3">
                        <label className="form-label">Status *</label>
                        <select
                          className="form-select"
                          name="status"
                          value={formData.status}
                          onChange={handleChange}
                          required
                          disabled={!isEditMode}
                        >
                          <option value="available">Available</option>
                          <option value="rented">Rented</option>
                          <option value="maintenance">Maintenance</option>
                          <option value="inactive">Inactive</option>
                        </select>
                      </div>
                      <div className="col-md-4 mb-3">
                        <label className="form-label">Mileage</label>
                        <input
                          type="number"
                          className="form-control"
                          name="mileage"
                          value={formData.mileage}
                          onChange={handleChange}
                          min="0"
                        />
                      </div>
                      <div className="col-md-12 mb-3">
                        <div className="form-check">
                          <input
                            className="form-check-input"
                            type="checkbox"
                            name="is_eligible_for_service"
                            checked={formData.is_eligible_for_service}
                            onChange={handleChange}
                          />
                          <label className="form-check-label">
                            Eligible for Service
                          </label>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>

                {/* Specifications */}
                <div className="card mb-4">
                  <div className="card-header">
                    <h5 className="mb-0">Specifications</h5>
                  </div>
                  <div className="card-body">
                    <div className="row">
                      <div className="col-md-6 mb-3">
                        <label className="form-label">Body Style</label>
                        <input
                          type="text"
                          className="form-control"
                          name="body_style"
                          value={formData.body_style}
                          onChange={handleChange}
                          placeholder="e.g., Sedan, SUV, Truck"
                        />
                      </div>
                      <div className="col-md-6 mb-3">
                        <label className="form-label">Transmission</label>
                        <input
                          type="text"
                          className="form-control"
                          name="transmission"
                          value={formData.transmission}
                          onChange={handleChange}
                          placeholder="e.g., Automatic, Manual"
                        />
                      </div>
                      <div className="col-md-6 mb-3">
                        <label className="form-label">Drivetrain</label>
                        <input
                          type="text"
                          className="form-control"
                          name="drivetrain"
                          value={formData.drivetrain}
                          onChange={handleChange}
                          placeholder="e.g., FWD, AWD, RWD"
                        />
                      </div>
                      <div className="col-md-6 mb-3">
                        <label className="form-label">Fuel Type</label>
                        <input
                          type="text"
                          className="form-control"
                          name="fuel_type"
                          value={formData.fuel_type}
                          onChange={handleChange}
                          placeholder="e.g., Gasoline, Diesel, Electric"
                        />
                      </div>
                      <div className="col-md-12 mb-3">
                        <label className="form-label">Engine</label>
                        <input
                          type="text"
                          className="form-control"
                          name="engine"
                          value={formData.engine}
                          onChange={handleChange}
                          placeholder="e.g., 2.5L 4-Cylinder"
                        />
                      </div>
                      <div className="col-md-6 mb-3">
                        <label className="form-label">MPG City</label>
                        <input
                          type="number"
                          className="form-control"
                          name="mpg_city"
                          value={formData.mpg_city}
                          onChange={handleChange}
                          min="0"
                        />
                      </div>
                      <div className="col-md-6 mb-3">
                        <label className="form-label">MPG Highway</label>
                        <input
                          type="number"
                          className="form-control"
                          name="mpg_highway"
                          value={formData.mpg_highway}
                          onChange={handleChange}
                          min="0"
                        />
                      </div>
                      <div className="col-md-6 mb-3">
                        <label className="form-label">Exterior Color</label>
                        <input
                          type="text"
                          className="form-control"
                          name="color_exterior"
                          value={formData.color_exterior}
                          onChange={handleChange}
                        />
                      </div>
                      <div className="col-md-6 mb-3">
                        <label className="form-label">Interior Color</label>
                        <input
                          type="text"
                          className="form-control"
                          name="color_interior"
                          value={formData.color_interior}
                          onChange={handleChange}
                        />
                      </div>
                      <div className="col-md-6 mb-3">
                        <label className="form-label">Seats</label>
                        <input
                          type="number"
                          className="form-control"
                          name="seats"
                          value={formData.seats}
                          onChange={handleChange}
                          min="1"
                          max="20"
                        />
                      </div>
                      <div className="col-md-6 mb-3">
                        <label className="form-label">Doors</label>
                        <input
                          type="number"
                          className="form-control"
                          name="doors"
                          value={formData.doors}
                          onChange={handleChange}
                          min="2"
                          max="6"
                        />
                      </div>
                    </div>
                  </div>
                </div>

                {/* Pricing */}
                <div className="card mb-4">
                  <div className="card-header">
                    <h5 className="mb-0">Rental Pricing</h5>
                  </div>
                  <div className="card-body">
                    <div className="row">
                      <div className="col-md-4 mb-3">
                        <label className="form-label">Daily Rate ($)</label>
                        <input
                          type="number"
                          className="form-control"
                          name="daily_rate"
                          value={formData.daily_rate}
                          onChange={handleChange}
                          min="0"
                          step="0.01"
                        />
                      </div>
                      <div className="col-md-4 mb-3">
                        <label className="form-label">Weekly Rate ($)</label>
                        <input
                          type="number"
                          className="form-control"
                          name="weekly_rate"
                          value={formData.weekly_rate}
                          onChange={handleChange}
                          min="0"
                          step="0.01"
                        />
                      </div>
                      <div className="col-md-4 mb-3">
                        <label className="form-label">Monthly Rate ($)</label>
                        <input
                          type="number"
                          className="form-control"
                          name="monthly_rate"
                          value={formData.monthly_rate}
                          onChange={handleChange}
                          min="0"
                          step="0.01"
                        />
                      </div>
                    </div>
                  </div>
                </div>

                {/* Description & Features */}
                <div className="card mb-4">
                  <div className="card-header">
                    <h5 className="mb-0">Description & Features</h5>
                  </div>
                  <div className="card-body">
                    <div className="mb-3">
                      <label className="form-label">Description</label>
                      <textarea
                        className="form-control"
                        name="description"
                        value={formData.description}
                        onChange={handleChange}
                        rows="4"
                        placeholder="Detailed description of the vehicle..."
                      ></textarea>
                    </div>
                    <div className="mb-3">
                      <label className="form-label">Features (one per line)</label>
                      <textarea
                        className="form-control"
                        name="features"
                        value={formData.features}
                        onChange={handleChange}
                        rows="6"
                        placeholder="Leather Seats&#10;Backup Camera&#10;Sunroof&#10;Bluetooth&#10;Navigation System"
                      ></textarea>
                      <small className="text-muted">Enter each feature on a new line</small>
                    </div>
                    <div className="mb-3">
                      <label className="form-label">Images (URLs, one per line)</label>
                      <textarea
                        className="form-control"
                        name="images"
                        value={formData.images}
                        onChange={handleChange}
                        rows="4"
                        placeholder="https://example.com/image1.jpg&#10;https://example.com/image2.jpg"
                      ></textarea>
                      <small className="text-muted">Enter each image URL on a new line</small>
                    </div>
                  </div>
                </div>

                {/* Action Buttons */}
                <div className="d-flex justify-content-between">
                  <button
                    type="button"
                    className="btn btn-outline-secondary"
                    onClick={() => navigate('/vehicles')}
                  >
                    Cancel
                  </button>
                  <button type="submit" className="btn btn-primary">
                    {isEditMode ? 'Update Vehicle' : 'Create Vehicle'}
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </section>
    </main>
  );
}

export default VehicleForm;
