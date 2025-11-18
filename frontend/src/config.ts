let appConfig: { baseURL: string, ssePath: string, chartRefreshMs: number } = { baseURL: 'http://localhost:8080', ssePath: '/sse/tasks', chartRefreshMs: 5000 }
fetch('/config.json').then(r=>r.json()).then(c=>{ appConfig = { ...appConfig, ...c } }).catch(()=>{})
export function getBaseURL(){ return appConfig.baseURL }
export function getSSEPath(){ return appConfig.ssePath }
export function getChartRefreshMs(){ return appConfig.chartRefreshMs }
