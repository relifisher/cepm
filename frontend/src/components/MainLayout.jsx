import React, { useState, useEffect } from 'react';
import { Layout, Menu, Dropdown, Avatar, Space, Select } from 'antd';
import { Outlet, useNavigate } from 'react-router-dom';
import { UserOutlined, EditOutlined, HistoryOutlined, TeamOutlined, LineChartOutlined, LogoutOutlined } from '@ant-design/icons';

const { Header, Content, Footer } = Layout;
const { Option } = Select;

const MainLayout = () => {
  const navigate = useNavigate();
  const [currentUser, setCurrentUser] = useState(null);

  useEffect(() => {
    const storedUser = localStorage.getItem('currentUser');
    if (storedUser) {
      setCurrentUser(JSON.parse(storedUser));
    } else {
      // If no user in localStorage, redirect to login
      navigate('/login');
    }
  }, [navigate]);

  const handleMenuClick = (e) => {
    if (e.key === '/logout') {
      localStorage.removeItem('jwtToken');
      localStorage.removeItem('currentUser');
      navigate('/logout');
    } else {
      navigate(`/app${e.key}`); // Navigate relative to /app
    }
  };

  const getMenuItems = () => {
    const items = [
      {
        key: '',
        label: '填写/更新绩效计划',
        icon: <EditOutlined />,
      },
      {
        key: '/history',
        label: '查看历史绩效',
        icon: <HistoryOutlined />,
      },
      {
        key: '/dashboard',
        label: '半年度绩效统计',
        icon: <LineChartOutlined />,
      },
    ];

    if (currentUser) {
      const isManager = currentUser.Role?.Name === '组长' || currentUser.Role?.Name === '总监';
      if (isManager) {
        items.push({
          key: '/team',
          label: '团队绩效管理',
          icon: <TeamOutlined />,
        });
      }

      const isHR = currentUser.Role?.Name === '人事' || currentUser.Role?.Name === 'HR';
      if (isHR) {
        items.push({
          key: '/hr-view',
          label: '所有绩效评估 (人事)',
          icon: <TeamOutlined />,
        });
      }

      const isAdmin = currentUser.Role?.Name === '管理员';
      if (isAdmin) {
        items.push({
          key: '/admin',
          label: '管理员设置',
          icon: <UserOutlined />,
        });
      }
    }

    items.push({ type: 'divider' });
    items.push({
      key: '/logout',
      label: '退出登录',
      icon: <LogoutOutlined />,
      danger: true,
    });

    return items;
  };

  if (!currentUser) {
    return null; // Or a loading spinner
  }

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', background: '#fff', borderBottom: '1px solid #f0f0f0' }}>
        <div className="logo" style={{ cursor: 'pointer' }} onClick={() => navigate('/app')}>
          <img src="/logo.svg" alt="Logo" style={{ height: '32px' }} /> 
          <span style={{ marginLeft: '16px', fontSize: '18px', fontWeight: 'bold' }}>月度绩效管理系统</span>
        </div>
        <Space>
          {/* Removed user selection dropdown */}
          <Dropdown menu={{ items: getMenuItems(), onClick: handleMenuClick }} trigger={['click']}>
            <a onClick={e => e.preventDefault()}>
              <Space>
                <Avatar icon={<UserOutlined />} />
                <span>{currentUser.Name}</span>
              </Space>
            </a>
          </Dropdown>
        </Space>
      </Header>
      <Content>
        <Outlet context={{ currentUserId: currentUser.ID, currentUser, isManager: currentUser.Role?.Name === '组长' || currentUser.Role?.Name === '总监' }} />
      </Content>
      <Footer style={{ textAlign: 'center' }}>
        CEPM ©2025 Created by Gemini
      </Footer>
    </Layout>
  );
};

export default MainLayout;
