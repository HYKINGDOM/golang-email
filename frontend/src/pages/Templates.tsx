import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { Table, Button, Modal, Form, Input, Switch } from 'antd'
import ReactQuill from 'react-quill'
import 'react-quill/dist/quill.snow.css'
import api from '../services/api'

export default function Templates() {
  const qc = useQueryClient()
  const { data, isLoading } = useQuery({ queryKey: ['templates'], queryFn: async () => (await api.get('/templates')).data })
  const create = useMutation({ mutationFn: async (payload: any) => (await api.post('/templates', payload)).data, onSuccess: () => qc.invalidateQueries({ queryKey: ['templates'] }) })
  let content = ''
  return (
    <div>
      <Button type="primary" onClick={() => Modal.confirm({
        title: '新建模板',
        content: <Form id="tplForm" layout="vertical">
          <Form.Item name="name" label="名称" rules={[{ required: true }]}><Input /></Form.Item>
          <Form.Item name="subject" label="主题" rules={[{ required: true }]}><Input /></Form.Item>
          <div className="mb-4"><ReactQuill onChange={(v)=>{content=v}} /></div>
          <Form.Item name="is_rich_text" label="富文本" valuePropName="checked"><Switch defaultChecked /></Form.Item>
          <Form.Item name="tracking_enabled" label="追踪" valuePropName="checked"><Switch /></Form.Item>
        </Form>,
        onOk: async () => {
          const form = (document.getElementById('tplForm') as any)?.__INTERNAL__.formStore
          const values = await form.validateFields(); await create.mutateAsync({ ...values, content })
        }
      })}>新建</Button>
      <Table loading={isLoading} rowKey="id" dataSource={data || []} columns={[
        { title: '名称', dataIndex: 'name' },
        { title: '主题', dataIndex: 'subject' },
        { title: '追踪', dataIndex: 'tracking_enabled', render: (v:boolean)=> v? '开启':'关闭' },
      ]} />
    </div>
  )
}