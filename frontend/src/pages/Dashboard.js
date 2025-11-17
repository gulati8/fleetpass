import React from 'react';
import { useAuth } from '../context/AuthContext';

function Dashboard() {
  const { user } = useAuth();

  return (
    <main className="flex-shrink-0">
      <section className="py-5">
        <div className="container px-5 my-5">
          <div className="row gx-5">
            <div className="col-lg-12">
              <h1 className="fw-bolder mb-4">Welcome to FleetPass Dashboard</h1>
              <div className="card mb-4">
                <div className="card-body">
                  <h5 className="card-title">Your Profile</h5>
                  <p className="card-text">
                    <strong>Email:</strong> {user?.email}<br />
                    <strong>Role:</strong> <span className="badge bg-primary">{user?.role}</span><br />
                    <strong>User ID:</strong> {user?.id}
                  </p>
                </div>
              </div>

              <div className="row gx-5">
                <div className="col-md-4 mb-4">
                  <div className="card bg-primary text-white">
                    <div className="card-body">
                      <h5 className="card-title">Organizations</h5>
                      <p className="display-4">0</p>
                      <p className="card-text">Active Organizations</p>
                    </div>
                  </div>
                </div>
                <div className="col-md-4 mb-4">
                  <div className="card bg-success text-white">
                    <div className="card-body">
                      <h5 className="card-title">Vehicles</h5>
                      <p className="display-4">0</p>
                      <p className="card-text">Total Vehicles</p>
                    </div>
                  </div>
                </div>
                <div className="col-md-4 mb-4">
                  <div className="card bg-info text-white">
                    <div className="card-body">
                      <h5 className="card-title">Customers</h5>
                      <p className="display-4">0</p>
                      <p className="card-text">Active Customers</p>
                    </div>
                  </div>
                </div>
              </div>

              <div className="alert alert-info" role="alert">
                <h5 className="alert-heading">Protected Route!</h5>
                <p>This page is only accessible to authenticated users. Your JWT token is being sent with each request to verify your identity.</p>
                <hr />
                <p className="mb-0">Next steps: Connect to the database and build out your vehicle management features!</p>
              </div>
            </div>
          </div>
        </div>
      </section>
    </main>
  );
}

export default Dashboard;
