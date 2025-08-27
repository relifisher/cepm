import React, { useState, useEffect } from 'react';
import { Table, Card, message, Tag, Space, Button, DatePicker } from 'antd';
import { EyeOutlined } from '@ant-design/icons';
import { getAllReviewsByPeriod } from '../services/api';
import { useOutletContext, useNavigate } from 'react-router-dom';
import dayjs from 'dayjs';

const HRViewPage = () => {
  const [reviews, setReviews] = useState([]);
  const [loading, setLoading] = useState(false);
  const [selectedMonth, setSelectedMonth] = useState(dayjs().format('YYYY-MM')); // Default to current month
  const { currentUserId } = useOutletContext();
  const navigate = useNavigate();

  useEffect(() => {
    if (selectedMonth) {
      fetchReviews(selectedMonth);
    }
  }, [selectedMonth]); // Fetch reviews when selectedMonth changes

  const fetchReviews = async (month) => {
    setLoading(true);
    try {
      const response = await getAllReviewsByPeriod(month);
      setReviews(response.data);
    } catch (error) {
      message.error(error.response?.data?.error || '获取所有绩效评估失败');
    }
    setLoading(false);
  };

  const handleMonthChange = (date, dateString) => {
    setSelectedMonth(dateString);
  };

  const handleViewDetails = (reviewId) => {
    navigate(`/reviews/${reviewId}/score`); // Navigate to the scoring page for viewing
  };

  const columns = [
    {
      title: '员工姓名',
      dataIndex: ['User', 'Name'],
      key: 'userName',
      width: '15%',
    },
    {
      title: '绩效周期',
      dataIndex: 'Period',
      key: 'period',
      width: '15%',
      render: (text) => dayjs(text).format('YYYY年MM月'),
    },
    {
      title: '状态',
      dataIndex: 'Status',
      key: 'status',
      width: '15%',
      render: (status) => {
        let color;
        switch (status) {
          case '草稿': color = 'default'; break;
          case '待审批': color = 'processing'; break;
          case '已批准': color = 'success'; break;
          case '已完成': color = 'success'; break;
          case '待人事确认': color = 'warning'; break;
          case '已驳回': color = 'error'; break;
          default: color = 'default';
        }
        return <Tag color={color}>{status}</Tag>;
      },
    },
    {
      title: '总分',
      dataIndex: 'TotalScore',
      key: 'totalScore',
      width: '10%',
      render: (score) => score ? score.toFixed(2) : 'N/A',
    },
    {
      title: '绩点',
      dataIndex: 'GradePoint',
      key: 'gradePoint',
      width: '10%',
      render: (gp) => gp ? gp.toFixed(2) : 'N/A',
    },
    {
      title: '操作',
      key: 'action',
      width: '15%',
      render: (_, record) => (
        <Space size="middle">
          <Button 
            icon={<EyeOutlined />} 
            onClick={() => handleViewDetails(record.ID)}
            disabled={record.Status === '草稿'}
          >
            查看详情
          </Button>
        </Space>
      ),
    },
  ];

  return (
    <div style={{ padding: '24px', background: '#f0f2f5' }}>
      <Card title="所有绩效评估" bordered={false} style={{ width: '100%' }}>
        <Space style={{ marginBottom: 16 }}>
          <DatePicker
            picker="month"
            defaultValue={dayjs(selectedMonth)}
            onChange={handleMonthChange}
            allowClear={false}
          />
        </Space>
        <Table 
          columns={columns} 
          dataSource={reviews} 
          loading={loading} 
          rowKey="ID" 
          pagination={{ pageSize: 10 }} 
        />
      </Card>
    </div>
  );
};

export default HRViewPage;
