<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import Button from 'primevue/button'
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import InputText from 'primevue/inputtext'
import SelectButton from 'primevue/selectbutton'
import ToggleSwitch from 'primevue/toggleswitch'
import { deleteAdminItem, deleteMedia, getAdminContent, getAdminItems, getAdminOverview, getMedia, setPublication, signIn, signOut, updateAdminContent, updateAdminItem, uploadMedia, type AdminCounts, type AdminItem, type AdminResource, type MediaItem } from '@/api/admin'

type Item = AdminItem

const sections: { label: string; value: AdminResource }[] = [
  { label: 'Занятия', value: 'classes' },
  { label: 'Новости', value: 'news' },
  { label: 'Услуги', value: 'services' },
  { label: 'Сотрудники', value: 'staff' },
  { label: 'Документы', value: 'documents' },
]
const user = ref('')
const password = ref('')
const token = ref(sessionStorage.getItem('admin-token') ?? '')
const authorized = computed(() => Boolean(token.value))
const loading = ref(false)
const error = ref('')
const active = ref<AdminResource>('classes')
const counts = ref<AdminCounts | null>(null)
const items = ref<Item[]>([])
const query = ref('')
const media = ref<MediaItem[]>([])
const mediaLoading = ref(false)
const fileInput = ref<HTMLInputElement | null>(null)
const editingId = ref<number | null>(null)
const contentDraft = ref('')

const filtered = computed(() => {
  const value = query.value.trim().toLowerCase()
  if (!value) return items.value
  return items.value.filter((item) => item.title.toLowerCase().includes(value))
})

async function loadItems() {
  items.value = await getAdminItems(active.value, token.value)
}

async function authenticate() {
  loading.value = true
  error.value = ''
  try {
    const session = await signIn(user.value, password.value)
    token.value = session.token
    sessionStorage.setItem('admin-token', session.token)
    counts.value = await getAdminOverview(token.value)
    await loadItems()
  } catch (reason) {
    error.value = reason instanceof Error ? reason.message : 'Не удалось войти'
  } finally {
    loading.value = false
  }
}

async function changeSection() {
  query.value = ''
  loading.value = true
  try { await loadItems() } catch (reason) { error.value = reason instanceof Error ? reason.message : 'Ошибка загрузки' } finally { loading.value = false }
}

async function toggle(item: Item, value: boolean) {
  loading.value = true
  error.value = ''
  try {
    await setPublication(active.value, item.id, value, token.value)
    await loadItems()
    counts.value = await getAdminOverview(token.value)
  } catch (reason) {
    error.value = reason instanceof Error ? reason.message : 'Ошибка сохранения'
  } finally {
    loading.value = false
  }
}
async function save(item: Item) { loading.value = true; error.value = ''; try { await updateAdminItem(active.value, item, token.value); await loadItems() } catch (reason) { error.value = reason instanceof Error ? reason.message : 'Не удалось сохранить запись' } finally { loading.value = false } }
async function editContent(item: Item) { loading.value = true; error.value = ''; try { contentDraft.value = JSON.stringify(await getAdminContent(active.value, item.id, token.value), null, 2); editingId.value = item.id } catch (reason) { error.value = reason instanceof Error ? reason.message : 'Не удалось открыть редактор' } finally { loading.value = false } }
async function saveContent() { if (editingId.value === null) return; loading.value = true; error.value = ''; try { const payload = JSON.parse(contentDraft.value) as Record<string, unknown>; await updateAdminContent(active.value, editingId.value, payload, token.value); editingId.value = null; await loadItems() } catch (reason) { error.value = reason instanceof Error ? (reason instanceof SyntaxError ? 'Проверьте синтаксис JSON' : reason.message) : 'Не удалось сохранить материал' } finally { loading.value = false } }
async function remove(item: Item) { if (!confirm(`Удалить «${item.title}» без возможности восстановления?`)) return; loading.value = true; try { await deleteAdminItem(active.value, item.id, token.value); await loadItems(); counts.value = await getAdminOverview(token.value) } catch (reason) { error.value = reason instanceof Error ? reason.message : 'Не удалось удалить запись' } finally { loading.value = false } }

async function loadMedia() { mediaLoading.value = true; try { media.value = (await getMedia(token.value, query.value)).items } finally { mediaLoading.value = false } }
async function onFile(event: Event) { const file = (event.target as HTMLInputElement).files?.[0]; if (!file) return; mediaLoading.value = true; try { await uploadMedia(token.value, file); await loadMedia() } catch (reason) { error.value = reason instanceof Error ? reason.message : 'Не удалось загрузить файл' } finally { mediaLoading.value = false; if (fileInput.value) fileInput.value.value = '' } }
async function removeMedia(id: number) { if (!confirm('Удалить файл без возможности восстановления?')) return; try { await deleteMedia(token.value, id); await loadMedia() } catch (reason) { error.value = reason instanceof Error ? reason.message : 'Не удалось удалить файл' } }
async function leave() { try { await signOut(token.value) } finally { sessionStorage.removeItem('admin-token'); token.value = ''; items.value = []; media.value = [] } }

onMounted(() => { if (token.value) { void loadItems(); void getAdminOverview(token.value).then(v => counts.value = v).catch(() => leave()); void loadMedia() } })
</script>

<template>
  <div class="admin-shell">
    <header class="admin-header">
      <RouterLink to="/" class="admin-brand">Мы Друзья <span>· админка</span></RouterLink>
      <Button v-if="authorized" label="Выйти" text severity="secondary" @click="leave" />
    </header>

    <main class="admin-main">
      <section v-if="!authorized" class="login-card">
        <p class="eyebrow">Управление сайтом</p>
        <h1>Вход в админку</h1>
        <p>Используйте логин и пароль из переменных <code>ADMIN_USER</code> и <code>ADMIN_PASSWORD</code>.</p>
        <label>Логин<InputText v-model="user" autocomplete="username" /></label>
        <label>Пароль<InputText v-model="password" type="password" autocomplete="current-password" @keyup.enter="authenticate" /></label>
        <small v-if="error" class="form-error">{{ error }}</small>
        <Button label="Войти" :loading="loading" @click="authenticate" />
      </section>

      <template v-else>
        <section class="admin-intro">
          <div><p class="eyebrow">Контент сайта</p><h1>Публикации</h1></div>
          <p>Редактируйте названия, порядок и видимость прямо в списке. Удаление требует отдельного подтверждения.</p>
        </section>
        <section v-if="counts" class="stats">
          <article v-for="(label, key) in { classes: 'Занятия', news: 'Новости', services: 'Услуги', staff: 'Сотрудники', documents: 'Документы', gallery: 'Фото' }" :key="key">
            <strong>{{ counts[key as keyof AdminCounts] }}</strong><span>{{ label }}</span>
          </article>
        </section>
        <section class="content-card">
          <div class="content-toolbar">
            <SelectButton v-model="active" :options="sections" option-label="label" option-value="value" @change="changeSection" />
            <InputText v-model="query" placeholder="Найти запись" />
          </div>
          <small v-if="error" class="form-error">{{ error }}</small>
          <DataTable :value="filtered" :loading="loading" striped-rows data-key="id" empty-message="Записей нет">
            <Column field="id" header="ID" style="width: 5rem" />
            <Column header="Название"><template #body="{ data }"><InputText v-model="data.title" :aria-label="`Название ${data.id}`" /></template></Column>
            <Column header="Порядок" style="width: 7rem"><template #body="{ data }"><input v-model.number="data.position" class="order-input" type="number" min="0" /></template></Column>
            <Column header="Показывать" style="width: 10rem"><template #body="{ data }"><ToggleSwitch :model-value="data.isPublished" @update:model-value="toggle(data, $event)" /></template></Column>
            <Column header="Действия" style="width: 17rem"><template #body="{ data }"><Button label="Изменить" size="small" text @click="editContent(data)" /><Button label="Сохранить" size="small" text @click="save(data)" /><Button label="Удалить" size="small" text severity="danger" @click="remove(data)" /></template></Column>
          </DataTable>
        </section>
        <section v-if="editingId !== null" class="content-card editor-card">
          <div class="content-toolbar"><div><p class="eyebrow">Редактор материала</p><h2>Все поля записи</h2></div><Button label="Закрыть" text @click="editingId = null" /></div>
          <p class="media-hint">Можно менять текст, даты, статус публикации, порядок и привязанные <code>media_id</code>. Массивы записываются как JSON-массивы.</p>
          <textarea v-model="contentDraft" class="json-editor" spellcheck="false" aria-label="JSON редактор материала" />
          <div class="editor-actions"><Button label="Сохранить все поля" :loading="loading" @click="saveContent" /><Button label="Отмена" severity="secondary" text @click="editingId = null" /></div>
        </section>
        <section class="content-card media-card">
          <div class="content-toolbar"><div><p class="eyebrow">Медиатека</p><h2>Файлы</h2></div><Button label="Загрузить файл" icon="pi pi-upload" @click="fileInput?.click()" /></div>
          <input ref="fileInput" type="file" hidden accept="image/*,video/*,application/pdf" @change="onFile" />
          <p class="media-hint">Перетащите файл в окно выбора или нажмите кнопку. Неиспользуемые файлы можно удалить.</p>
          <div class="media-grid" :aria-busy="mediaLoading"><article v-for="file in media" :key="file.id"><img v-if="file.mime.startsWith('image/')" :src="file.url" :alt="file.alt" /><div v-else class="media-file">{{ file.mime }}</div><small>{{ file.alt || file.url.split('/').pop() }}</small><Button label="Удалить" severity="danger" text size="small" @click="removeMedia(file.id)" /></article></div>
        </section>
      </template>
    </main>
  </div>
</template>

<style scoped>
.admin-shell { min-height: 100vh; background: #f7f6f2; color: #173d35; }
.admin-header { display:flex; justify-content:space-between; align-items:center; padding:1.1rem clamp(1.25rem, 5vw, 5rem); background:#fff; border-bottom:1px solid #e6e2d9; }
.admin-brand { color:#1a6757; font-weight:800; text-decoration:none; font-size:1.2rem; }.admin-brand span { color:#6f766e; font-weight:500; }
.admin-main { width:min(1100px, calc(100% - 2.5rem)); margin:0 auto; padding:3.5rem 0 5rem; }
.login-card,.content-card { background:#fff; border:1px solid #e6e2d9; border-radius:20px; padding:clamp(1.5rem,4vw,3rem); box-shadow:0 14px 35px #173d3510; }
.login-card { width:min(460px,100%); margin:7vh auto; display:grid; gap:1rem; }.login-card h1,.admin-intro h1 { margin:0; }.login-card p,.admin-intro p { color:#647067; line-height:1.5; }
label { display:grid; gap:.45rem; font-weight:700; }.eyebrow { color:#e46129; font-weight:800; text-transform:uppercase; letter-spacing:.08em; font-size:.75rem; margin:0 0 .35rem; }.form-error { color:#b42318; }
.admin-intro { display:flex; justify-content:space-between; align-items:end; gap:2rem; margin-bottom:1.8rem; }.admin-intro > p { max-width:430px; margin:0; }
.stats { display:grid; grid-template-columns:repeat(6,1fr); gap:.75rem; margin-bottom:1.8rem; }.stats article { background:#e9f3ef; border-radius:14px; padding:1rem; display:grid; gap:.2rem; }.stats strong{font-size:1.6rem}.stats span{font-size:.8rem;color:#527168}
.content-toolbar { display:flex; justify-content:space-between; gap:1rem; margin-bottom:1.25rem; }.content-card small { display:block; margin-top:.2rem; } :deep(.p-datatable td small) { display:block; color:#6f766e; margin-top:.25rem; }
.media-card { margin-top:1.8rem }.media-card h2 { margin:0 }.media-hint { color:#647067 }.media-grid { display:grid; grid-template-columns:repeat(auto-fill,minmax(140px,1fr)); gap:1rem }.media-grid article { border:1px solid #e6e2d9; padding:.65rem; border-radius:12px; overflow:hidden; display:grid; gap:.45rem }.media-grid img,.media-file { width:100%; height:100px; object-fit:cover; background:#e9f3ef; border-radius:8px }.media-file { display:grid; place-items:center; padding:.5rem; text-align:center; font-size:.72rem }
.order-input { width:100%; padding:.55rem; border:1px solid #d9d4ca; border-radius:6px; color:#173d35 }
.editor-card { margin-top:1.8rem }.editor-card h2 { margin:0 }.json-editor { width:100%; min-height:28rem; font:13px/1.5 ui-monospace,SFMono-Regular,Menlo,monospace; padding:1rem; border:1px solid #d9d4ca; border-radius:10px; resize:vertical; box-sizing:border-box; color:#173d35 }.editor-actions { display:flex; gap:.5rem; margin-top:1rem }
:deep(.p-inputtext) { width:100%; background:#fff; color:#173d35; border-color:#d9d4ca; }
:deep(.p-inputtext:hover), :deep(.p-inputtext:enabled:focus) { border-color:#1a6757; }
:deep(.p-inputtext::placeholder) { color:#7b857e; }
@media(max-width:760px){.admin-intro,.content-toolbar{align-items:stretch;flex-direction:column}.stats{grid-template-columns:repeat(2,1fr)}}
</style>
