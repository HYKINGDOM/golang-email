import { useEffect, useState } from 'react'
import { Card, Table, Tag, Button, InputNumber } from 'antd'
import api from '../services/api'
import { getBaseURL, getSSEPath } from '../config'

type Event = { task_id: number, type: string, data: any }

export default function Monitoring() {
  const [taskId, setTaskId] = useState<number>(1)
  const [items, setItems] = useState<any[]>([])
  useEffect(()=>{
    const es = new EventSource(`${getBaseURL()}${getSSEPath()}/${taskId}`)
    es.onmessage = (e)=>{
      const evt: Event = JSON.parse(e.data)
      if (evt.type === 'progress') {
        setItems(prev => [{ recipient: evt.data.recipient, status: evt.data.status, error: evt.data.error }, ...prev].slice(0,100))
      }
    }
    return ()=> es.close()
  },[taskId])
  return (
    <div className="grid grid-cols-2 gap-4">
      <Card title="任务实时监控">
        <div className="mb-4">任务ID：<InputNumber value={taskId} onChange={(v)=> setTaskId(Number(v))} /></div>
        <Table rowKey={(r)=>r.recipient+r.status} dataSource={items} columns={[
          { title: '收件人', dataIndex: 'recipient' },
          { title: '状态', dataIndex: 'status', render: (v:string)=> v==='failed'? <Tag color="red">失败</Tag>:<Tag color="green">成功</Tag> },
          { title: '错误', dataIndex: 'error' },
          { title: '操作', render: (_:any, r:any)=> r.status==='failed' && r.id ? <Button onClick={()=> api.post('/records/'+r.id+'/retry')}>重试</Button> : null }
        ]} />
      </Card>
    </div>
  )
}
