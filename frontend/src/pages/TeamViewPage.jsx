import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { Table, Button, message, Spin, Result, Tag, Space, Modal, Input } from 'antd';
import { listTeamReviews, approvePerformanceReview, rejectPerformanceReview } from '../services/api';
import { useOutletContext } from 'react-router-dom';

const statusTags = {
  'Draft': 'default',
  'PendingApproval': 'processing',
  'Approved': 'processing',
  'Evaluating': 'warning',
  'Completed': 'success',
  'Rejected': 'error',
  '草稿': 'default',
  '待打分': 'warning',
  '待审批': 'processing',
};

const TeamViewPage = () => {
  const [reviews, setReviews] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [currentReviewId, setCurrentReviewId] = useState(null);
  const [comment, setComment] = useState('');
  const [actionType, setActionType] = useState(''); // 'approve' or 'reject'

  const navigate = useNavigate();

  const { currentUserId, isManager } = useOutletContext();

  const fetchReviews = async () => {
    if (!isManager) {
      setError('您没有权限查看团队绩效。');
      setLoading(false);
      return;
    }
    try {
      setLoading(true);
      const response = await listTeamReviews(currentUserId);
      const reviewsData = Array.isArray(response.data) ? response.data : [];

      // *** DEBUGGING STEP: Remove nested Items array to prevent tree-data logic in Table ***
      const simplifiedReviews = reviewsData.map(({ Items, ...rest }) => rest);

      setReviews(simplifiedReviews);
    } catch (err) {
      setError('无法加载团队绩效记录，请稍后再试。');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchReviews();
  }, [currentUserId, isManager]); // Depend on currentUserId and isManager to refetch when user/role changes

  const showApprovalModal = (reviewId, type) => {
    setCurrentReviewId(reviewId);
    setActionType(type);
    setIsModalVisible(true);
  };

  const handleModalOk = async () => {
    try {
      if (actionType === 'approve') {
        await approvePerformanceReview(currentReviewId, currentUserId, comment);
        message.success('绩效评估已成功批准！');
      } else if (actionType === 'reject') {
        await rejectPerformanceReview(currentReviewId, currentUserId, comment);
        message.success('绩效评估已成功驳回！');
      }
      setIsModalVisible(false);
      setComment('');
      fetchReviews(); // Refresh the list
    } catch (err) {
      const errorMsg = err.response?.data?.error || '操作失败，请重试。';
      message.error(errorMsg);
    }
  };

  const handleModalCancel = () => {
    setIsModalVisible(false);
    setComment('');
  };

  const columns = [
    {
      title: '员工姓名',
      dataIndex: ['User', 'Name'], // Access nested data
      key: 'userName',
    },
    {
      title: '绩效周期',
      dataIndex: 'Period',
      key: 'period',
      // sorter: (a, b) => a.Period.localeCompare(b.Period),
      // defaultSortOrder: 'descend',
    },
    {
      title: '状态',
      dataIndex: 'Status',
      key: 'status',
      render: (status) => <Tag color={statusTags[status] || 'default'}>{status}</Tag>,
    },
    {
      title: '总分',
      dataIndex: 'TotalScore',
      key: 'totalScore',
      render: (score) => (score ? score.toFixed(2) : '--'),
    },
    {
      title: '操作',
      key: 'action',
      render: (_, record) => (
        <Space size="middle">
          <Button type="link" onClick={() => navigate(`/reviews/${record.ID}/score`)}>查看/打分</Button>
          {isManager && record.Status === '待审批' && (
            <>
              <Button type="link" onClick={() => showApprovalModal(record.ID, 'approve')}>批准</Button>
              <Button type="link" danger onClick={() => showApprovalModal(record.ID, 'reject')}>驳回</Button>
            </>
          )}
        </Space>
      ),
    },
  ];

  if (loading) {
    return <div style={{ textAlign: 'center', padding: '50px' }}><Spin size="large" /></div>;
  }

  if (error) {
    return <Result status="error" title="加载失败" subTitle={error} />;
  }

  return (
    <div style={{ padding: '24px' }}>
      <h1>团队绩效管理</h1>
      <Table columns={columns} dataSource={reviews} rowKey="ID" />

      <Modal
        title={actionType === 'approve' ? '批准绩效评估' : '驳回绩效评估'}
        open={isModalVisible}
        onOk={handleModalOk}
        onCancel={handleModalCancel}
        okText={actionType === 'approve' ? '批准' : '驳回'}
        cancelText="取消"
      >
        <Input.TextArea
          rows={4}
          placeholder={actionType === 'approve' ? '请输入批准意见 (可选)' : '请输入驳回原因 (必填)'}
          value={comment}
          onChange={e => setComment(e.target.value)}
        />
      </Modal>
    </div>
  );
};

export default TeamViewPage;