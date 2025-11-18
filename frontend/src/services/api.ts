import axios from 'axios'
import { getBaseURL } from '../config'

const api = axios.create({ baseURL: getBaseURL() })
api.interceptors.request.use(cfg => { cfg.baseURL = getBaseURL(); return cfg })
export default api
