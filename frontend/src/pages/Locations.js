import React, { useState, useEffect } from 'react';
import api from '../services/api';

function Locations() {
  const [locations, setLocations] = useState([]);
  const [organizations, setOrganizations] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState({
    organization_id: '',
    name: '',
    address_line1: '',
    address_line2: '',
    city: '',
    state: '',
    zip_code: '',
    country: 'USA',
    phone: '',
    email: '',
  });
  const [error, setError] = useState('');

  useEffect(() => {
    fetchData();
  }, []);

  const fetchData = async () => {
    try {
      const [locsRes, orgsRes] = await Promise.all([
        api.get('/api/locations'),
        api.get('/api/organizations'),
      ]);
      setLocations(locsRes.data || []);
      setOrganizations(orgsRes.data || []);
    } catch (err) {
      setError('Failed to fetch data');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    try {
      await api.post('/api/locations', formData);
      setFormData({
        organization_id: '',
        name: '',
        address_line1: '',
        address_line2: '',
        city: '',
        state: '',
        zip_code: '',
        country: 'USA',
        phone: '',
        email: '',
      });
      setShowForm(false);
      fetchData();
    } catch (err) {
      setError(err.response?.data || 'Failed to create location');
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Are you sure you want to delete this location?')) {
      return;
    }

    try {
      await api.delete(`/api/locations/${id}`);
      fetchData();
    } catch (err) {
      setError('Failed to delete location');
    }
  };

  const getOrgName = (orgId) => {
    const org = organizations.find((o) => o.id === orgId);
    return org ? org.name : 'Unknown';
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
                <h1 className="fw-bolder">Locations</h1>
                <button
                  className="btn btn-primary"
                  onClick={() => setShowForm(!showForm)}
                  disabled={organizations.length === 0}
                >
                  {showForm ? 'Cancel' : '+ New Location'}
                </button>
              </div>

              {organizations.length === 0 && (
                <div className="alert alert-warning" role="alert">
                  You need to create an organization first before adding locations.
                </div>
              )}

              {error && (
                <div className="alert alert-danger" role="alert">
                  {error}
                </div>
              )}

              {showForm && (
                <div className="card mb-4">
                  <div className="card-body">
                    <h5 className="card-title">Create New Location</h5>
                    <form onSubmit={handleSubmit}>
                      <div className="row">
                        <div className="col-md-6 mb-3">
                          <label className="form-label">Organization *</label>
                          <select
                            className="form-select"
                            value={formData.organization_id}
                            onChange={(e) => setFormData({ ...formData, organization_id: e.target.value })}
                            required
                          >
                            <option value="">Select an organization</option>
                            {organizations.map((org) => (
                              <option key={org.id} value={org.id}>
                                {org.name}
                              </option>
                            ))}
                          </select>
                        </div>
                        <div className="col-md-6 mb-3">
                          <label className="form-label">Location Name *</label>
                          <input
                            type="text"
                            className="form-control"
                            value={formData.name}
                            onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                            placeholder="e.g., Downtown Showroom"
                            required
                          />
                        </div>
                      </div>
                      <div className="mb-3">
                        <label className="form-label">Address Line 1</label>
                        <input
                          type="text"
                          className="form-control"
                          value={formData.address_line1}
                          onChange={(e) => setFormData({ ...formData, address_line1: e.target.value })}
                        />
                      </div>
                      <div className="mb-3">
                        <label className="form-label">Address Line 2</label>
                        <input
                          type="text"
                          className="form-control"
                          value={formData.address_line2}
                          onChange={(e) => setFormData({ ...formData, address_line2: e.target.value })}
                        />
                      </div>
                      <div className="row">
                        <div className="col-md-4 mb-3">
                          <label className="form-label">City</label>
                          <input
                            type="text"
                            className="form-control"
                            value={formData.city}
                            onChange={(e) => setFormData({ ...formData, city: e.target.value })}
                          />
                        </div>
                        <div className="col-md-4 mb-3">
                          <label className="form-label">State</label>
                          <input
                            type="text"
                            className="form-control"
                            value={formData.state}
                            onChange={(e) => setFormData({ ...formData, state: e.target.value })}
                            placeholder="CA"
                          />
                        </div>
                        <div className="col-md-4 mb-3">
                          <label className="form-label">Zip Code</label>
                          <input
                            type="text"
                            className="form-control"
                            value={formData.zip_code}
                            onChange={(e) => setFormData({ ...formData, zip_code: e.target.value })}
                          />
                        </div>
                      </div>
                      <div className="row">
                        <div className="col-md-6 mb-3">
                          <label className="form-label">Phone</label>
                          <input
                            type="tel"
                            className="form-control"
                            value={formData.phone}
                            onChange={(e) => setFormData({ ...formData, phone: e.target.value })}
                          />
                        </div>
                        <div className="col-md-6 mb-3">
                          <label className="form-label">Email</label>
                          <input
                            type="email"
                            className="form-control"
                            value={formData.email}
                            onChange={(e) => setFormData({ ...formData, email: e.target.value })}
                          />
                        </div>
                      </div>
                      <button type="submit" className="btn btn-primary">
                        Create Location
                      </button>
                    </form>
                  </div>
                </div>
              )}

              <div className="card">
                <div className="card-body">
                  {locations.length === 0 ? (
                    <div className="text-center py-5">
                      <p className="text-muted">No locations yet. Create one to get started!</p>
                    </div>
                  ) : (
                    <div className="table-responsive">
                      <table className="table table-hover">
                        <thead>
                          <tr>
                            <th>Name</th>
                            <th>Organization</th>
                            <th>Address</th>
                            <th>Phone</th>
                            <th>Status</th>
                            <th>Actions</th>
                          </tr>
                        </thead>
                        <tbody>
                          {locations.map((loc) => (
                            <tr key={loc.id}>
                              <td>{loc.name}</td>
                              <td>{getOrgName(loc.organization_id)}</td>
                              <td>
                                {loc.city && loc.state ? `${loc.city}, ${loc.state}` : 'N/A'}
                              </td>
                              <td>{loc.phone || 'N/A'}</td>
                              <td>
                                {loc.is_active ? (
                                  <span className="badge bg-success">Active</span>
                                ) : (
                                  <span className="badge bg-secondary">Inactive</span>
                                )}
                              </td>
                              <td>
                                <button
                                  className="btn btn-sm btn-danger"
                                  onClick={() => handleDelete(loc.id)}
                                >
                                  Delete
                                </button>
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

export default Locations;
