export type AdminResource = 'classes' | 'summer-shifts' | 'news' | 'services' | 'staff' | 'documents' | 'reviews' | 'gallery' | 'leads'
export type AdminCounts = Record<AdminResource | 'gallery', number>
export type AdminItem = { id: number; title: string; isPublished: boolean; position: number }
export type AdminSession = { token: string; expiresAt: string; user: { id: number; login: string; role: 'administrator' | 'editor' } }
export type MediaItem = { id: number; url: string; mime: string; sizeBytes: number; alt: string }
export type AdminClassForm = {
  id?: number
  title: string
  ageRange: string
  description: string
  schedule: string
  price: string
  emoji: string
  color: string
  coverMediaId: number | null
  photoMediaIds: number[]
  fullDescription: string[]
  highlights: string[]
  position: number
  isPublished: boolean
  coverMedia?: MediaItem
  photoMedia?: MediaItem[]
}
export type AdminNewsForm = {
  id?: number
  title: string
  date: string
  category: string
  excerpt: string
  emoji: string
  coverMediaId: number | null
  photoMediaIds: number[]
  content: string[]
  isPublished: boolean
  coverMedia?: MediaItem
  photoMedia?: MediaItem[]
}
export type AdminStaffForm = { id?:number; name:string; role:string; bio:string; initials:string; color:string; photoMediaId:number|null; position:number; isPublished:boolean; photoMedia?:MediaItem }
export type AdminDocumentForm = { id?:number; title:string; description:string; fileMediaId:number; displayFileType:string; displayFileSize:string; position:number; isPublished:boolean; fileMedia?:MediaItem }
export type AdminReviewForm = { id?:number; author:string; role:string; text:string; rating:number; position:number; isPublished:boolean }
export type AdminGalleryItem = { id:number; mediaId:number; title:string; position:number; media:MediaItem }
export type AdminLead = { id:number; name:string; phone:string; message:string; isProcessed:boolean; createdAt:string }

function headers(token: string, json = false): HeadersInit {
  return { Authorization: `Bearer ${token}`, ...(json ? { 'Content-Type': 'application/json' } : {}) }
}
async function request<T>(path: string, token: string, init: RequestInit = {}): Promise<T> {
  const response = await fetch(`/api/admin${path}`, { ...init, headers: { ...headers(token), ...init.headers } })
  if (!response.ok) {
    if (response.status === 401) throw new Error('Сессия истекла — войдите снова')
    if (response.status === 413) throw new Error('Файл слишком большой — максимальный размер 50 МБ')
    const payload = await response.json().catch(() => null) as { error?: { message?: string } } | null
    throw new Error(payload?.error?.message || 'Не удалось выполнить запрос')
  }
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
export async function uploadMedia(token: string, file: File, alt = ''): Promise<MediaItem> {
  const form = new FormData(); form.set('file', file); form.set('alt', alt)
  return request<MediaItem>('/media', token, { method: 'POST', body: form })
}
export const deleteMedia = (token: string, id: number) => request<void>(`/media/${id}`, token, { method: 'DELETE' })
export const getAdminClassForm = (token: string, id: number) => request<AdminClassForm>(`/classes/${id}/form`, token)
export const saveAdminClassForm = (token: string, value: AdminClassForm) => request<{ id: number }>(value.id ? `/classes/${value.id}/form` : '/classes', token, {
  method: value.id ? 'PUT' : 'POST', headers: headers(token, true), body: JSON.stringify(value),
})
export const deleteAdminClassMedia = (token: string, classId: number, mediaId: number) => request<void>(`/classes/${classId}/media/${mediaId}`, token, { method: 'DELETE' })
export const getAdminSummerShiftForm = (token: string, id: number) => request<AdminClassForm>(`/summer-shifts/${id}/form`, token)
export const saveAdminSummerShiftForm = (token: string, value: AdminClassForm) => request<{ id: number }>(value.id ? `/summer-shifts/${value.id}/form` : '/summer-shifts', token, {
  method: value.id ? 'PUT' : 'POST', headers: headers(token, true), body: JSON.stringify(value),
})
export const deleteAdminSummerShiftMedia = (token: string, shiftId: number, mediaId: number) => request<void>(`/summer-shifts/${shiftId}/media/${mediaId}`, token, { method: 'DELETE' })
export const getAdminNewsForm = (token: string, id: number) => request<AdminNewsForm>(`/news/${id}/form`, token)
export const getAdminNewsCategories = (token: string) => request<string[]>('/news/categories', token)
export const saveAdminNewsForm = (token: string, value: AdminNewsForm) => request<{ id: number }>(value.id ? `/news/${value.id}/form` : '/news', token, {
  method: value.id ? 'PUT' : 'POST', headers: headers(token, true), body: JSON.stringify(value),
})
export const deleteAdminNewsMedia = (token: string, newsId: number, mediaId: number) => request<void>(`/news/${newsId}/media/${mediaId}`, token, { method: 'DELETE' })
export const getAdminStaffForm = (token:string,id:number)=>request<AdminStaffForm>(`/staff/${id}/form`,token)
export const saveAdminStaffForm = (token:string,value:AdminStaffForm)=>request<{id:number}>(value.id?`/staff/${value.id}/form`:'/staff',token,{method:value.id?'PUT':'POST',headers:headers(token,true),body:JSON.stringify(value)})
export const deleteAdminStaffMedia = (token:string,staffId:number,mediaId:number)=>request<void>(`/staff/${staffId}/media/${mediaId}`,token,{method:'DELETE'})
export const getAdminDocumentForm = (token:string,id:number)=>request<AdminDocumentForm>(`/documents/${id}/form`,token)
export const saveAdminDocumentForm = (token:string,value:AdminDocumentForm)=>request<{id:number}>(value.id?`/documents/${value.id}/form`:'/documents',token,{method:value.id?'PUT':'POST',headers:headers(token,true),body:JSON.stringify(value)})
export const getAdminReviewForm=(token:string,id:number)=>request<AdminReviewForm>(`/reviews/${id}/form`,token)
export const saveAdminReviewForm=(token:string,value:AdminReviewForm)=>request<{id:number}>(value.id?`/reviews/${value.id}/form`:'/reviews',token,{method:value.id?'PUT':'POST',headers:headers(token,true),body:JSON.stringify(value)})
export const getAdminGallery=(token:string)=>request<AdminGalleryItem[]>('/gallery',token)
export const addAdminGallery=(token:string,mediaId:number,title='')=>request<{id:number}>('/gallery',token,{method:'POST',headers:headers(token,true),body:JSON.stringify({mediaId,title})})
export const deleteAdminGallery=(token:string,id:number)=>request<void>(`/gallery/${id}`,token,{method:'DELETE'})
export const getAdminLeads=(token:string)=>request<AdminLead[]>('/leads',token)
export const setAdminLeadProcessed=(token:string,id:number,isProcessed:boolean)=>request<void>(`/leads/${id}`,token,{method:'PATCH',headers:headers(token,true),body:JSON.stringify({isProcessed})})
export const deleteAdminLead=(token:string,id:number)=>request<void>(`/leads/${id}`,token,{method:'DELETE'})
