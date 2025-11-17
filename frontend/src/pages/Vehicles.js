import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../services/api';

function Vehicles() {
  const [vehicles, setVehicles] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const navigate = useNavigate();

  useEffect(() => {
    fetchVehicles();
  }, []);

  const fetchVehicles = async () => {
    try {
      const response = await api.get('/api/vehicles');
      setVehicles(response.data || []);
    } catch (err) {
      setError('Failed to fetch vehicles');
    } finally {
      setLoading(false);
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Are you sure you want to delete this vehicle?')) {
      return;
    }

    try {
      await api.delete(`/api/vehicles/${id}`);
      fetchVehicles();
    } catch (err) {
      setError('Failed to delete vehicle');
    }
  };

  const getStatusBadge = (status) => {
    const badges = {
      available: 'bg-success',
      rented: 'bg-warning',
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

  return (
    <main className="flex-shrink-0">
      <section className="py-5">
        <div className="container px-5">
          <div className="row gx-5">
            <div className="col-lg-12">
              <div className="d-flex justify-content-between align-items-center mb-4">
                <h1 className="fw-bolder">Vehicles</h1>
                <button
                  className="btn btn-primary"
                  onClick={() => navigate('/vehicles/new')}
                >
                  + New Vehicle
                </button>
              </div>

              {error && (
                <div className="alert alert-danger" role="alert">
                  {error}
                </div>
              )}

              <div className="card">
                <div className="card-body">
                  {vehicles.length === 0 ? (
                    <div className="text-center py-5">
                      <p className="text-muted">No vehicles yet. Add one to get started!</p>
                    </div>
                  ) : (
                    <div className="table-responsive">
                      <table className="table table-hover">
                        <thead>
                          <tr>
                            <th>Vehicle</th>
                            <th>VIN</th>
                            <th>Year</th>
                            <th>Mileage</th>
                            <th>Daily Rate</th>
                            <th>Status</th>
                            <th>Actions</th>
                          </tr>
                        </thead>
                        <tbody>
                          {vehicles.map((vehicle) => (
                            <tr key={vehicle.id}>
                              <td>
                                <strong>{vehicle.make} {vehicle.model}</strong>
                                {vehicle.trim && <div className="text-muted small">{vehicle.trim}</div>}
                              </td>
                              <td><code className="small">{vehicle.vin}</code></td>
                              <td>{vehicle.year}</td>
                              <td>{vehicle.mileage?.toLocaleString()} mi</td>
                              <td>${vehicle.daily_rate?.toFixed(2)}/day</td>
                              <td>
                                <span className={`badge ${getStatusBadge(vehicle.status)}`}>
                                  {vehicle.status}
                                </span>
                              </td>
                              <td>
                                <div className="btn-group" role="group">
                                  <button
                                    className="btn btn-sm btn-outline-primary"
                                    onClick={() => navigate(`/vehicles/${vehicle.id}`)}
                                  >
                                    View
                                  </button>
                                  <button
                                    className="btn btn-sm btn-outline-secondary"
                                    onClick={() => navigate(`/vehicles/${vehicle.id}/edit`)}
                                  >
                                    Edit
                                  </button>
                                  <button
                                    className="btn btn-sm btn-outline-danger"
                                    onClick={() => handleDelete(vehicle.id)}
                                  >
                                    Delete
                                  </button>
                                </div>
                              </td>
                            </tr>
                          ))}
                        </tbody>
                      </table>
                    </div>
                  )}
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </main>
  );
}

export default Vehicles;
