import React from 'react';
import { Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

function Home() {
  const { isAuthenticated } = useAuth();

  return (
    <main className="flex-shrink-0">
      {/* Header */}
      <header className="bg-dark py-5">
        <div className="container px-5">
          <div className="row gx-5 align-items-center justify-content-center">
            <div className="col-lg-8 col-xl-7 col-xxl-6">
              <div className="my-5 text-center text-xl-start">
                <h1 className="display-5 fw-bolder text-white mb-2">
                  FleetPass Marketplace
                </h1>
                <p className="lead fw-normal text-white-50 mb-4">
                  The modern solution for automotive dealerships to maximize fleet utilization
                  and build stronger customer relationships.
                </p>
                <div className="d-grid gap-3 d-sm-flex justify-content-sm-center justify-content-xl-start">
                  {isAuthenticated ? (
                    <Link className="btn btn-primary btn-lg px-4 me-sm-3" to="/dashboard">
                      Go to Dashboard
                    </Link>
                  ) : (
                    <>
                      <Link className="btn btn-primary btn-lg px-4 me-sm-3" to="/login">
                        Get Started
                      </Link>
                      <Link className="btn btn-outline-light btn-lg px-4" to="/login">
                        Login
                      </Link>
                    </>
                  )}
                </div>
              </div>
            </div>
            <div className="col-xl-5 col-xxl-6 d-none d-xl-block text-center">
              <img
                className="img-fluid rounded-3 my-5"
                src="https://images.unsplash.com/photo-1492144534655-ae79c964c9d7?w=600&h=400&fit=crop"
                alt="Automotive"
              />
            </div>
          </div>
        </div>
      </header>

      {/* Features Section */}
      <section className="py-5" id="features">
        <div className="container px-5 my-5">
          <div className="row gx-5">
            <div className="col-lg-4 mb-5 mb-lg-0">
              <h2 className="fw-bolder mb-0">White-Label Ready</h2>
            </div>
            <div className="col-lg-8">
              <div className="row gx-5 row-cols-1 row-cols-md-2">
                <div className="col mb-5 h-100">
                  <div className="feature bg-primary bg-gradient text-white rounded-3 mb-3">
                    <i className="bi bi-building"></i>
                  </div>
                  <h2 className="h5">Multi-Organization Support</h2>
                  <p className="mb-0">
                    Manage multiple dealership organizations with complete data isolation and privacy.
                  </p>
                </div>
                <div className="col mb-5 h-100">
                  <div className="feature bg-primary bg-gradient text-white rounded-3 mb-3">
                    <i className="bi bi-car-front"></i>
                  </div>
                  <h2 className="h5">Fleet Management</h2>
                  <p className="mb-0">
                    Track and manage your entire vehicle inventory across multiple locations.
                  </p>
                </div>
                <div className="col mb-5 mb-md-0 h-100">
                  <div className="feature bg-primary bg-gradient text-white rounded-3 mb-3">
                    <i className="bi bi-people"></i>
                  </div>
                  <h2 className="h5">Customer Relationships</h2>
                  <p className="mb-0">
                    Build stronger, long-term relationships with your customers through personalized service.
                  </p>
                </div>
                <div className="col h-100">
                  <div className="feature bg-primary bg-gradient text-white rounded-3 mb-3">
                    <i className="bi bi-shield-check"></i>
                  </div>
                  <h2 className="h5">Role-Based Access</h2>
                  <p className="mb-0">
                    Secure RBAC system for org owners, admins, and users with granular permissions.
                  </p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>
    </main>
  );
}

export default Home;
