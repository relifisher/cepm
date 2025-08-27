import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Form, Input, Button, Table, InputNumber, message, Descriptions, Card, Spin, Result, Typography } from 'antd';
import { getPerformanceReview, scorePerformanceReview } from '../services/api';
import { useOutletContext } from 'react-router-dom';
import * as XLSX from 'xlsx'; // Import xlsx library

const { Title } = Typography;

const calculateGradePoint = (totalScore) => {
  if (totalScore >= 90 && totalScore <= 100) {
    return 1.0;
  } else if (totalScore >= 60 && totalScore < 90) {
    return 0.8;
  } else if (totalScore < 60) {
    return 0;
  } else if (totalScore > 100) {
    return totalScore / 100; // Total score / 100 * 100% as per requirement
  }
  return 0; // Default or error case
};

const ScorePerformancePage = () => {
  const [form] = Form.useForm();
  const { id } = useParams(); // Get review ID from URL
  const navigate = useNavigate();

  const [review, setReview] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [dynamicTotalScore, setDynamicTotalScore] = useState(0);
  const [dynamicGradePoint, setDynamicGradePoint] = useState(0);
  const [isReadOnly, setIsReadOnly] = useState(true); // Default to read-only for safety

  const { currentUserId, isManager } = useOutletContext(); // Get current user info

  const calculateDynamicScores = (currentItems) => {
    let total = 0;
    currentItems.forEach(item => {
      if (item.Weight && item.Score !== null && item.Score !== undefined) {
        total += (item.Weight / 100.0) * item.Score;
      }
    });
    setDynamicTotalScore(total);
    setDynamicGradePoint(calculateGradePoint(total));
  };

  useEffect(() => {
    const fetchReview = async () => {
      try {
        setLoading(true);
        const response = await getPerformanceReview(id);
        const fetchedReview = response.data;
        setReview(fetchedReview);

        // Determine if the page should be in read-only mode
        // A user cannot score their own review.
        if (currentUserId === fetchedReview.UserID) {
          setIsReadOnly(true);
        } else {
          // In a real-world scenario, you might have more complex logic
          // to check if the current user is the designated scorer (e.g., the direct manager).
          // For now, we'll assume if you're not the owner, you're the scorer.
          setIsReadOnly(false);
        }
        
        // Set initial form values from fetched data
        const initialValues = {};
        fetchedReview.Items.forEach(item => {
          initialValues[`completion_${item.ID}`] = item.CompletionDetails;
          initialValues[`score_${item.ID}`] = item.Score;
        });
        initialValues.finalComment = fetchedReview.FinalComment;
        form.setFieldsValue(initialValues);

        // Calculate initial dynamic scores
        calculateDynamicScores(fetchedReview.Items);

      } catch (err) {
        setError('无法加载绩效评估详情，请检查ID是否正确或稍后再试。');
      } finally {
        setLoading(false);
      }
    };
    fetchReview();
  }, [id, form, currentUserId]);

  const onFinish = async (values) => {
    if (isReadOnly) return; // Do not submit if in read-only mode
    setLoading(true);
    const itemsToScore = review.Items.map(item => ({
      id: item.ID,
      completionDetails: values[`completion_${item.ID}`],
      score: values[`score_${item.ID}`],
    }));

    const scoreData = {
      items: itemsToScore,
      finalComment: values.finalComment,
    };

    try {
      await scorePerformanceReview(id, scoreData);
      message.success('绩效评估打分成功!');
      navigate('/history'); // Redirect to history page after success
    } catch (err) {
      const errorMsg = err.response?.data?.error || '打分失败，请重试。';
      message.error(errorMsg);
      setLoading(false);
    }
  };

  const handleValuesChange = (changedValues, allValues) => {
    if (isReadOnly) return;
    const isScoreChanged = Object.keys(changedValues).some(key => key.startsWith('score_'));
    if (isScoreChanged) {
      const updatedItems = review.Items.map(item => {
        const score = allValues[`score_${item.ID}`];
        return {
          ...item,
          Score: score,
        };
      });
      calculateDynamicScores(updatedItems);
    }
  };

  const handlePrint = () => {
    window.print();
  };

  const handleExport = () => {
    if (!review) {
      message.error('绩效数据未加载，无法导出。');
      return;
    }

    const data = [];

    // Add basic info
    data.push(['姓名', review.User?.Name || 'N/A']);
    data.push(['部门', review.User?.Department?.Name || 'N/A']);
    data.push(['岗位', review.User?.Role?.Name || 'N/A']);
    data.push(['绩效周期', review.Period]);
    data.push(['总分', dynamicTotalScore.toFixed(2)]);
    data.push(['绩点', dynamicGradePoint.toFixed(2)]);
    data.push([]); // Empty row for separation

    // Add performance items header
    data.push(['绩效项详情']);
    data.push(['考核指标', '指标描述', '目标/衡量标准', '权重 (%)', '完成情况', '得分']);

    // Add performance items data
    review.Items.forEach(item => {
      data.push([
        item.Title,
        item.Description,
        item.Target,
        item.Weight,
        item.CompletionDetails || '',
        item.Score !== null && item.Score !== undefined ? item.Score : 'N/A',
      ]);
    });
    data.push([]); // Empty row for separation

    // Add final comment
    data.push(['最终评语']);
    data.push([review.FinalComment || '']);

    const ws = XLSX.utils.aoa_to_sheet(data);
    const wb = XLSX.utils.book_new();
    XLSX.utils.book_append_sheet(wb, ws, "绩效评估");
    XLSX.writeFile(wb, `绩效评估_${review.User?.Name}_${review.Period}.xlsx`);

    message.success('绩效评估已导出为Excel文件。');
  };

  const columns = [
    { title: '考核指标', dataIndex: 'Title', width: '20%' },
    { title: '指标描述', dataIndex: 'Description', width: '25%' },
    { title: '目标/衡量标准', dataIndex: 'Target', width: '25%' },
    { title: '权重 (%)', dataIndex: 'Weight', width: '5%' },
    {
      title: '完成情况',
      dataIndex: 'CompletionDetails',
      render: (_, record) => (
        <Form.Item name={`completion_${record.ID}`} noStyle><Input.TextArea rows={2} disabled={isReadOnly} /></Form.Item>
      ),
    },
    {
      title: '得分',
      dataIndex: 'Score',
      width: '5%',
      render: (_, record) => (
        <Form.Item name={`score_${record.ID}`} noStyle rules={[{ type: 'number', min: 0, max: 120, message: '分数需在0-120之间' }]}><InputNumber min={0} max={120} disabled={isReadOnly} /></Form.Item>
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
    <div style={{ padding: '24px', background: '#f0f2f5' }}>
      <Card>
        <Title level={2} style={{ textAlign: 'center', marginBottom: '24px' }}>
          {isReadOnly ? '月度绩效详情' : '月度绩效打分'}
        </Title>
        <Form form={form} layout="vertical" onFinish={onFinish} onValuesChange={handleValuesChange}>
          <Descriptions bordered column={{ xxl: 4, xl: 3, lg: 3, md: 3, sm: 2, xs: 1 }}>
            <Descriptions.Item label="姓名">{review.User?.Name || 'N/A'}</Descriptions.Item>
            <Descriptions.Item label="部门">{review.User?.Department?.Name || 'N/A'}</Descriptions.Item>
            <Descriptions.Item label="岗位">{review.User?.Role?.Name || 'N/A'}</Descriptions.Item>
            <Descriptions.Item label="绩效周期">{review.Period}</Descriptions.Item>
            <Descriptions.Item label="总分">{dynamicTotalScore.toFixed(2)}</Descriptions.Item>
            <Descriptions.Item label="绩点">{dynamicGradePoint.toFixed(2)}</Descriptions.Item>
          </Descriptions>

          <Title level={4} style={{ marginTop: '24px' }}>绩效项详情</Title>
          <Table columns={columns} dataSource={review.Items} pagination={false} rowKey="ID" />

          <Title level={4} style={{ marginTop: '24px' }}>最终评语</Title>
          <Form.Item name="finalComment">
            <Input.TextArea rows={4} placeholder={isReadOnly ? "" : "请输入对本次绩效的最终评语"} disabled={isReadOnly} />
          </Form.Item>

          {!isReadOnly && (
            <Form.Item style={{ marginTop: 24, textAlign: 'center' }}>
              <Button type="primary" htmlType="submit" loading={loading} size="large">提交分数</Button>
            </Form.Item>
          )}

          {review && review.Status !== '草稿' && (
            <Form.Item style={{ marginTop: 24, textAlign: 'center' }}>
              <Button onClick={handlePrint} size="large" style={{ marginRight: '10px' }}>打印</Button>
              <Button onClick={handleExport} size="large">导出</Button>
            </Form.Item>
          )}
        </Form>
      </Card>
    </div>
  );
};

export default ScorePerformancePage;
