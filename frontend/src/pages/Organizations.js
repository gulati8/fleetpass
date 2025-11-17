import React, { useState, useEffect } from 'react';
import api from '../services/api';

function Organizations() {
  const [organizations, setOrganizations] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [formData, setFormData] = useState({ name: '', slug: '' });
  const [error, setError] = useState('');

  useEffect(() => {
    fetchOrganizations();
  }, []);

  const fetchOrganizations = async () => {
    try {
      const response = await api.get('/api/organizations');
      setOrganizations(response.data || []);
    } catch (err) {
      setError('Failed to fetch organizations');
    } finally {
      setLoading(false);
    }
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    try {
      await api.post('/api/organizations', formData);
      setFormData({ name: '', slug: '' });
      setShowForm(false);
      fetchOrganizations();
    } catch (err) {
      setError(err.response?.data || 'Failed to create organization');
    }
  };

  const handleDelete = async (id) => {
    if (!window.confirm('Are you sure you want to delete this organization?')) {
      return;
    }

    try {
      await api.delete(`/api/organizations/${id}`);
      fetchOrganizations();
    } catch (err) {
      setError('Failed to delete organization');
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
          <div className="row gx-5">
            <div className="col-lg-12">
              <div className="d-flex justify-content-between align-items-center mb-4">
                <h1 className="fw-bolder">Organizations</h1>
                <button
                  className="btn btn-primary"
                  onClick={() => setShowForm(!showForm)}
                >
                  {showForm ? 'Cancel' : '+ New Organization'}
                </button>
              </div>

              {error && (
                <div className="alert alert-danger" role="alert">
                  {error}
                </div>
              )}

              {showForm && (
                <div className="card mb-4">
                  <div className="card-body">
                    <h5 className="card-title">Create New Organization</h5>
                    <form onSubmit={handleSubmit}>
                      <div className="mb-3">
                        <label className="form-label">Name</label>
                        <input
                          type="text"
                          className="form-control"
                          value={formData.name}
                          onChange={(e) => setFormData({ ...formData, name: e.target.value })}
                          required
                        />
                      </div>
                      <div className="mb-3">
                        <label className="form-label">Slug (URL-friendly identifier)</label>
                        <input
                          type="text"
                          className="form-control"
                          value={formData.slug}
                          onChange={(e) => setFormData({ ...formData, slug: e.target.value })}
                          placeholder="e.g., acme-motors"
                          required
                        />
                      </div>
                      <button type="submit" className="btn btn-primary">
                        Create Organization
                      </button>
                    </form>
                  </div>
                </div>
              )}

              <div className="card">
                <div className="card-body">
                  {organizations.length === 0 ? (
                    <div className="text-center py-5">
                      <p className="text-muted">No organizations yet. Create one to get started!</p>
                    </div>
                  ) : (
                    <div className="table-responsive">
                      <table className="table table-hover">
                        <thead>
                          <tr>
                            <th>Name</th>
                            <th>Slug</th>
                            <th>Status</th>
                            <th>Created</th>
                            <th>Actions</th>
                          </tr>
                        </thead>
                        <tbody>
                          {organizations.map((org) => (
                            <tr key={org.id}>
                              <td>{org.name}</td>
                              <td><code>{org.slug}</code></td>
                              <td>
                                {org.is_active ? (
                                  <span className="badge bg-success">Active</span>
                                ) : (
                                  <span className="badge bg-secondary">Inactive</span>
                                )}
                              </td>
                              <td>{new Date(org.created_at).toLocaleDateString()}</td>
                              <td>
                                <button
                                  className="btn btn-sm btn-danger"
                                  onClick={() => handleDelete(org.id)}
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

export default Organizations;
