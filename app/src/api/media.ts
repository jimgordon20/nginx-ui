import { http } from '@uozi-admin/request'

export interface MediaRoot {
  name: string
  path: string
  size: number
  count: number
}

export interface MediaEntry {
  name: string
  path: string
  is_dir: boolean
  size: number
  mod_time: number
  mime_type?: string
  is_image: boolean
  extension?: string
}

const mediaApi = {
  listRoots(): Promise<{ data: MediaRoot[] }> {
    return http.get('/media/roots')
  },

  list(root: string, path = ''): Promise<{ root: string, path: string, entries: MediaEntry[] }> {
    return http.get('/media/list', { params: { root, path } })
  },

  mkdir(root: string, path: string, name: string): Promise<{ message: string }> {
    return http.post('/media/mkdir', { root, path, name })
  },

  delete(root: string, paths: string[]): Promise<{ message: string, count: number }> {
    return http.post('/media/delete', { root, paths })
  },

  upload(root: string, path: string, files: File[]): Promise<{ message: string, uploaded: number }> {
    const form = new FormData()
    form.append('root', root)
    form.append('path', path)
    files.forEach(file => form.append('files', file))
    return http.post('/media/upload', form, {
      headers: { 'Content-Type': 'multipart/form-data' },
    })
  },

  fileUrl(root: string, path: string): string {
    const params = new URLSearchParams({ root, path })
    return `/api/media/file?${params.toString()}`
  },

  async fetchPreviewBlob(root: string, path: string): Promise<Blob> {
    return http.get('/media/file', {
      params: { root, path },
      responseType: 'blob',
    })
  },
}

export default mediaApi
