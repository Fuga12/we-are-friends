export type AdminResource = 'classes' | 'news' | 'services' | 'staff' | 'documents'
export type AdminCounts = Record<AdminResource | 'gallery', number>
export type AdminItem = { id: number; title: string; isPublished: boolean; position: number }
export type AdminSession = { token: string; expiresAt: string; user: { id: number; login: string; role: 'administrator' | 'editor' } }
export type MediaItem = { id: number; url: string; mime: string; sizeBytes: number; alt: string }
export type MediaPage = { items: MediaItem[]; page: number; perPage: number; total: number }

function headers(token: string, json = false): HeadersInit {
  return { Authorization: `Bearer ${token}`, ...(json ? { 'Content-Type': 'application/json' } : {}) }
}
async function request<T>(path: string, token: string, init: RequestInit = {}): Promise<T> {
  const response = await fetch(`/api/admin${path}`, { ...init, headers: { ...headers(token), ...init.headers } })
  if (!response.ok) throw new Error(response.status === 401 ? 'Сессия истекла — войдите снова' : 'Не удалось выполнить запрос')
  return response.status === 204 ? undefined as T : response.json() as Promise<T>
}
export async function signIn(login: string, password: string): Promise<AdminSession> {
  const response = await fetch('/api/admin/login', { method: 'POST', headers: { 'Content-Type': 'application/json' }, body: JSON.stringify({ login, password }) })
  if (!response.ok) throw new Error('Неверный логин или пароль')
  return response.json() as Promise<AdminSession>
}
export const signOut = (token: string) => request<void>('/logout', token, { method: 'POST' })
export const getAdminOverview = (token: string) => request<AdminCounts>('/overview', token)
export const getAdminItems = (resource: AdminResource, token: string) => request<AdminItem[]>(`/${resource}`, token)
export const setPublication = (resource: AdminResource, id: number, isPublished: boolean, token: string) => request<void>(`/${resource}/${id}`, token, { method: 'PATCH', headers: headers(token, true), body: JSON.stringify({ isPublished }) })
export const updateAdminItem = (resource: AdminResource, item: AdminItem, token: string) => request<void>(`/${resource}/${item.id}`, token, { method: 'PUT', headers: headers(token, true), body: JSON.stringify(item) })
export const deleteAdminItem = (resource: AdminResource, id: number, token: string) => request<void>(`/${resource}/${id}`, token, { method: 'DELETE' })
export const getAdminContent = (resource: AdminResource, id: number, token: string) => request<Record<string, unknown>>(`/${resource}/${id}/content`, token)
export const updateAdminContent = (resource: AdminResource, id: number, content: Record<string, unknown>, token: string) => request<void>(`/${resource}/${id}/content`, token, { method: 'PATCH', headers: headers(token, true), body: JSON.stringify(content) })
export const getMedia = (token: string, q = '') => request<MediaPage>(`/media?${new URLSearchParams({ q })}`, token)
export async function uploadMedia(token: string, file: File, alt = ''): Promise<MediaItem> {
  const form = new FormData(); form.set('file', file); form.set('alt', alt)
  return request<MediaItem>('/media', token, { method: 'POST', body: form })
}
export const deleteMedia = (token: string, id: number) => request<void>(`/media/${id}`, token, { method: 'DELETE' })
