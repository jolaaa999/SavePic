const BASE = ''

async function request(url, options = {}) {
  const res = await fetch(`${BASE}${url}`, options)
  const json = await res.json().catch(() => ({}))
  if (json.code !== 200) {
    throw new Error(json.msg || `请求失败 (${res.status})`)
  }
  return json.data
}

export function fetchCategories() {
  return request('/api/categories')
}

export function createCategory(payload) {
  return request('/api/categories', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })
}

export function deleteCategory(id) {
  return request(`/api/categories/${id}`, { method: 'DELETE' })
}

export function updateCategory(id, payload) {
  return request(`/api/categories/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })
}

export function fetchTags() {
  return request('/api/tags')
}

export function updateTag(id, payload) {
  return request(`/api/tags/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })
}

/** @param {{ categoryId?: number, tagIds?: number[], sort?: string }} params */
export function fetchMemes(params = {}) {
  const q = new URLSearchParams()
  if (params.categoryId) q.set('category_id', String(params.categoryId))
  if (params.tagIds?.length) q.set('tag_ids', params.tagIds.join(','))
  if (params.sort) q.set('sort', params.sort)
  const qs = q.toString()
  return request(`/api/memes${qs ? `?${qs}` : ''}`)
}

export function uploadMeme(categoryId, file, tags = []) {
  const form = new FormData()
  form.append('file', file)
  form.append('category_id', String(categoryId))
  if (tags.length) {
    form.append('tags', JSON.stringify(tags))
  }
  return request('/api/memes/upload', { method: 'POST', body: form })
}

export function updateMeme(id, payload) {
  return request(`/api/memes/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(payload),
  })
}

export function deleteMeme(id) {
  return request(`/api/memes/${id}`, { method: 'DELETE' })
}

export function resolveFileUrl(path) {
  if (!path) return ''
  if (path.startsWith('http')) return path
  return `${window.location.origin}${path}`
}

/** 从上传响应中取出 meme 对象 */
export function unwrapUploadResult(data) {
  return data?.meme ?? data
}
