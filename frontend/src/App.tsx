import { Layout, Menu } from 'antd'
import { BrowserRouter, Routes, Route, Link } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import Dashboard from './pages/Dashboard'
import EmailConfigs from './pages/EmailConfigs'
import Tasks from './pages/Tasks'
import Templates from './pages/Templates'
import Monitoring from './pages/Monitoring'
import DBTest from './pages/DBTest'
import './index.css'

const qc = new QueryClient()

export default function App() {
  return (
    <QueryClientProvider client={qc}>
      <BrowserRouter>
        <Layout style={{ minHeight: '100vh' }}>
          <Layout.Header>
            <Menu theme="dark" mode="horizontal">
              <Menu.Item key="dashboard"><Link to="/">仪表盘</Link></Menu.Item>
              <Menu.Item key="configs"><Link to="/configs">邮箱配置</Link></Menu.Item>
              <Menu.Item key="tasks"><Link to="/tasks">任务管理</Link></Menu.Item>
              <Menu.Item key="templates"><Link to="/templates">模板管理</Link></Menu.Item>
              <Menu.Item key="monitor"><Link to="/monitor">实时监控</Link></Menu.Item>
              <Menu.Item key="dbtest"><Link to="/dbtest">数据库测试</Link></Menu.Item>
            </Menu>
          </Layout.Header>
          <Layout.Content className="p-6">
            <Routes>
              <Route path="/" element={<Dashboard />} />
              <Route path="/configs" element={<EmailConfigs />} />
              <Route path="/tasks" element={<Tasks />} />
              <Route path="/templates" element={<Templates />} />
              <Route path="/monitor" element={<Monitoring />} />
              <Route path="/dbtest" element={<DBTest />} />
            </Routes>
          </Layout.Content>
        </Layout>
      </BrowserRouter>
    </QueryClientProvider>
  )
}
