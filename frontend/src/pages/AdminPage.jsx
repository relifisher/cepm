import React, { useState, useEffect } from 'react';
import { Card, Tabs, Typography, Table, Button, Modal, Form, Input, Select, message, Switch, Tree } from 'antd';
import { EditOutlined } from '@ant-design/icons';
import { getUsers, updateUser, createDepartment, updateSystemSetting, getDepartments, getRoles } from '../services/adminApi';

const { Title } = Typography;
const { TabPane } = Tabs;
const { Option } = Select;

const AdminPage = () => {
  const [users, setUsers] = useState([]);
  const [loading, setLoading] = useState(false);
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [editingUser, setEditingUser] = useState(null);
  const [form] = Form.useForm();

  const [departments, setDepartments] = useState([]);
  const [departmentForm] = Form.useForm();

  const [roles, setRoles] = useState([]); // New state for roles

  useEffect(() => {
    fetchUsers();
    fetchDepartments();
    fetchRoles(); // Fetch roles on component mount
  }, []);

  const fetchUsers = async () => {
    setLoading(true);
    try {
      const data = await getUsers();
      setUsers(data);
    } catch (error) {
      message.error('Failed to fetch users.');
    } finally {
      setLoading(false);
    }
  };

  const fetchDepartments = async () => {
    try {
      const data = await getDepartments();
      setDepartments(data);
    } catch (error) {
      message.error('Failed to fetch departments.');
    }
  };

  const fetchRoles = async () => {
    try {
      const data = await getRoles();
      setRoles(data);
    } catch (error) {
      message.error('Failed to fetch roles.');
    }
  };

  const handleEdit = (user) => {
    setEditingUser(user);
    form.setFieldsValue(user);
    setIsModalVisible(true);
  };

  const handleModalOk = async () => {
    try {
      const values = await form.validateFields();
      const updatedUser = { ...editingUser, ...values };
      await updateUser(updatedUser.ID, updatedUser);
      message.success('User updated successfully!');
      setIsModalVisible(false);
      fetchUsers(); // Refresh user list
    } catch (error) {
      message.error('Failed to update user.');
    }
  };

  const handleModalCancel = () => {
    setIsModalVisible(false);
    setEditingUser(null);
    form.resetFields();
  };

  // Helper to build tree data from flat list
  const buildDepartmentTree = (data, parentId = null) => {
    return data
      .filter(item => item.ParentID === parentId)
      .map(item => ({
        key: item.ID,
        title: item.Name,
        children: buildDepartmentTree(data, item.ID),
      }));
  };

  const handleCreateDepartment = async (values) => {
    try {
      await createDepartment({ Name: values.departmentName, ParentID: values.parentDepartmentId || null });
      message.success('Department created successfully!');
      departmentForm.resetFields();
      fetchDepartments(); // Refresh department list
    } catch (error) {
      message.error('Failed to create department.');
    }
  };

  const handleUpdateSystemSetting = async (key, value) => {
    try {
      await updateSystemSetting({ Key: key, Value: value });
      message.success('System setting updated successfully!');
    } catch (error) {
      message.error('Failed to update system setting.');
    }
  };

  const userColumns = [
    { title: 'ID', dataIndex: 'ID', key: 'ID' },
    { title: '姓名', dataIndex: 'Name', key: 'Name' },
    { title: '邮箱', dataIndex: 'Email', key: 'Email' },
    { title: '角色', dataIndex: ['Role', 'Name'], key: 'RoleName' },
    { title: '部门', dataIndex: ['Department', 'Name'], key: 'DepartmentName' },
    { title: '操作', key: 'actions', render: (_, record) => (
      <Button icon={<EditOutlined />} onClick={() => handleEdit(record)}>编辑</Button>
    )},
  ];

  return (
    <div style={{ padding: '24px' }}>
      <Title level={2}>管理员设置</Title>
      <Card>
        <Tabs defaultActiveKey="1">
          <TabPane tab="用户管理" key="1">
            <Table
              columns={userColumns}
              dataSource={users}
              loading={loading}
              rowKey="ID"
              pagination={{ pageSize: 10 }}
            />
            <Modal
              title="编辑用户"
              visible={isModalVisible}
              onOk={handleModalOk}
              onCancel={handleModalCancel}
            >
              <Form form={form} layout="vertical">
                <Form.Item name="Name" label="姓名" rules={[{ required: true, message: '请输入姓名' }]}>
                  <Input />
                </Form.Item>
                <Form.Item name="Email" label="邮箱" rules={[{ required: true, message: '请输入邮箱' }]}>
                  <Input />
                </Form.Item>
                <Form.Item name="RoleID" label="角色" rules={[{ required: true, message: '请选择角色' }]}>
                  <Select placeholder="选择角色">
                    {roles.map(role => (
                      <Option key={role.ID} value={role.ID}>
                        {role.Name}
                      </Option>
                    ))}
                  </Select>
                </Form.Item>
                <Form.Item name="DepartmentID" label="部门">
                  <Select placeholder="选择部门 (可选)" allowClear>
                    {departments.map(dept => (
                      <Option key={dept.ID} value={dept.ID}>
                        {dept.Name}
                      </Option>
                    ))}
                  </Select>
                </Form.Item>
              </Form>
            </Modal>
          </TabPane>
          <TabPane tab="组织架构管理" key="2">
            <Form form={departmentForm} layout="inline" onFinish={handleCreateDepartment} style={{ marginBottom: '20px' }}>
              <Form.Item
                name="departmentName"
                rules={[{ required: true, message: '请输入部门名称' }]}
              >
                <Input placeholder="部门名称" />
              </Form.Item>
              <Form.Item name="parentDepartmentId">
                <Select
                  placeholder="选择上级部门 (可选)"
                  style={{ width: 200 }}
                  allowClear
                >
                  {departments.map(dept => (
                    <Option key={dept.ID} value={dept.ID}>
                      {dept.Name}
                    </Option>
                  ))}
                </Select>
              </Form.Item>
              <Form.Item>
                <Button type="primary" htmlType="submit">创建部门</Button>
              </Form.Item>
            </Form>
            <Tree
              showLine
              defaultExpandAll
              treeData={buildDepartmentTree(departments)}
            />
          </TabPane>
          <TabPane tab="系统设置" key="3">
            <Form layout="vertical">
              <Form.Item label="限制用户只能填报当前月份">
                <Switch
                  checkedChildren="开启"
                  unCheckedChildren="关闭"
                  // You'd typically fetch the initial state from backend
                  onChange={(checked) => handleUpdateSystemSetting('restrict_current_month_only', checked ? 'true' : 'false')}
                />
              </Form.Item>
              <Button type="primary" onClick={() => message.info('System setting update logic not fully implemented yet.')}>保存设置</Button>
            </Form>
          </TabPane>
        </Tabs>
      </Card>
    </div>
  );
};

export default AdminPage;
