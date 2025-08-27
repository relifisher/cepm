import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import MainLayout from './components/MainLayout';
import PerformancePlanPage from './pages/PerformancePlanPage';
import HistoryPage from './pages/HistoryPage';
import DashboardPage from './pages/DashboardPage';
import TeamViewPage from './pages/TeamViewPage';
import ScorePerformancePage from './pages/ScorePerformancePage';
import HRViewPage from './pages/HRViewPage'; // Import the new HRViewPage
import AdminPage from './pages/AdminPage'; // Import the new AdminPage
import LoginPage from './pages/LoginPage'; // Import the new LoginPage
import { Result, Button } from 'antd';

// A simple placeholder for the logout functionality
const LogoutPage = () => (
  <Result
    status="success"
    title="您已成功退出登录"
    extra={[
      <Button type="primary" key="console" href="/">
        重新登录
      </Button>,
    ]}
  />
);

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<LoginPage />} />
        <Route path="/" element={<LoginPage />} /> {/* Default route to login page */}
        <Route path="/app" element={<MainLayout />}> {/* Main application routes under /app */}
          <Route index element={<PerformancePlanPage />} />
          <Route path="history" element={<HistoryPage />} />
          <Route path="dashboard" element={<DashboardPage />} />
          <Route path="team" element={<TeamViewPage />} />
          <Route path="hr-view" element={<HRViewPage />} />
          <Route path="admin" element={<AdminPage />} />
          <Route path="reviews/:id/score" element={<ScorePerformancePage />} />
        </Route>
        <Route path="/logout" element={<LogoutPage />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;