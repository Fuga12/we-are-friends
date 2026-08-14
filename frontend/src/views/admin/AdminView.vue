<script setup lang="ts">
import { computed, onMounted, ref } from 'vue'
import Button from 'primevue/button'
import Column from 'primevue/column'
import DataTable from 'primevue/datatable'
import InputText from 'primevue/inputtext'
import SelectButton from 'primevue/selectbutton'
import ToggleSwitch from 'primevue/toggleswitch'
import DeleteButton from '@/components/admin/DeleteButton.vue'
import { deleteAdminItem, getAdminContent, getAdminItems, getAdminOverview, setPublication, signIn, signOut, updateAdminContent, updateAdminItem, type AdminCounts, type AdminItem, type AdminResource } from '@/api/admin'
import ClassForm from './ClassForm.vue'
import NewsForm from './NewsForm.vue'
import StaffForm from './StaffForm.vue'
import DocumentForm from './DocumentForm.vue'
import ReviewForm from './ReviewForm.vue'
import GalleryManager from './GalleryManager.vue'
import LeadManager from './LeadManager.vue'

type Item = AdminItem

const sections: { label: string; value: AdminResource }[] = [
  { label: 'Занятия', value: 'classes' },
  { label: 'Летние смены', value: 'summer-shifts' },
  { label: 'Новости', value: 'news' },
  { label: 'Сотрудники', value: 'staff' },
  { label: 'Документы', value: 'documents' },
  { label: 'Отзывы', value: 'reviews' },
  { label: 'Галерея', value: 'gallery' },
  { label: 'Заявки', value: 'leads' },
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
const editingId = ref<number | null>(null)
const contentDraft = ref('')
const classFormId = ref<number | null | undefined>(undefined)
const summerFormId = ref<number | null | undefined>(undefined)
const newsFormId = ref<number | null | undefined>(undefined)
const staffFormId = ref<number | null | undefined>(undefined)
const documentFormId = ref<number | null | undefined>(undefined)
const reviewFormId = ref<number | null | undefined>(undefined)

const filtered = computed(() => {
  const value = query.value.trim().toLowerCase()
  if (!value) return items.value
  return items.value.filter((item) => item.title.toLowerCase().includes(value))
})

async function loadItems() {
  if(['gallery','leads'].includes(active.value)){items.value=[];return}
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
  classFormId.value = undefined
  summerFormId.value = undefined
  newsFormId.value = undefined
  staffFormId.value = undefined
  documentFormId.value = undefined
  reviewFormId.value = undefined
  editingId.value = null
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
async function reorderClasses(event: { value: Item[] }) {
  if (['news','gallery'].includes(active.value) || query.value.trim()) return
  const previous = items.value.map(item => ({ ...item }))
  const reordered = event.value.map((item, position) => ({ ...item, position }))
  items.value = reordered
  loading.value = true
  error.value = ''
  try {
    await Promise.all(reordered.map(item => updateAdminItem(active.value, item, token.value)))
    await loadItems()
  } catch (reason) {
    items.value = previous
    error.value = reason instanceof Error ? reason.message : 'Не удалось сохранить новый порядок'
  } finally {
    loading.value = false
  }
}
async function editContent(item: Item) { if (active.value === 'classes') { classFormId.value = item.id; return } if (active.value === 'summer-shifts') { summerFormId.value = item.id; return } if (active.value === 'news') { newsFormId.value = item.id; return } if(active.value==='staff'){staffFormId.value=item.id;return} if(active.value==='documents'){documentFormId.value=item.id;return} if(active.value==='reviews'){reviewFormId.value=item.id;return} loading.value = true; error.value = ''; try { contentDraft.value = JSON.stringify(await getAdminContent(active.value, item.id, token.value), null, 2); editingId.value = item.id } catch (reason) { error.value = reason instanceof Error ? reason.message : 'Не удалось открыть редактор' } finally { loading.value = false } }
async function classSaved() { classFormId.value = undefined; await loadItems(); counts.value = await getAdminOverview(token.value) }
async function summerSaved() { summerFormId.value = undefined; await loadItems(); counts.value = await getAdminOverview(token.value) }
async function staffSaved(){staffFormId.value=undefined;await loadItems();counts.value=await getAdminOverview(token.value)}
async function documentSaved(){documentFormId.value=undefined;await loadItems();counts.value=await getAdminOverview(token.value)}
async function reviewSaved(){reviewFormId.value=undefined;await loadItems();counts.value=await getAdminOverview(token.value)}
async function galleryChanged(){counts.value=await getAdminOverview(token.value)}
async function newsSaved() { newsFormId.value = undefined; await loadItems(); counts.value = await getAdminOverview(token.value) }
async function saveContent() { if (editingId.value === null) return; loading.value = true; error.value = ''; try { const payload = JSON.parse(contentDraft.value) as Record<string, unknown>; await updateAdminContent(active.value, editingId.value, payload, token.value); editingId.value = null; await loadItems() } catch (reason) { error.value = reason instanceof Error ? (reason instanceof SyntaxError ? 'Проверьте синтаксис JSON' : reason.message) : 'Не удалось сохранить материал' } finally { loading.value = false } }
async function remove(item: Item) { loading.value = true; try { await deleteAdminItem(active.value, item.id, token.value); await loadItems(); counts.value = await getAdminOverview(token.value) } catch (reason) { error.value = reason instanceof Error ? reason.message : 'Не удалось удалить запись' } finally { loading.value = false } }

async function leave() { try { await signOut(token.value) } finally { sessionStorage.removeItem('admin-token'); token.value = ''; items.value = [] } }

onMounted(() => { if (token.value) { void loadItems(); void getAdminOverview(token.value).then(v => counts.value = v).catch(() => leave()) } })
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
          <p>{{ active==='leads' ? 'Здесь находятся заявки, отправленные посетителями через форму записи.' : ['classes','summer-shifts'].includes(active) ? 'Основные данные меняются через форму. Видимость можно быстро переключить в списке.' : active === 'news' ? 'Содержание новости и фотографии меняются через форму. Видимость можно переключить в списке.' : 'Редактируйте названия, порядок и видимость прямо в списке. Удаление требует отдельного подтверждения.' }}</p>
        </section>
        <section v-if="counts" class="stats">
          <article v-for="(label, key) in { classes: 'Занятия', 'summer-shifts': 'Летние смены', news: 'Новости', staff: 'Сотрудники', documents: 'Документы', reviews:'Отзывы', gallery: 'Фото', leads:'Новые заявки' }" :key="key">
            <strong>{{ counts[key as keyof AdminCounts] }}</strong><span>{{ label }}</span>
          </article>
        </section>
        <section class="content-card admin-navigation">
          <div class="content-toolbar">
            <SelectButton v-model="active" :options="sections" option-label="label" option-value="value" @change="changeSection" />
            <div v-if="!['gallery','leads'].includes(active)" class="toolbar-actions"><InputText v-model="query" placeholder="Найти запись" /><Button v-if="active === 'classes'" label="Новое занятие" @click="classFormId = null" /><Button v-if="active === 'summer-shifts'" label="Новая смена" @click="summerFormId = null" /><Button v-if="active === 'news'" label="Новая новость" @click="newsFormId = null" /><Button v-if="active==='staff'" label="Новый сотрудник" @click="staffFormId=null"/><Button v-if="active==='documents'" label="Новый документ" @click="documentFormId=null"/><Button v-if="active==='reviews'" label="Новый отзыв" @click="reviewFormId=null"/></div>
          </div>
        </section>
        <section v-if="!['gallery','leads'].includes(active)" class="content-card content-list">
          <p v-if="active!=='news'" class="order-hint">Записи отображаются в порядке, согласно этому списку.<span v-if="query.trim()"> Очистите поиск, чтобы включить перетаскивание.</span></p>
          <small v-if="error" class="form-error">{{ error }}</small>
          <DataTable :value="filtered" :loading="loading" striped-rows data-key="id" empty-message="Записей нет" @row-reorder="reorderClasses">
            <Column v-if="active!=='news'" :row-reorder="!query.trim()" header="" style="width: 3rem" />
            <Column field="id" header="ID" style="width: 5rem" />
            <Column header="Название"><template #body="{ data }"><strong class="item-title">{{ data.title }}</strong></template></Column>
            <Column header="Показывать" style="width: 10rem"><template #body="{ data }"><ToggleSwitch :model-value="data.isPublished" @update:model-value="toggle(data, $event)" /></template></Column>
            <Column header="Действия" style="width: 17rem"><template #body="{ data }"><Button label="Изменить" size="small" text @click="editContent(data)" /><DeleteButton size="small" @confirm="remove(data)" /></template></Column>
          </DataTable>
        </section>
        <section v-if="active==='gallery'" class="content-card"><GalleryManager :token="token" @changed="galleryChanged"/></section>
        <section v-if="active==='leads'" class="content-card"><LeadManager :token="token" @changed="galleryChanged"/></section>
        <section v-if="classFormId !== undefined" class="content-card editor-card">
          <ClassForm :key="classFormId ?? 'new'" :token="token" :id="classFormId" @saved="classSaved" @cancel="classFormId = undefined" />
        </section>
        <section v-if="summerFormId !== undefined" class="content-card editor-card">
          <ClassForm :key="summerFormId ?? 'new-summer'" :token="token" :id="summerFormId" kind="summer" @saved="summerSaved" @cancel="summerFormId = undefined" />
        </section>
        <section v-if="newsFormId !== undefined" class="content-card editor-card">
          <NewsForm :key="newsFormId ?? 'new'" :token="token" :id="newsFormId" @saved="newsSaved" @cancel="newsFormId = undefined" />
        </section>
        <section v-if="staffFormId!==undefined" class="content-card editor-card"><StaffForm :key="staffFormId??'new-staff'" :token="token" :id="staffFormId" @saved="staffSaved" @cancel="staffFormId=undefined"/></section>
        <section v-if="documentFormId!==undefined" class="content-card editor-card"><DocumentForm :key="documentFormId??'new-document'" :token="token" :id="documentFormId" @saved="documentSaved" @cancel="documentFormId=undefined"/></section>
        <section v-if="reviewFormId!==undefined" class="content-card editor-card"><ReviewForm :key="reviewFormId??'new-review'" :token="token" :id="reviewFormId" @saved="reviewSaved" @cancel="reviewFormId=undefined"/></section>
        <section v-if="editingId !== null" class="content-card editor-card">
          <div class="content-toolbar"><div><p class="eyebrow">Редактор материала</p><h2>Все поля записи</h2></div><Button label="Закрыть" text @click="editingId = null" /></div>
          <p class="media-hint">Можно менять текст, даты, статус публикации, порядок и привязанные <code>media_id</code>. Массивы записываются как JSON-массивы.</p>
          <textarea v-model="contentDraft" class="json-editor" spellcheck="false" aria-label="JSON редактор материала" />
          <div class="editor-actions"><Button label="Сохранить все поля" :loading="loading" @click="saveContent" /><Button label="Отмена" severity="secondary" text @click="editingId = null" /></div>
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
label { display:grid; gap:.45rem; font-weight:700; }.eyebrow { color:#e46129; font-weight:800; text-transform:uppercase; letter-spacing:.08em; font-size:.75rem; margin:0 0 .35rem; }
:deep(.form-error) { display:block; margin:.75rem 0; padding:.85rem 1rem; border:2px solid #b42318; border-radius:10px; background:#fff1f0; color:#b42318; font-size:1.1rem; font-weight:800; line-height:1.4; }
.admin-intro { display:flex; justify-content:space-between; align-items:end; gap:2rem; margin-bottom:1.8rem; }.admin-intro > p { max-width:430px; margin:0; }
.stats { display:grid; grid-template-columns:repeat(auto-fit,minmax(110px,1fr)); gap:.75rem; margin-bottom:1.8rem; }.stats article { background:#e9f3ef; border-radius:14px; padding:1rem; display:grid; gap:.2rem; }.stats strong{font-size:1.6rem}.stats span{font-size:.8rem;color:#527168}
.content-toolbar { display:flex; justify-content:space-between; gap:1rem; margin-bottom:1.25rem; }.content-card small { display:block; margin-top:.2rem; } :deep(.p-datatable td small) { display:block; color:#6f766e; margin-top:.25rem; }
.admin-navigation{margin-bottom:1rem}.admin-navigation .content-toolbar{margin-bottom:0;flex-direction:column;align-items:stretch}.admin-navigation :deep(.p-selectbutton){display:grid;grid-template-columns:repeat(8,minmax(0,1fr));width:100%}.admin-navigation :deep(.p-selectbutton .p-button){min-width:0;padding:.7rem .45rem;white-space:normal;line-height:1.2}.content-list{padding-top:1.5rem}
.toolbar-actions { display:flex; width:100%; justify-content:flex-end; flex-wrap:wrap; gap:.75rem; align-items:center }.toolbar-actions :deep(.p-inputtext) { min-width:13rem }
.order-hint { margin:0 0 1.25rem; padding:.85rem 1rem; border-radius:10px; background:#f1f6f3; color:#527168; line-height:1.45 }.item-title { font-weight:700 }
.order-input { width:100%; padding:.55rem; border:1px solid #d9d4ca; border-radius:6px; color:#173d35 }
.editor-card { margin-top:1.8rem }.editor-card h2 { margin:0 }.json-editor { width:100%; min-height:28rem; font:13px/1.5 ui-monospace,SFMono-Regular,Menlo,monospace; padding:1rem; border:1px solid #d9d4ca; border-radius:10px; resize:vertical; box-sizing:border-box; color:#173d35 }.editor-actions { display:flex; gap:.5rem; margin-top:1rem }
:deep(.p-inputtext) { width:100%; background:#fff; color:#173d35; border-color:#d9d4ca; }
:deep(.p-inputtext:hover), :deep(.p-inputtext:enabled:focus) { border-color:#1a6757; }
@media(max-width:850px){.admin-navigation{overflow:hidden}.admin-navigation :deep(.p-selectbutton){display:flex;overflow-x:auto;justify-content:flex-start}.admin-navigation :deep(.p-selectbutton .p-button){flex:0 0 auto;padding:.7rem 1rem}.toolbar-actions{justify-content:stretch}.toolbar-actions>*{flex:1 1 12rem}}
:deep(.p-inputtext::placeholder) { color:#7b857e; }
@media(max-width:760px){.admin-intro,.content-toolbar,.toolbar-actions{align-items:stretch;flex-direction:column}.stats{grid-template-columns:repeat(2,1fr)}}
</style>
