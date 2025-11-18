import { Card } from 'antd'
import { useQuery } from '@tanstack/react-query'
import { Bar } from 'react-chartjs-2'
import { Chart, BarElement, CategoryScale, LinearScale, Tooltip, Legend } from 'chart.js'
import api from '../services/api'
import { getChartRefreshMs } from '../config'

Chart.register(BarElement, CategoryScale, LinearScale, Tooltip, Legend)

export default function Dashboard() {
  const { data } = useQuery({ queryKey: ['stats'], queryFn: async () => (await api.get('/stats/send')).data, refetchInterval: getChartRefreshMs() })
  const chartData = {
    labels: ['总量', '成功', '失败'],
    datasets: [{ label: '发送统计', data: [data?.total||0, data?.sent||0, data?.failed||0], backgroundColor: ['#93c5fd','#86efac','#fca5a5'] }]
  }
  return (
    <div className="grid grid-cols-3 gap-4">
      <Card title="发送统计">
        <Bar data={chartData} options={{ responsive: true, plugins: { legend: { display: false } } }} />
      </Card>
      <Card title="任务状态">暂无数据</Card>
      <Card title="实时监控">待接入</Card>
    </div>
  )
}
