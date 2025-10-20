export default {
  server: {
    proxy: {
      '/leaderboard': 'http://localhost:8080',
      '/ws': {
        target: 'ws://localhost:8080',
        ws: true
      }
    }
  }
};
