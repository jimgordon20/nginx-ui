import { http } from '@uozi-admin/request'

export interface ShieldConfig {
  site: string
  applied: boolean
  hotlink: {
    enabled: boolean
    allow_empty: boolean
    allow_referers: string[]
    extensions: string[]
  }
  rate_limit: {
    enabled: boolean
    rate: string
    burst: number
    path_regex: string
  }
  waf: {
    enabled: boolean
    block_empty_ua: boolean
    block_bad_bots: boolean
    blocked_user_agents: string[]
    blocked_paths: string[]
  }
  under_attack: {
    enabled: boolean
    rate: string
    burst: number
  }
  ip_block: {
    enabled: boolean
    deny_ips: string[]
  }
  security_headers: {
    enabled: boolean
  }
  cache: {
    enabled: boolean
    path_regex: string
    expires: string
    cache_control: string
  }
  preview?: string
}

export interface AlertRule {
  id?: string
  name: string
  enabled: boolean
  site: string
  log_path: string
  metric: string
  threshold: number
  window_minutes: number
  notify_on_trigger: boolean
}

export interface AlertEvaluation {
  rule: AlertRule
  value: number
  triggered: boolean
  message: string
}

export interface HealthItem {
  name: string
  url: string
  status: string
  status_code: number
  latency_ms: number
  error?: string
}

const shieldApi = {
  get(site: string): Promise<ShieldConfig> {
    return http.get(`/shield/${encodeURIComponent(site)}`)
  },

  save(site: string, config: ShieldConfig): Promise<ShieldConfig> {
    return http.post(`/shield/${encodeURIComponent(site)}`, config)
  },

  apply(site: string, config: ShieldConfig): Promise<ShieldConfig> {
    return http.post(`/shield/${encodeURIComponent(site)}/apply`, config)
  },

  detach(site: string): Promise<ShieldConfig> {
    return http.post(`/shield/${encodeURIComponent(site)}/detach`)
  },

  applyImagePreset(site: string): Promise<ShieldConfig> {
    return http.post(`/shield/${encodeURIComponent(site)}/preset/image`)
  },

  health(): Promise<{ data: HealthItem[] }> {
    return http.get('/shield/health')
  },

  listAlerts(): Promise<{ data: AlertRule[] }> {
    return http.get('/shield/alerts')
  },

  saveAlert(rule: AlertRule): Promise<{ data: AlertRule[] }> {
    return http.post('/shield/alerts', rule)
  },

  deleteAlert(id: string): Promise<{ data: AlertRule[] }> {
    return http.post('/shield/alerts/delete', { id })
  },

  evaluateAlerts(): Promise<{ data: AlertEvaluation[] }> {
    return http.post('/shield/alerts/evaluate')
  },
}

export default shieldApi
