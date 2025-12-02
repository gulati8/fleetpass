describe('API Service', () => {
  beforeEach(() => {
    localStorage.clear();
    delete window.location;
    window.location = { href: '' };
  });

  test('API module loads successfully', () => {
    const api = require('./api').default;
    expect(api).toBeDefined();
  });

  test('authAPI is exported', () => {
    const { authAPI } = require('./api');
    expect(authAPI).toBeDefined();
    expect(authAPI.login).toBeDefined();
    expect(authAPI.getProfile).toBeDefined();
  });
});
