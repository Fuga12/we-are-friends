<script setup lang="ts">
import { onMounted, ref } from 'vue'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import ToggleSwitch from 'primevue/toggleswitch'
import DeleteButton from '@/components/admin/DeleteButton.vue'
import { deleteAdminClassMedia, deleteAdminSummerShiftMedia, deleteMedia, getAdminClassForm, getAdminSummerShiftForm, saveAdminClassForm, saveAdminSummerShiftForm, uploadMedia, type AdminClassForm, type MediaItem } from '@/api/admin'

const props = withDefaults(defineProps<{ token: string; id: number | null; kind?: 'class' | 'summer' }>(), { kind: 'class' })
const emit = defineEmits<{ saved: []; cancel: [] }>()
const blank = (): AdminClassForm => ({ title:'', ageRange:'', description:'', schedule:'', price:'', emoji:'🎨', color:'#e46129', coverMediaId:null, photoMediaIds:[], fullDescription:[], highlights:[], position:0, isPublished:false })
const form = ref<AdminClassForm>(blank())
const coverMedia = ref<MediaItem | null>(null)
const photoMedia = ref<MediaItem[]>([])
const uploadedIds = new Set<number>()
const originalAttachedIds = new Set<number>()
const loading = ref(false)
const uploading = ref(false)
const error = ref('')
const coverInput = ref<HTMLInputElement | null>(null)
const galleryInput = ref<HTMLInputElement | null>(null)
const paragraphs = ref('')
const highlights = ref('')
const emojiPicker = ref<HTMLDetailsElement | null>(null)
const emojiOptions = '🎨 🖌️ 🖍️ ✏️ ✂️ 📐 🧵 🧶 🪡 🧩 🧱 🪵 🪚 🔨 🛠️ ⚙️ 🤖 💻 🔬 🔭 🧪 🧬 🧮 📚 📖 📝 🔤 🔢 🗣️ 🌍 🏛️ 🏺 🎭 🎬 🎤 🎵 🎹 🎸 🥁 🎻 🎧 💃 🕺 🩰 🤸 🧘 🏃 ⚽ 🏀 🏐 🎾 🏓 🏸 🥋 ⛸️ ♟️ 🎲 🎯 🏆 🍎 🍓 🥕 🍪 🧁 🍰 🥣 👩‍🍳 🌱 🌿 🌳 🌻 🐝 🦋 🐾 ☀️ 🌈 ⭐ ✨ 🚀 ✈️ 🚂 🚗 ⛵ 🏕️ 🏠 ❤️ 😊 🤗 🥳 🧠 💡 🔎 🎓'.split(' ')

function lines(value: string) { return value.split('\n').map(line => line.trim()).filter(Boolean) }
function chooseEmoji(value: string) { form.value.emoji = value; if (emojiPicker.value) emojiPicker.value.open = false }
async function load() {
  loading.value = true
  try {
    if (props.id) {
      form.value = props.kind === 'summer' ? await getAdminSummerShiftForm(props.token, props.id) : await getAdminClassForm(props.token, props.id)
      coverMedia.value = form.value.coverMedia ?? null
      photoMedia.value = form.value.photoMedia ?? []
      if (coverMedia.value) originalAttachedIds.add(coverMedia.value.id)
      photoMedia.value.forEach(item => originalAttachedIds.add(item.id))
    }
    paragraphs.value = form.value.fullDescription.join('\n\n')
    highlights.value = form.value.highlights.join('\n')
  } catch (reason) { error.value = reason instanceof Error ? reason.message : 'Не удалось открыть форму' }
  finally { loading.value = false }
}
async function uploadCover(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0]
  if (!file) return
  uploading.value = true; error.value = ''
  try {
    const item = await uploadMedia(props.token, file, `${form.value.title} — обложка`)
    uploadedIds.add(item.id)
    coverMedia.value = item
    form.value.coverMediaId = item.id
  } catch (reason) { error.value = reason instanceof Error ? reason.message : 'Не удалось загрузить обложку' }
  finally { uploading.value = false; if (coverInput.value) coverInput.value.value = '' }
}
async function uploadGallery(event: Event) {
  const files = Array.from((event.target as HTMLInputElement).files ?? [])
  if (!files.length) return
  uploading.value = true; error.value = ''
  try {
    for (const file of files) {
      const item = await uploadMedia(props.token, file, form.value.title)
      uploadedIds.add(item.id)
      photoMedia.value.push(item)
      form.value.photoMediaIds.push(item.id)
    }
  } catch (reason) { error.value = reason instanceof Error ? reason.message : 'Не удалось загрузить фото' }
  finally { uploading.value = false; if (galleryInput.value) galleryInput.value.value = '' }
}
async function removePhoto(item: MediaItem) {
  loading.value = true; error.value = ''
  try {
    if (props.id && originalAttachedIds.has(item.id)) {
      if (props.kind === 'summer') await deleteAdminSummerShiftMedia(props.token, props.id, item.id)
      else await deleteAdminClassMedia(props.token, props.id, item.id)
      originalAttachedIds.delete(item.id)
    } else {
      await deleteMedia(props.token, item.id)
      uploadedIds.delete(item.id)
    }
    if (coverMedia.value?.id === item.id) { coverMedia.value = null; form.value.coverMediaId = null }
    photoMedia.value = photoMedia.value.filter(photo => photo.id !== item.id)
    form.value.photoMediaIds = form.value.photoMediaIds.filter(id => id !== item.id)
  } catch (reason) { error.value = reason instanceof Error ? reason.message : 'Не удалось удалить фотографию' }
  finally { loading.value = false }
}
async function cancel() {
  await Promise.allSettled(Array.from(uploadedIds).map(id => deleteMedia(props.token, id)))
  emit('cancel')
}
async function save() {
  error.value = ''
  if (!form.value.title.trim()) { error.value = `Не хватает обязательного поля: ${props.kind === 'summer' ? 'название смены' : 'название занятия'}`; return }
  loading.value = true
  try {
    form.value.fullDescription = lines(paragraphs.value)
    form.value.highlights = lines(highlights.value)
    if (props.kind === 'summer') await saveAdminSummerShiftForm(props.token, form.value)
    else await saveAdminClassForm(props.token, form.value)
    const retained = new Set([form.value.coverMediaId, ...form.value.photoMediaIds].filter((id): id is number => id !== null))
    await Promise.allSettled(Array.from(originalAttachedIds).filter(id => !retained.has(id)).map(id => deleteMedia(props.token, id)))
    uploadedIds.clear()
    emit('saved')
  } catch (reason) { error.value = reason instanceof Error ? reason.message : `Не удалось сохранить ${props.kind === 'summer' ? 'смену' : 'занятие'}` }
  finally { loading.value = false }
}
onMounted(load)
</script>

<template>
  <section class="class-form" :aria-busy="loading">
    <div class="form-heading"><div><p class="eyebrow">{{ id ? 'Редактирование' : kind === 'summer' ? 'Новая смена' : 'Новое занятие' }}</p><h2>{{ id ? form.title || (kind === 'summer' ? 'Смена' : 'Занятие') : kind === 'summer' ? 'Создать смену' : 'Создать занятие' }}</h2></div><Button label="Закрыть" text severity="secondary" @click="cancel" /></div>
    <p v-if="error" class="form-error">{{ error }}</p>
    <div class="fields two"><label>Название *<InputText v-model="form.title" placeholder="Например, Столярка" /><small>Единственное обязательное поле.</small></label><label>Возраст<InputText v-model="form.ageRange" placeholder="8–13 лет" /></label></div>
    <label>Краткое описание<Textarea v-model="form.description" rows="4" auto-resize placeholder="Текст для карточки занятия" /></label>
    <div class="fields two"><label>{{ kind === 'summer' ? 'Даты' : 'Расписание' }}<InputText v-model="form.schedule" :placeholder="kind === 'summer' ? '13–24 июля' : 'Чт 17:00, Сб 12:00'" /></label><label>Стоимость<InputText v-model="form.price" placeholder="по запросу" /></label></div>
    <div class="fields compact"><label>Значок<details ref="emojiPicker" class="emoji-picker"><summary>{{ form.emoji || '＋' }}</summary><div class="emoji-grid"><button v-for="emoji in emojiOptions" :key="emoji" type="button" :class="{ selected: form.emoji === emoji }" @click="chooseEmoji(emoji)">{{ emoji }}</button></div></details></label><label>Цвет карточки<div class="color-field"><input v-model="form.color" type="color" /><InputText v-model="form.color" /></div></label></div>
    <label>Подробное описание <small>Один абзац — один блок, разделяйте пустой строкой</small><Textarea v-model="paragraphs" rows="8" auto-resize /></label>
    <label>Что получит ребёнок <small>Каждый пункт с новой строки</small><Textarea v-model="highlights" rows="5" auto-resize /></label>

    <div class="photo-section">
      <div class="form-heading"><div><h3>Обложка</h3><p>Основная фотография для карточки и отдельной страницы.</p></div><Button :label="coverMedia ? 'Заменить обложку' : 'Загрузить обложку'" :loading="uploading" @click="coverInput?.click()" /></div>
      <input ref="coverInput" hidden type="file" accept="image/*" @change="uploadCover" />
      <article v-if="coverMedia" class="cover-photo"><img :src="coverMedia.url" :alt="coverMedia.alt" /><DeleteButton label="Удалить фотографию" @confirm="removePhoto(coverMedia)" /></article>
      <p v-else class="empty-photo">Обложка не загружена. На сайте будет показана заглушка с эмодзи.</p>
    </div>
    <div class="photo-section">
      <div class="form-heading"><div><h3>Галерея {{ kind === 'summer' ? 'смены' : 'занятия' }}</h3><p>Дополнительные фотографии загружаются отдельно от обложки.</p></div><Button label="Добавить фотографии" :loading="uploading" @click="galleryInput?.click()" /></div>
      <input ref="galleryInput" hidden type="file" accept="image/*" multiple @change="uploadGallery" />
      <div v-if="photoMedia.length" class="photo-grid"><article v-for="item in photoMedia" :key="item.id"><img :src="item.url" :alt="item.alt" /><DeleteButton size="small" @confirm="removePhoto(item)" /></article></div>
      <p v-else class="empty-photo">В галерее пока нет фотографий.</p>
    </div>
    <div class="publish-row"><div><strong>Показывать на сайте</strong><small>{{ form.isPublished ? (kind === 'summer' ? 'Смена опубликована' : 'Занятие опубликовано') : 'Сохранится как черновик' }}</small></div><ToggleSwitch v-model="form.isPublished" /></div>
    <div class="actions"><Button :label="id ? 'Сохранить изменения' : kind === 'summer' ? 'Создать смену' : 'Создать занятие'" :loading="loading" @click="save" /><Button label="Отмена" text severity="secondary" @click="cancel" /></div>
  </section>
</template>

<style scoped>
.class-form{display:grid;gap:1.35rem}.form-heading{display:flex;align-items:flex-start;justify-content:space-between;gap:1rem}.form-heading h2,.form-heading h3{margin:0}.form-heading p:not(.eyebrow){margin:.35rem 0 0;color:#647067}.fields{display:grid;gap:1rem}.fields.two{grid-template-columns:1fr 1fr}.fields.compact{grid-template-columns:minmax(12rem,1fr) minmax(0,2fr)}.fields>*{min-width:0}.class-form label{display:grid;gap:.45rem;font-weight:700}.class-form label small{font-weight:400;color:#738078}.emoji-picker{position:relative}.emoji-picker summary{height:2.65rem;box-sizing:border-box;padding:.55rem .8rem;border:1px solid #d9d4ca;border-radius:7px;background:#fff;cursor:pointer;font-size:1.35rem;list-style:none}.emoji-picker summary::-webkit-details-marker{display:none}.emoji-picker summary::after{content:'⌄';float:right;color:#647067;font-size:1rem}.emoji-grid{position:absolute;z-index:20;top:calc(100% + .35rem);left:0;width:min(31rem,75vw);max-height:18rem;overflow:auto;display:grid;grid-template-columns:repeat(10,1fr);gap:.25rem;padding:.65rem;background:#fff;border:1px solid #d9d4ca;border-radius:12px;box-shadow:0 14px 35px #173d3525}.emoji-grid button{display:grid;place-items:center;aspect-ratio:1;border:0;border-radius:7px;background:transparent;font-size:1.35rem;cursor:pointer}.emoji-grid button:hover,.emoji-grid button.selected{background:#e1f1eb}.color-field{display:flex;gap:.5rem;min-width:0}.color-field :deep(.p-inputtext){min-width:0}.color-field input[type=color]{width:3rem;min-width:3rem;height:2.65rem;border:1px solid #d9d4ca;border-radius:7px;background:#fff}.photo-section{border-top:1px solid #e6e2d9;padding-top:1.35rem}.cover-photo{width:min(24rem,100%);margin-top:1rem;border:1px solid #e6e2d9;border-radius:12px;padding:.65rem;display:grid;gap:.45rem}.cover-photo img{width:100%;height:220px;object-fit:cover;border-radius:8px}.photo-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(180px,1fr));gap:1rem;margin-top:1rem}.photo-grid article{border:1px solid #e6e2d9;border-radius:12px;padding:.65rem;display:grid;gap:.45rem}.photo-grid img{width:100%;height:130px;object-fit:cover;border-radius:8px}.empty-photo{padding:1rem;background:#f3f5f1;border-radius:10px;color:#647067}.publish-row{display:flex;justify-content:space-between;align-items:center;padding:1rem;background:#e9f3ef;border-radius:12px}.publish-row div{display:grid;gap:.2rem}.publish-row small{color:#527168}.actions{display:flex;gap:.5rem;position:sticky;bottom:1rem;background:#fff;padding:1rem 0 0;border-top:1px solid #e6e2d9}@media(max-width:700px){.fields.two,.fields.compact{grid-template-columns:1fr}.form-heading{align-items:stretch;flex-direction:column}.emoji-grid{grid-template-columns:repeat(7,1fr)}}
</style>
