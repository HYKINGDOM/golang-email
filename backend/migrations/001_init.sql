-- 邮箱配置表
CREATE TABLE IF NOT EXISTS email_configs (
    id SERIAL PRIMARY KEY,
    provider VARCHAR(50) NOT NULL,
    host VARCHAR(255) NOT NULL,
    port INT NOT NULL,
    username VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    daily_limit INT DEFAULT 200,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- 邮件模板表
CREATE TABLE IF NOT EXISTS email_templates (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    subject VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    is_rich_text BOOLEAN DEFAULT TRUE,
    tracking_enabled BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- 发送任务表
CREATE TABLE IF NOT EXISTS send_tasks (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    sender_configs TEXT NOT NULL,
    recipient_list TEXT NOT NULL,
    template_id INT NOT NULL,
    scheduled_time TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- 发送记录表
CREATE TABLE IF NOT EXISTS send_records (
    id SERIAL PRIMARY KEY,
    task_id INT NOT NULL,
    sender_email VARCHAR(255) NOT NULL,
    recipient_email VARCHAR(255) NOT NULL,
    send_time TIMESTAMP,
    status VARCHAR(20) DEFAULT 'pending',
    error_message TEXT,
    retry_count INT DEFAULT 0,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

-- 邮件追踪表
CREATE TABLE IF NOT EXISTS email_tracking (
    id SERIAL PRIMARY KEY,
    record_id INT NOT NULL,
    open_time TIMESTAMP,
    open_count INT DEFAULT 0,
    read_duration INT DEFAULT 0,
    last_open_time TIMESTAMP,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);