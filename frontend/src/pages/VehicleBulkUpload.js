import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../services/api';

function VehicleBulkUpload() {
  const navigate = useNavigate();
  const [organizations, setOrganizations] = useState([]);
  const [locations, setLocations] = useState([]);
  const [filteredLocations, setFilteredLocations] = useState([]);
  const [selectedOrganization, setSelectedOrganization] = useState('');
  const [selectedLocation, setSelectedLocation] = useState('');
  const [file, setFile] = useState(null);
  const [uploading, setUploading] = useState(false);
  const [result, setResult] = useState(null);
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    fetchData();
  }, []);

  useEffect(() => {
    if (selectedOrganization) {
      const filtered = locations.filter(
        (loc) => loc.organization_id === selectedOrganization
      );
      setFilteredLocations(filtered);
      setSelectedLocation('');
    } else {
      setFilteredLocations([]);
      setSelectedLocation('');
    }
  }, [selectedOrganization, locations]);

  const fetchData = async () => {
    try {
      const [orgsResponse, locsResponse] = await Promise.all([
        api.get('/api/organizations'),
        api.get('/api/locations'),
      ]);
      setOrganizations(orgsResponse.data || []);
      setLocations(locsResponse.data || []);
    } catch (err) {
      setError('Failed to fetch organizations and locations');
    } finally {
      setLoading(false);
    }
  };

  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0];
    if (selectedFile && selectedFile.type !== 'text/csv') {
      setError('Please select a CSV file');
      setFile(null);
      return;
    }
    setFile(selectedFile);
    setError('');
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setResult(null);

    if (!file || !selectedOrganization || !selectedLocation) {
      setError('Please select organization, location, and CSV file');
      return;
    }

    setUploading(true);

    try {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('organization_id', selectedOrganization);
      formData.append('location_id', selectedLocation);

      const response = await api.post('/api/vehicles/bulk-upload', formData, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });

      setResult(response.data);
      setFile(null);

      // Reset file input
      const fileInput = document.getElementById('csvFile');
      if (fileInput) {
        fileInput.value = '';
      }

      // If all successful, redirect after a delay
      if (response.data.failed === 0) {
        setTimeout(() => navigate('/vehicles'), 2000);
      }
    } catch (err) {
      setError(err.response?.data || 'Failed to upload vehicles');
    } finally {
      setUploading(false);
    }
  };

  const downloadTemplate = () => {
    const headers = [
      'vin',
      'make',
      'model',
      'year',
      'trim',
      'color_exterior',
      'color_interior',
      'condition',
      'mileage',
      'license_plate',
      'body_style',
      'transmission',
      'drivetrain',
      'fuel_type',
      'engine',
      'mpg_city',
      'mpg_highway',
      'seats',
      'doors',
      'stock_number',
      'description',
      'daily_rate',
      'weekly_rate',
      'monthly_rate',
      'features',
    ];

    const exampleRow = [
      '1HGBH41JXMN109186',
      'Honda',
      'Accord',
      '2022',
      'EX-L',
      'Silver',
      'Black',
      'used',
      '15000',
      'ABC123',
      'Sedan',
      'Automatic',
      'FWD',
      'Gasoline',
      '2.0L I4',
      '30',
      '38',
      '5',
      '4',
      'STK001',
      'Well-maintained Honda Accord',
      '45.00',
      '280.00',
      '1000.00',
      'Bluetooth|Backup Camera|Heated Seats',
    ];

    const csvContent = [headers.join(','), exampleRow.join(',')].join('\n');
    const blob = new Blob([csvContent], { type: 'text/csv' });
    const url = window.URL.createObjectURL(blob);
    const a = document.createElement('a');
    a.href = url;
    a.download = 'vehicle_upload_template.csv';
    a.click();
    window.URL.revokeObjectURL(url);
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
            <div className="col-lg-8 mx-auto">
              <div className="d-flex justify-content-between align-items-center mb-4">
                <h1 className="fw-bolder">Bulk Upload Vehicles</h1>
                <button
                  className="btn btn-outline-secondary"
                  onClick={() => navigate('/vehicles')}
                >
                  Back to Vehicles
                </button>
              </div>

              {error && (
                <div className="alert alert-danger" role="alert">
                  {error}
                </div>
              )}

              {result && (
                <div className={`alert ${result.failed === 0 ? 'alert-success' : 'alert-warning'}`} role="alert">
                  <h5 className="alert-heading">Upload Complete</h5>
                  <p>
                    <strong>Total:</strong> {result.total} |{' '}
                    <strong>Success:</strong> {result.success} |{' '}
                    <strong>Failed:</strong> {result.failed}
                  </p>
                  {result.errors && result.errors.length > 0 && (
                    <div className="mt-3">
                      <strong>Errors:</strong>
                      <ul className="mb-0 mt-2">
                        {result.errors.map((err, idx) => (
                          <li key={idx}>{err}</li>
                        ))}
                      </ul>
                    </div>
                  )}
                  {result.failed === 0 && (
                    <p className="mb-0 mt-2">Redirecting to vehicles page...</p>
                  )}
                </div>
              )}

              <div className="card mb-4">
                <div className="card-body">
                  <h5 className="card-title mb-3">CSV Format</h5>
                  <p className="text-muted">
                    Download the template CSV file to see the required format. Required fields are:
                    <code>vin</code>, <code>make</code>, <code>model</code>, <code>year</code>.
                  </p>
                  <p className="text-muted mb-3">
                    For the <code>condition</code> field, use: <code>new</code>, <code>used</code>, or <code>certified_pre_owned</code>.
                    For <code>features</code>, separate multiple features with a pipe character (|).
                  </p>
                  <button className="btn btn-outline-primary" onClick={downloadTemplate}>
                    <i className="bi bi-download me-2"></i>
                    Download Template CSV
                  </button>
                </div>
              </div>

              <div className="card">
                <div className="card-body">
                  <h5 className="card-title mb-4">Upload Vehicles</h5>
                  <form onSubmit={handleSubmit}>
                    <div className="mb-3">
                      <label className="form-label">Organization *</label>
                      <select
                        className="form-select"
                        value={selectedOrganization}
                        onChange={(e) => setSelectedOrganization(e.target.value)}
                        required
                      >
                        <option value="">Select an organization...</option>
                        {organizations.map((org) => (
                          <option key={org.id} value={org.id}>
                            {org.name}
                          </option>
                        ))}
                      </select>
                    </div>

                    <div className="mb-3">
                      <label className="form-label">Location *</label>
                      <select
                        className="form-select"
                        value={selectedLocation}
                        onChange={(e) => setSelectedLocation(e.target.value)}
                        required
                        disabled={!selectedOrganization}
                      >
                        <option value="">
                          {selectedOrganization
                            ? 'Select a location...'
                            : 'Select organization first...'}
                        </option>
                        {filteredLocations.map((loc) => (
                          <option key={loc.id} value={loc.id}>
                            {loc.name}
                          </option>
                        ))}
                      </select>
                    </div>

                    <div className="mb-4">
                      <label htmlFor="csvFile" className="form-label">
                        CSV File *
                      </label>
                      <input
                        type="file"
                        className="form-control"
                        id="csvFile"
                        accept=".csv"
                        onChange={handleFileChange}
                        required
                      />
                      <div className="form-text">
                        Upload a CSV file containing vehicle data
                      </div>
                    </div>

                    <button
                      type="submit"
                      className="btn btn-primary"
                      disabled={uploading || !file || !selectedOrganization || !selectedLocation}
                    >
                      {uploading ? (
                        <>
                          <span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                          Uploading...
                        </>
                      ) : (
                        <>
                          <i className="bi bi-upload me-2"></i>
                          Upload Vehicles
                        </>
                      )}
                    </button>
                  </form>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </main>
  );
}

export default VehicleBulkUpload;
