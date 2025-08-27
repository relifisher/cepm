import React, { useState, useRef, useEffect } from 'react';
import { Form, Input, Button, DatePicker, Table, InputNumber, Popconfirm, message, Descriptions, Card, Space } from 'antd';
import { PlusOutlined, DeleteOutlined, PrinterOutlined } from '@ant-design/icons';
import { useReactToPrint } from 'react-to-print';
import dayjs from 'dayjs';
import { 
  createPerformanceReview, 
  updatePerformanceReview, 
  getReviewByPeriod
} from '../services/api';
import { useOutletContext } from 'react-router-dom';

const PerformancePlanPage = () => {
  const [form] = Form.useForm();
  const [workItems, setWorkItems] = useState([]);
  const [loading, setLoading] = useState(false);
  const [counter, setCounter] = useState(1);
  const [currentWorkWeight, setCurrentWorkWeight] = useState(0);
  const [weightValidateStatus, setWeightValidateStatus] = useState(''); // '' or 'error'

  const [activeReview, setActiveReview] = useState(null); // Holds the review being edited/viewed
  const [isReadOnly, setIsReadOnly] = useState(false); // Controls form editability

  const { currentUserId, currentUser } = useOutletContext();
  const componentRef = useRef();

  // Effect to update the form when a review is loaded or cleared
  useEffect(() => {
    if (activeReview) {
      form.setFieldsValue({ period: dayjs(activeReview.Period, 'YYYY-MM') });
      const items = activeReview.Items || [];
      const work = items.filter(item => item.Category === '工作业绩').map((item, index) => ({ ...item, key: item.ID || `loaded-${index}` }));
      const total = work.reduce((sum, item) => sum + (item.Weight || 0), 0);
      
      setWorkItems(work);
      setCurrentWorkWeight(total);

      const isEditable = activeReview.Status === '草稿' || activeReview.Status === '已驳回';
      setIsReadOnly(!isEditable);

      if (isEditable) {
        message.success(`加载了 ${activeReview.Period} 的可编辑记录`);
      } else {
        message.info(`'${activeReview.Period}' 的评估已提交，仅可查看。`);
      }
    } else {
      // Clear the table when there is no active review
      setWorkItems([]);
      setCurrentWorkWeight(0);
      setIsReadOnly(false);
    }
  }, [activeReview, form]);

  const handlePrint = useReactToPrint({
    content: () => componentRef.current,
    documentTitle: '月度绩效考核表',
  });

  const globalLargeModelItem = { key: 'global-1', Title: '大模型使用能力', Description: '衡量员工利用公司引入的大模型工具提升工作效率的能力', Weight: 10, Category: '大模型' };
  const globalValuesItem = { key: 'global-2', Title: '价值观践行', Description: '评估员工在工作中对公司价值观的理解和实践程度', Weight: 10, Category: '价值观' };

  const handleMonthChange = async (date) => {
    if (!date) {
      setActiveReview(null);
      form.resetFields();
      return;
    }

    const period = date.format('YYYY-MM');
    setLoading(true);

    try {
        const response = await getReviewByPeriod(currentUserId, period);
        // Check if the response data is an empty object (indicating no record found)
        if (Object.keys(response.data).length === 0) {
          setActiveReview(null); // Treat as no active review, clear form
          form.resetFields(); // Reset all form fields
          form.setFieldsValue({ period: date }); // Re-set the selected date in the picker
        } else {
          setActiveReview(response.data);
        }
      } catch (error) {
        // This catch block will now only handle actual errors, not 404s for "not found"
        message.error('查询评估记录时出错');
      } finally {
      setLoading(false);
    }
  };

  const handleAddItem = () => {
    if (workItems.length >= 10) {
      message.warning('最多只能添加10个业绩考核项。');
      return;
    }
    const newItem = { key: counter, Title: '', Description: '', Weight: null, Target: '' };
    setWorkItems([...workItems, newItem]);
    setCounter(counter + 1);
  };

  const handleDeleteItem = (key) => {
    setWorkItems(workItems.filter(item => item.key !== key));
  };

  const handleItemChange = (key, dataIndex, value) => {
    if (isReadOnly) return;
    if (dataIndex === 'Weight') setWeightValidateStatus('');
    const newItems = workItems.map(item => item.key === key ? { ...item, [dataIndex]: value } : item);
    setWorkItems(newItems);
    // Recalculate weight sum on change
    if (dataIndex === 'Weight') {
      const total = newItems.reduce((sum, item) => sum + (item.Weight || 0), 0);
      setCurrentWorkWeight(total);
    }
  };

  const handleSave = async (status) => {
    if (isReadOnly) return;

    let values;
    try {
      values = await form.validateFields();
    } catch (errorInfo) {
      console.log('表单校验失败:', errorInfo);
      return;
    }

    for (const item of workItems) {
      if (!item.Title || !item.Description || !item.Target || !item.Weight) {
        message.error('所有绩效项的字段均不能为空，请填写完整。');
        return;
      }
    }
    const workTotalWeight = workItems.reduce((sum, item) => sum + (item.Weight || 0), 0);
    if (workTotalWeight !== 80) {
      message.error(`“工作业绩”部分的各项权重之和必须等于 80%，当前为 ${workTotalWeight}%。`);
      setWeightValidateStatus('error');
      return;
    }

    setLoading(true);

    const finalWorkItems = workItems.map(({ key, ...rest }) => ({ ...rest, Category: '工作业绩' }));
    const finalGlobalItems = [globalLargeModelItem, globalValuesItem].map(({ key, ...rest }) => ({...rest, Target: rest.Description, Title: rest.Title, Weight: rest.Weight, Description: rest.Description, Category: rest.Category}));
    const reviewData = { UserID: currentUserId, period: values.period.format('YYYY-MM'), status: status, items: [...finalWorkItems, ...finalGlobalItems] };

    try {
      if (activeReview && activeReview.ID) {
        await updatePerformanceReview(activeReview.ID, { ...reviewData, ID: activeReview.ID });
        message.success(`绩效评估${status === '草稿' ? '保存' : '提交'}成功!`);
      } else {
        await createPerformanceReview(reviewData);
        message.success(`绩效评估${status === '草稿' ? '创建' : '提交'}成功!`);
      }
      form.resetFields();
      setActiveReview(null);
    } catch (error) {
      const errorMsg = error.response?.data?.error || (activeReview ? '更新失败' : '创建失败');
      message.error(errorMsg);
    } finally {
      setLoading(false);
    }
  };

  const handleSaveDraft = () => handleSave('草稿');
  const handleSubmit = () => handleSave('待审批');

  const workColumns = [
    { title: '考核指标 (KPI)', dataIndex: 'Title', width: '20%', render: (_, record) => <Input value={record.Title} onChange={e => handleItemChange(record.key, 'Title', e.target.value)} placeholder="例如：完成XX功能模块" disabled={isReadOnly} /> },
    { title: '指标描述', dataIndex: 'Description', width: '40%', render: (_, record) => <Input.TextArea value={record.Description} onChange={e => handleItemChange(record.key, 'Description', e.target.value)} placeholder="指标的详细描述" disabled={isReadOnly} /> },
    {
      title: '权重 (%)',
      dataIndex: 'Weight',
      width: '10%',
      render: (_, record) => (
        <Form.Item noStyle validateStatus={weightValidateStatus}>
          <InputNumber min={0} max={80} value={record.Weight} onChange={value => handleItemChange(record.key, 'Weight', value)} disabled={isReadOnly} />
        </Form.Item>
      ),
    },
    { title: '目标/衡量标准', dataIndex: 'Target', width: '25%', render: (_, record) => <Input.TextArea value={record.Target} onChange={e => handleItemChange(record.key, 'Target', e.target.value)} placeholder="例如：月底前上线" disabled={isReadOnly} /> },
    { title: '操作', dataIndex: 'action', width: '5%', render: (_, record) => !isReadOnly && <Popconfirm title="确认删除?" onConfirm={() => handleDeleteItem(record.key)}><Button icon={<DeleteOutlined />} danger className="no-print" /></Popconfirm> },
  ];

  const globalColumns = [
    { title: '考核指标', dataIndex: 'Title', width: '20%' },
    { title: '指标描述', dataIndex: 'Description', width: '40%' },
    { title: '权重 (%)', dataIndex: 'Weight', width: '10%' },
    // Placeholders for alignment
    { title: '', dataIndex: 'placeholder1', width: '25%', render: () => null },
    { title: '', dataIndex: 'placeholder2', width: '5%', render: () => null },
  ];

  return (
    <div style={{ padding: '24px', background: '#f0f2f5' }}>
      <Card ref={componentRef}>
        <h1 style={{ textAlign: 'center', marginBottom: '24px' }}>月度绩效考核表</h1>
        <Form form={form} layout="vertical">
          <Descriptions bordered column={{ xxl: 4, xl: 3, lg: 3, md: 3, sm: 2, xs: 1 }}>
            <Descriptions.Item label="姓名">{currentUser.name}</Descriptions.Item>
            <Descriptions.Item label="部门">{currentUser.department}</Descriptions.Item>
            <Descriptions.Item label="岗位">{currentUser.role}</Descriptions.Item>
            <Descriptions.Item label="绩效周期">
              <Form.Item name="period" noStyle rules={[{ required: true, message: '请选择绩效周期!' }]}>
                <DatePicker picker="month" onChange={handleMonthChange} />
              </Form.Item>
            </Descriptions.Item>
          </Descriptions>

          <h2 style={{ marginTop: '24px' }}>一、工作业绩 (当前总和: {currentWorkWeight}% / 要求: 80%)</h2>
          {!isReadOnly && <Button onClick={handleAddItem} type="dashed" style={{ marginBottom: 16 }} icon={<PlusOutlined />} className="no-print">添加业绩考核项</Button>}
          <Table columns={workColumns} dataSource={workItems} pagination={false} rowKey="key" />

          <h2 style={{ marginTop: '24px' }}>二、大模型 (权重 10%)</h2>
          <Table columns={globalColumns} dataSource={[globalLargeModelItem]} pagination={false} rowKey="key" />

          <h2 style={{ marginTop: '24px' }}>三、价值观 (权重 10%)</h2>
          <Table columns={globalColumns} dataSource={[globalValuesItem]} pagination={false} rowKey="key" />

          <Form.Item style={{ marginTop: 24, textAlign: 'center' }} className="no-print">
            <Space>
              {!isReadOnly && (
                <>
                  <Button onClick={handleSaveDraft} loading={loading} size="large">保存草稿</Button>
                  <Button type="primary" onClick={handleSubmit} loading={loading} size="large">提交填报</Button>
                </>
              )}
              {isReadOnly && (
                <Button icon={<PrinterOutlined />} onClick={handlePrint} size="large">打印/导出PDF</Button>
              )}
            </Space>
          </Form.Item>
        </Form>
      </Card>
    </div>
  );
};

export default PerformancePlanPage;
