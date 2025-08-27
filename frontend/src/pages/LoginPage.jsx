import React, { useEffect, useState } from 'react';
import { Button, Card, Spin, Typography, Result, Space, message } from 'antd';
import { useNavigate, useLocation } from 'react-router-dom';
import axios from 'axios';

const { Title, Paragraph } = Typography;

const WECHAT_WORK_CORP_ID = 'ww40bcae13177b01e9'; // Replace with your actual CorpID
const WECHAT_WORK_AGENT_ID = '1000062'; // Replace with your actual AgentID
const REDIRECT_URI = encodeURIComponent('http://localhost:3100/login'); // Your frontend login callback URL

const BACKEND_LOGIN_URL = 'http://localhost:8090/api/v1/wechat/login';

const LoginPage = () => {
  const navigate = useNavigate();
  const location = useLocation();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    const queryParams = new URLSearchParams(location.search);
    const code = queryParams.get('code');

    if (code) {
      // This is a callback from WeChat Work, try to log in
      handleWechatLoginCallback(code);
    } else if (localStorage.getItem('jwtToken')) {
      // If already logged in, redirect to app
      navigate('/app');
    }
  }, [location, navigate]);

  const handleWechatLoginCallback = async (code) => {
    setLoading(true);
    setError(null);
    try {
      const response = await axios.get(`${BACKEND_LOGIN_URL}?code=${code}`);
      const { token, user } = response.data;

      localStorage.setItem('jwtToken', token);
      localStorage.setItem('currentUser', JSON.stringify(user)); // Store user info
      message.success('登录成功！');
      navigate('/app');
    } catch (err) {
      console.error('WeChat login failed:', err);
      setError(err.response?.data?.error || '微信登录失败，请重试。');
      message.error(err.response?.data?.error || '微信登录失败，请重试。');
    } finally {
      setLoading(false);
    }
  };

  const handleMobileLogin = () => {
    // For mobile, redirect directly to WeChat Work authorization URL
    const authUrl = `https://open.weixin.qq.com/connect/oauth2/authorize?appid=${WECHAT_WORK_CORP_ID}&redirect_uri=${REDIRECT_URI}&response_type=code&scope=snsapi_base&agentid=${WECHAT_WORK_AGENT_ID}#wechat_redirect`;
    window.location.href = authUrl;
  };

  const handlePCLogin = () => {
    // For PC, open the QR code login page in a new window/tab
    // Note: This is a simplified approach. For a true embedded QR code, you'd use WeChat Work's JS-SDK for PC.
    const qrLoginUrl = `https://open.work.weixin.qq.com/wwopen/sso/qrConnect?appid=${WECHAT_WORK_CORP_ID}&agentid=${WECHAT_WORK_AGENT_ID}&redirect_uri=${REDIRECT_URI}&state=STATE#wechat_redirect`;
    window.open(qrLoginUrl, '_blank');
  };

  if (loading) {
    return (
      <div style={{ textAlign: 'center', padding: '50px' }}>
        <Spin size="large" tip="登录中..." />
      </div>
    );
  }

  if (error) {
    return (
      <Result
        status="error"
        title="登录失败"
        subTitle={error}
        extra={[
          <Button type="primary" key="retry" onClick={() => setError(null)}>
            重试
          </Button>,
        ]}
      />
    );
  }

  return (
    <div style={{ textAlign: 'center', padding: '50px', background: '#f0f2f5', minHeight: '100vh', display: 'flex', flexDirection: 'column', justifyContent: 'center', alignItems: 'center' }}>
      <Card style={{ width: 400, padding: '20px' }}>
        <Title level={2}>月度绩效管理系统</Title>
        <Paragraph>请选择登录方式：</Paragraph>
        <Space direction="vertical" size="large" style={{ width: '100%' }}>
          <Button type="primary" size="large" onClick={handleMobileLogin} style={{ width: '100%' }}>
            企业微信手机端登录
          </Button>
          <Button size="large" onClick={handlePCLogin} style={{ width: '100%' }}>
            企业微信PC端扫码登录
          </Button>
        </Space>
      </Card>
    </div>
  );
};

export default LoginPage;
