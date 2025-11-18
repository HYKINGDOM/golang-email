import { useQuery, useMutation } from '@tanstack/react-query'
import { Table, Button, Modal, Form, Input, Select, DatePicker } from 'antd'
import api from '../services/api'

export default function Tasks() {
  const { data: configs } = useQuery({ queryKey: ['configs'], queryFn: async () => (await api.get('/email-configs')).data })
  const { data: templates } = useQuery({ queryKey: ['templates'], queryFn: async () => (await api.get('/templates')).data, enabled: false })
  const create = useMutation({ mutationFn: async (payload: any) => (await api.post('/tasks', payload)).data })
  return (
    <div>
      <Button type="primary" onClick={() => Modal.confirm({
        title: '创建任务',
        content: <Form id="taskForm" layout="vertical">
          <Form.Item name="name" label="任务名称" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="sender_configs" label="发件邮箱" rules={[{ required: true }]}><Select mode="multiple" options={(configs||[]).map((c:any)=>({label:c.username,value:c.id}))} /></Form.Item>
          <Form.Item name="recipient_list" label="收件人" rules={[{ required: true }]}><Select mode="tags" placeholder="输入邮箱，回车确认" /></Form.Item>
          <Form.Item name="template_id" label="模板" rules={[{ required: true }]}><Select options={(templates||[]).map((t:any)=>({label:t.name,value:t.id}))} /></Form.Item>
          <Form.Item name="scheduled_time" label="定时发送"><DatePicker showTime /></Form.Item>
        </Form>,
        onOk: async () => {
          const form = (document.getElementById('taskForm') as any)?.__INTERNAL__.formStore
          const values = await form.validateFields(); await create.mutateAsync(values)
        }
      })}>创建任务</Button>
      <Table rowKey="id" dataSource={[]} columns={[{ title: '名称', dataIndex: 'name' }]} />
    </div>
  )
}