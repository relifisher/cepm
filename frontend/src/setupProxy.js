const { createProxyMiddleware } = require('http-proxy-middleware');

// A more specific proxy configuration
module.exports = function(app) {
  app.use(
    '/api/v1',
    createProxyMiddleware({
      target: 'http://localhost:8090',
      changeOrigin: true,
      // No need to rewrite path, as the backend router is also under /api/v1
    })
  );
};
