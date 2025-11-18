import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Table, Button, Modal, Form, Input, InputNumber, Switch } from 'antd'
import api from '../services/api'

export default function EmailConfigs() {
  const qc = useQueryClient()
  const { data, isLoading } = useQuery({ queryKey: ['configs'], queryFn: async () => (await api.get('/email-configs')).data })
  const create = useMutation({ mutationFn: async (payload: any) => (await api.post('/email-configs', payload)).data, onSuccess: () => qc.invalidateQueries({ queryKey: ['configs'] }) })
  return (
    <div>
      <Button type="primary" onClick={() => Modal.confirm({
        title: '新建邮箱配置',
        content: <Form id="configForm" layout="vertical">
          <Form.Item name="provider" label="服务商" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="host" label="SMTP主机" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="port" label="端口" rules={[{ required: true }]}><InputNumber style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="username" label="用户名" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="password" label="密码" rules={[{ required: true }]}><Input.Password /></Form.Item>
          <Form.Item name="daily_limit" label="每日上限"><InputNumber style={{ width: '100%' }} /></Form.Item>
          <Form.Item name="is_active" label="启用" valuePropName="checked"><Switch /></Form.Item>
        </Form>,
        onOk: async () => {
          const form = (document.getElementById('configForm') as any)?.__INTERNAL__.formStore
          const values = await form.validateFields(); await create.mutateAsync(values)
        }
      })}>新建</Button>
      <Table loading={isLoading} rowKey="id" dataSource={data || []} columns={[
        { title: '邮箱', dataIndex: 'username' },
        { title: '主机', dataIndex: 'host' },
        { title: '端口', dataIndex: 'port' },
        { title: '状态', dataIndex: 'is_active', render: (v: boolean) => v ? '启用' : '停用' },
      ]} />
    </div>
  )
}