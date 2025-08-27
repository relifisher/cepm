-- CEPM (Corporate Employee Performance Management) Database Schema
-- Target DBMS: PostgreSQL

-- 部门表 (Departments)
-- 存储公司的组织架构
CREATE TABLE departments (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    parent_id INTEGER REFERENCES departments(id) ON DELETE SET NULL, -- 父部门ID，形成树状结构
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE departments IS '部门表，存储公司组织架构';
COMMENT ON COLUMN departments.parent_id IS '父部门ID，顶级部门为NULL';

-- 角色表 (Roles)
-- 定义系统中的不同角色及其权限
CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE, -- 例如: "员工", "组长", "中心负责人", "人事"
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE roles IS '角色表';
COMMENT ON COLUMN roles.name IS '角色名称，唯一';

-- 员工表 (Users)
-- 存储员工基本信息，并关联部门、角色和上级
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    wechat_userid VARCHAR(255) UNIQUE, -- 企业微信的UserID，用于单点登录
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    avatar TEXT, -- 头像URL
    department_id INTEGER REFERENCES departments(id) ON DELETE SET NULL,
    role_id INTEGER REFERENCES roles(id) ON DELETE SET NULL,
    manager_id INTEGER REFERENCES users(id) ON DELETE SET NULL, -- 直属上级ID
    is_active BOOLEAN NOT NULL DEFAULT TRUE, -- 账号是否激活
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE users IS '员工信息表';
COMMENT ON COLUMN users.wechat_userid IS '企业微信的UserID，用于单点登录和身份识别';
COMMENT ON COLUMN users.manager_id IS '直属上级的ID，用于构建汇报关系';

-- 月度绩效评估主表 (Performance Reviews)
-- 每次绩效评估的核心记录
CREATE TABLE performance_reviews (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    period VARCHAR(7) NOT NULL, -- 绩效周期，格式 "YYYY-MM"
    status VARCHAR(50) NOT NULL DEFAULT 'Draft', -- Draft, PendingApproval, Approved, Evaluating, Completed, Rejected
    total_score NUMERIC(5, 2), -- 最终总分
    final_comment TEXT, -- 最终评语
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(user_id, period) -- 每个员工每个月只能有一份绩效
);
COMMENT ON TABLE performance_reviews IS '月度绩效评估主表';
COMMENT ON COLUMN performance_reviews.period IS '绩效周期，格式 YYYY-MM';
COMMENT ON COLUMN performance_reviews.status IS '绩效状态：草稿、待审批、已批准、评分中、已完成、已驳回';

-- 绩效评估项表 (Performance Items)
-- 具体的绩效指标（KPI）
CREATE TABLE performance_items (
    id SERIAL PRIMARY KEY,
    review_id INTEGER NOT NULL REFERENCES performance_reviews(id) ON DELETE CASCADE,
    category VARCHAR(50) NOT NULL DEFAULT '工作业绩', -- 新增的 category 列
    title VARCHAR(255) NOT NULL, -- 指标名称
    description TEXT, -- 指标的详细描述
    weight NUMERIC(5, 2) NOT NULL, -- 权重 (例如: 20.00 表示 20%)
    target TEXT, -- 目标或衡量标准
    completion_details TEXT, -- 实际完成情况
    score NUMERIC(5, 2), -- 单项得分
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE performance_items IS '具体的绩效指标项';
COMMENT ON COLUMN performance_items.weight IS '权重，所有项权重之和应为100';

-- 审批流转历史表 (Approval History)
-- 记录每次绩效评估的审批过程
CREATE TABLE approval_history (
    id SERIAL PRIMARY KEY,
    review_id INTEGER NOT NULL REFERENCES performance_reviews(id) ON DELETE CASCADE,
    approver_id INTEGER NOT NULL REFERENCES users(id),
    status VARCHAR(50) NOT NULL, -- 例如: "Approved", "Rejected"
    comment TEXT, -- 审批意见
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
COMMENT ON TABLE approval_history IS '审批流转历史记录';
COMMENT ON COLUMN approval_history.status IS '审批结果状态';

-- 创建索引以提高查询性能
CREATE INDEX idx_users_department_id ON users(department_id);
CREATE INDEX idx_reviews_user_period ON performance_reviews(user_id, period);
CREATE INDEX idx_items_review_id ON performance_items(review_id);


-- Initial data for CEPM system

-- Insert Roles
INSERT INTO roles (name, description) VALUES
('组员', '普通员工'),
('组长', '团队负责人'),
('总监', '部门或中心负责人'),
('人事', '人力资源部成员');

-- Insert Departments (forming a hierarchy)
INSERT INTO departments (id, name, parent_id) VALUES
(1, '公司总部', NULL),
(2, '人力资源部', 1),
(3, '研发中心', 1),
(4, '研发一部', 3),
(5, '研发二部', 3),
(6, '市场部', 1),
(7, '销售部', 1),
(8, '财务部', 1),
(9, '行政部', 1),
(10, '法务部', 1),
(11, '产品部', 1);

-- Insert Users (covering roles and departments)
-- Passwords are not handled here, as authentication will be via WeChat Work

-- CEO (总监, 公司总部)
INSERT INTO users (name, email, role_id, department_id) VALUES
('张总', 'zhangzong@example.com', (SELECT id FROM roles WHERE name = '总监'), (SELECT id FROM departments WHERE name = '公司总部'));

-- HR Manager (人事, 人力资源部)
INSERT INTO users (name, email, role_id, department_id, manager_id) VALUES
('李人事', 'lirenshi@example.com', (SELECT id FROM roles WHERE name = '人事'), (SELECT id FROM departments WHERE name = '人力资源部'), (SELECT id FROM users WHERE name = '张总'));

-- R&D Director (总监, 研发中心)
INSERT INTO users (name, email, role_id, department_id, manager_id) VALUES
('王总监', 'wangzongjian@example.com', (SELECT id FROM roles WHERE name = '总监'), (SELECT id FROM departments WHERE name = '研发中心'), (SELECT id FROM users WHERE name = '张总'));

-- R&D Team Lead 1 (组长, 研发一部)
INSERT INTO users (name, email, role_id, department_id, manager_id) VALUES
('赵组长', 'zhaozuzhang@example.com', (SELECT id FROM roles WHERE name = '组长'), (SELECT id FROM departments WHERE name = '研发一部'), (SELECT id FROM users WHERE name = '王总监'));

-- R&D Member 1 (组员, 研发一部)
INSERT INTO users (name, email, role_id, department_id, manager_id) VALUES
('钱成员', 'qianchengyuan@example.com', (SELECT id FROM roles WHERE name = '组员'), (SELECT id FROM departments WHERE name = '研发一部'), (SELECT id FROM users WHERE name = '赵组长'));

-- R&D Member 2 (组员, 研发一部)
INSERT INTO users (name, email, role_id, department_id, manager_id) VALUES
('孙成员', 'sunchengyuan@example.com', (SELECT id FROM roles WHERE name = '组员'), (SELECT id FROM departments WHERE name = '研发一部'), (SELECT id FROM users WHERE name = '赵组长'));

-- R&D Team Lead 2 (组长, 研发二部)
INSERT INTO users (name, email, role_id, department_id, manager_id) VALUES
('周组长', 'zhouzuzhang@example.com', (SELECT id FROM roles WHERE name = '组长'), (SELECT id FROM departments WHERE name = '研发二部'), (SELECT id FROM users WHERE name = '王总监'));

-- R&D Member 3 (组员, 研发二部)
INSERT INTO users (name, email, role_id, department_id, manager_id) VALUES
('吴成员', 'wuchengyuan@example.com', (SELECT id FROM roles WHERE name = '组员'), (SELECT id FROM departments WHERE name = '研发二部'), (SELECT id FROM users WHERE name = '周组长'));

-- Marketing Manager (组长, 市场部)
INSERT INTO users (name, email, role_id, department_id, manager_id) VALUES
('郑经理', 'zhengjingli@example.com', (SELECT id FROM roles WHERE name = '组长'), (SELECT id FROM departments WHERE name = '市场部'), (SELECT id FROM users WHERE name = '张总'));

-- Sales Member (组员, 销售部)
INSERT INTO users (name, email, role_id, department_id, manager_id) VALUES
('王销售', 'wangxiaoshou@example.com', (SELECT id FROM roles WHERE name = '组员'), (SELECT id FROM departments WHERE name = '销售部'), (SELECT id FROM users WHERE name = '郑经理'));

-- Example: A user without a manager (e.g., top-level or self-managed)
INSERT INTO users (name, email, role_id, department_id) VALUES
('陈独立', 'chenduli@example.com', (SELECT id FROM roles WHERE name = '组员'), (SELECT id FROM departments WHERE name = '产品部'));


-- Insert Performance Reviews and Items

-- Review 1: 钱成员 (qianchengyuan@example.com) - Draft
INSERT INTO performance_reviews (user_id, period, status, final_comment) VALUES
((SELECT id FROM users WHERE email = 'qianchengyuan@example.com'), '2025-07', '草稿', NULL);

INSERT INTO performance_items (review_id, category, title, description, weight, target) VALUES
((SELECT id FROM performance_reviews WHERE user_id = (SELECT id FROM users WHERE email = 'qianchengyuan@example.com') AND period = '2025-07'), '工作业绩', '完成项目A核心模块', '负责项目A的后端核心逻辑开发', 50, '模块功能通过所有单元测试'),
((SELECT id FROM performance_reviews WHERE user_id = (SELECT id FROM users WHERE email = 'qianchengyuan@example.com') AND period = '2025-07'), '工作业绩', '参与技术分享', '每月至少分享一次技术经验', 30, '完成2次内部技术分享'),
((SELECT id FROM performance_reviews WHERE user_id = (SELECT id FROM users WHERE email = 'qianchengyuan@example.com') AND period = '2025-07'), '大模型', '大模型工具应用', '在日常开发中积极尝试使用大模型工具提升效率', 10, '提交至少3个使用大模型工具的案例'),
((SELECT id FROM performance_reviews WHERE user_id = (SELECT id FROM users WHERE email = 'qianchengyuan@example.com') AND period = '2025-07'), '价值观', '团队协作', '积极与团队成员沟通协作，共同解决问题', 10, '获得至少3位同事的正面反馈');

-- Review 2: 孙成员 (sunchengyuan@example.com) - Completed
INSERT INTO performance_reviews (user_id, period, status, total_score, final_comment) VALUES
((SELECT id FROM users WHERE email = 'sunchengyuan@example.com'), '2025-06', '已完成', 85.5, '该员工表现优秀，超额完成任务。');

INSERT INTO performance_items (review_id, category, title, description, weight, target, completion_details, score) VALUES
((SELECT id FROM performance_reviews WHERE user_id = (SELECT id FROM users WHERE email = 'sunchengyuan@example.com') AND period = '2025-06'), '工作业绩', '完成项目B需求分析', '负责项目B的需求调研和文档编写', 40, '需求文档通过评审', '按时提交需求文档，并获得高层认可', 90),
((SELECT id FROM performance_reviews WHERE user_id = (SELECT id FROM users WHERE email = 'sunchengyuan@example.com') AND period = '2025-06'), '工作业绩', '优化系统性能', '提升核心模块响应速度', 40, '响应时间缩短20%', '响应时间缩短30%，效果显著', 95),
((SELECT id FROM performance_reviews WHERE user_id = (SELECT id FROM users WHERE email = 'sunchengyuan@example.com') AND period = '2025-06'), '大模型', '大模型学习与实践', '主动学习大模型相关知识并应用于工作', 10, '完成大模型课程学习', '完成课程学习并提交2个创新应用方案', 80),
((SELECT id FROM performance_reviews WHERE user_id = (SELECT id FROM users WHERE email = 'sunchengyuan@example.com') AND period = '2025-06'), '价值观', '客户导向', '积极响应客户需求，提供优质服务', 10, '客户满意度达到90%', '客户满意度达到95%，无客户投诉', 70);

-- Review 3: 吴成员 (wuchengyuan@example.com) - Pending Score (待打分)
INSERT INTO performance_reviews (user_id, period, status, final_comment) VALUES
((SELECT id FROM users WHERE email = 'wuchengyuan@example.com'), '2025-07', '待打分', NULL);

INSERT INTO performance_items (review_id, category, title, description, weight, target) VALUES
((SELECT id FROM performance_reviews WHERE user_id = (SELECT id FROM users WHERE email = 'wuchengyuan@example.com') AND period = '2025-07'), '工作业绩', '完成新功能开发', '负责新功能从设计到上线全流程', 60, '功能按时上线，无重大bug'),
((SELECT id FROM performance_reviews WHERE user_id = (SELECT id FROM users WHERE email = 'wuchengyuan@example.com') AND period = '2025-07'), '工作业绩', '参与代码评审', '积极参与团队代码评审，提升代码质量', 20, '每月至少评审10次'),
((SELECT id FROM performance_reviews WHERE user_id = (SELECT id FROM users WHERE email = 'wuchengyuan@example.com') AND period = '2025-07'), '大模型', '大模型知识分享', '组织一次大模型技术分享会', 10, '分享会参与人数超过20人'),
((SELECT id FROM performance_reviews WHERE user_id = (SELECT id FROM users WHERE email = 'wuchengyuan@example.com') AND period = '2025-07'), '价值观', '持续学习', '主动学习新知识，提升个人能力', 10, '完成2门在线课程学习');
