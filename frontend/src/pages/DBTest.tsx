import { Card, Button, Table } from 'antd'
import { useState } from 'react'
import api from '../services/api'

export default function DBTest() {
  const [ping, setPing] = useState<any>(null)
  const [perf, setPerf] = useState<any>(null)
  const [tx, setTx] = useState<any>(null)
  const [conc, setConc] = useState<any>(null)
  const runAll = async () => {
    const r1 = await api.get('/db/ping'); setPing(r1.data)
    const r2 = await api.get('/db/perf'); setPerf(r2.data)
    const r3 = await api.post('/db/transaction'); setTx(r3.data)
    const r4 = await api.get('/db/concurrency'); setConc(r4.data)
  }
  return (
    <div className="grid grid-cols-2 gap-4">
      <Card title="数据库连通性">
        <Button onClick={runAll} type="primary">运行所有测试</Button>
        <div className="mt-4">{ping && <pre>{JSON.stringify(ping,null,2)}</pre>}</div>
      </Card>
      <Card title="查询性能">{perf && <pre>{JSON.stringify(perf,null,2)}</pre>}</Card>
      <Card title="事务测试">{tx && <pre>{JSON.stringify(tx,null,2)}</pre>}</Card>
      <Card title="并发测试">{conc && <pre>{JSON.stringify(conc,null,2)}</pre>}</Card>
    </div>
  )
}

