<script setup lang="ts">
import { onMounted, ref } from 'vue'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import ToggleSwitch from 'primevue/toggleswitch'
import DeleteButton from '@/components/admin/DeleteButton.vue'
import { deleteAdminNewsMedia, deleteMedia, getAdminNewsCategories, getAdminNewsForm, saveAdminNewsForm, uploadMedia, type AdminNewsForm, type MediaItem } from '@/api/admin'

const props = defineProps<{ token: string; id: number | null }>()
const emit = defineEmits<{ saved: []; cancel: [] }>()
const today = new Date().toISOString().slice(0, 10)
const form = ref<AdminNewsForm>({ title:'', date:today, category:'', excerpt:'', emoji:'📰', coverMediaId:null, photoMediaIds:[], content:[], isPublished:false })
const coverMedia = ref<MediaItem|null>(null)
const photoMedia = ref<MediaItem[]>([])
const originalAttachedIds = new Set<number>()
const uploadedIds = new Set<number>()
const contentText = ref('')
const loading = ref(false)
const uploading = ref(false)
const error = ref('')
const coverInput = ref<HTMLInputElement|null>(null)
const galleryInput = ref<HTMLInputElement|null>(null)
const emojiPicker = ref<HTMLDetailsElement|null>(null)
const emojis = '📰 📣 📢 ✨ ⭐ 🎉 🥳 ❤️ 😊 🤗 👋 🎨 🎭 🎵 🏆 🎓 📚 ✏️ 🧩 🤖 🔬 🧪 🍪 🧁 🌱 🌿 🌳 🌻 ☀️ 🌈 ❄️ 🎄 🎁 🎂 🚀 ✈️ 🏕️ ⚽ 🏀 ♟️ 💡 📸 🎬 🏠'.split(' ')
const categories = ref<string[]>([])
const categorySelection = ref('__new__')

function paragraphs(value:string){return value.split(/\n\s*\n/).map(v=>v.trim()).filter(Boolean)}
function chooseEmoji(value:string){form.value.emoji=value;if(emojiPicker.value)emojiPicker.value.open=false}
function changeCategory(){if(categorySelection.value==='__new__')form.value.category='';else form.value.category=categorySelection.value}
async function load(){loading.value=true;try{categories.value=await getAdminNewsCategories(props.token);if(props.id){form.value=await getAdminNewsForm(props.token,props.id);categorySelection.value=categories.value.includes(form.value.category)?form.value.category:'__new__';coverMedia.value=form.value.coverMedia??null;photoMedia.value=form.value.photoMedia??[];if(coverMedia.value)originalAttachedIds.add(coverMedia.value.id);photoMedia.value.forEach(v=>originalAttachedIds.add(v.id));contentText.value=form.value.content.join('\n\n')}}catch(reason){error.value=reason instanceof Error?reason.message:'Не удалось открыть новость'}finally{loading.value=false}}
async function uploadCover(event:Event){const file=(event.target as HTMLInputElement).files?.[0];if(!file)return;uploading.value=true;error.value='';try{const item=await uploadMedia(props.token,file,`${form.value.title} — обложка`);uploadedIds.add(item.id);coverMedia.value=item;form.value.coverMediaId=item.id}catch(reason){error.value=reason instanceof Error?reason.message:'Не удалось загрузить обложку'}finally{uploading.value=false;if(coverInput.value)coverInput.value.value=''}}
async function uploadGallery(event:Event){const files=Array.from((event.target as HTMLInputElement).files??[]);if(!files.length)return;uploading.value=true;error.value='';try{for(const file of files){const item=await uploadMedia(props.token,file,form.value.title);uploadedIds.add(item.id);photoMedia.value.push(item);form.value.photoMediaIds.push(item.id)}}catch(reason){error.value=reason instanceof Error?reason.message:'Не удалось загрузить фотографии'}finally{uploading.value=false;if(galleryInput.value)galleryInput.value.value=''}}
async function removePhoto(item:MediaItem){loading.value=true;error.value='';try{if(props.id&&originalAttachedIds.has(item.id)){await deleteAdminNewsMedia(props.token,props.id,item.id);originalAttachedIds.delete(item.id)}else{await deleteMedia(props.token,item.id);uploadedIds.delete(item.id)}if(coverMedia.value?.id===item.id){coverMedia.value=null;form.value.coverMediaId=null}photoMedia.value=photoMedia.value.filter(v=>v.id!==item.id);form.value.photoMediaIds=form.value.photoMediaIds.filter(id=>id!==item.id)}catch(reason){error.value=reason instanceof Error?reason.message:'Не удалось удалить фотографию'}finally{loading.value=false}}
async function cancel(){await Promise.allSettled(Array.from(uploadedIds).map(id=>deleteMedia(props.token,id)));emit('cancel')}
async function save(){error.value='';const missing=[!form.value.title.trim()&&'заголовок',!form.value.date&&'дата публикации',!form.value.category.trim()&&'рубрика',!form.value.excerpt.trim()&&'краткий анонс'].filter(Boolean);if(missing.length){error.value=`Не хватает обязательных полей: ${missing.join(', ')}`;return}loading.value=true;try{form.value.content=paragraphs(contentText.value);await saveAdminNewsForm(props.token,form.value);const retained=new Set([form.value.coverMediaId,...form.value.photoMediaIds].filter((id):id is number=>id!==null));await Promise.allSettled(Array.from(originalAttachedIds).filter(id=>!retained.has(id)).map(id=>deleteMedia(props.token,id)));uploadedIds.clear();emit('saved')}catch(reason){error.value=reason instanceof Error?reason.message:'Не удалось сохранить новость'}finally{loading.value=false}}
onMounted(load)
</script>

<template>
  <section class="news-form" :aria-busy="loading">
    <div class="heading"><div><p class="eyebrow">{{ id?'Редактирование':'Новая публикация' }}</p><h2>{{ id?(form.title||'Новость'):'Создать новость' }}</h2></div><Button label="Закрыть" text severity="secondary" @click="cancel" /></div>
    <p v-if="error" class="form-error">{{ error }}</p>
    <label>Заголовок *<InputText v-model="form.title" placeholder="Заголовок новости" /></label>
    <div class="fields"><label>Дата публикации *<input v-model="form.date" class="native-input" type="date" /></label><label>Рубрика *<select v-model="categorySelection" class="native-input" @change="changeCategory"><option v-for="category in categories" :key="category" :value="category">{{ category }}</option><option value="__new__">＋ Добавить новую рубрику</option></select><InputText v-if="categorySelection==='__new__'" v-model="form.category" placeholder="Название новой рубрики" /></label><label>Значок<details ref="emojiPicker" class="emoji-picker"><summary>{{ form.emoji }}</summary><div><button v-for="emoji in emojis" :key="emoji" type="button" @click="chooseEmoji(emoji)">{{ emoji }}</button></div></details></label></div>
    <label>Краткий анонс *<Textarea v-model="form.excerpt" rows="4" auto-resize placeholder="Короткий текст для карточки новости" /></label>
    <label>Текст новости <small>Разделяйте абзацы пустой строкой</small><Textarea v-model="contentText" rows="12" auto-resize /></label>
    <div class="photo-section"><div class="heading"><div><h3>Обложка</h3><p>Основная фотография новости.</p></div><Button :label="coverMedia?'Заменить обложку':'Загрузить обложку'" :loading="uploading" @click="coverInput?.click()" /></div><input ref="coverInput" hidden type="file" accept="image/*" @change="uploadCover" /><article v-if="coverMedia" class="cover"><img :src="coverMedia.url" :alt="coverMedia.alt" /><DeleteButton label="Удалить фотографию" @confirm="removePhoto(coverMedia)" /></article><p v-else class="empty">Обложка не загружена. Новость всё равно будет отображаться — вместо фотографии используется выбранный значок.</p></div>
    <div class="photo-section"><div class="heading"><div><h3>Галерея новости</h3><p>Дополнительные фотографии публикации.</p></div><Button label="Добавить фотографии" :loading="uploading" @click="galleryInput?.click()" /></div><input ref="galleryInput" hidden type="file" accept="image/*" multiple @change="uploadGallery" /><div v-if="photoMedia.length" class="photo-grid"><article v-for="item in photoMedia" :key="item.id"><img :src="item.url" :alt="item.alt" /><DeleteButton @confirm="removePhoto(item)" /></article></div><p v-else class="empty">В галерее пока нет фотографий.</p></div>
    <div class="publish"><div><strong>Показывать на сайте</strong><small>{{ form.isPublished?'Новость опубликована':'Сохранится как черновик' }}</small></div><ToggleSwitch v-model="form.isPublished" /></div>
    <div class="actions"><Button :label="id?'Сохранить изменения':'Создать новость'" :loading="loading" @click="save" /><Button label="Отмена" text severity="secondary" @click="cancel" /></div>
  </section>
</template>

<style scoped>
.news-form{display:grid;gap:1.35rem}.heading{display:flex;justify-content:space-between;align-items:flex-start;gap:1rem}.heading h2,.heading h3{margin:0}.heading p{margin:.35rem 0 0;color:#647067}.news-form label{display:grid;gap:.45rem;font-weight:700}.news-form label small{font-weight:400;color:#738078}.fields{display:grid;grid-template-columns:1fr 1fr 8rem;gap:1rem}.native-input{box-sizing:border-box;width:100%;padding:.7rem;border:1px solid #d9d4ca;border-radius:7px;color:#173d35}.emoji-picker{position:relative}.emoji-picker summary{height:2.65rem;box-sizing:border-box;padding:.55rem .8rem;border:1px solid #d9d4ca;border-radius:7px;background:#fff;cursor:pointer;font-size:1.35rem;list-style:none}.emoji-picker>div{position:absolute;z-index:20;right:0;top:calc(100% + .35rem);width:22rem;display:grid;grid-template-columns:repeat(8,1fr);gap:.2rem;padding:.6rem;background:#fff;border:1px solid #d9d4ca;border-radius:12px;box-shadow:0 14px 35px #173d3525}.emoji-picker button{aspect-ratio:1;border:0;border-radius:6px;background:transparent;font-size:1.25rem;cursor:pointer}.emoji-picker button:hover{background:#e1f1eb}.photo-section{border-top:1px solid #e6e2d9;padding-top:1.35rem}.cover{width:min(24rem,100%);margin-top:1rem;border:1px solid #e6e2d9;border-radius:12px;padding:.65rem;display:grid}.cover img{width:100%;height:220px;object-fit:cover;border-radius:8px}.photo-grid{display:grid;grid-template-columns:repeat(auto-fill,minmax(180px,1fr));gap:1rem;margin-top:1rem}.photo-grid article{border:1px solid #e6e2d9;border-radius:12px;padding:.65rem;display:grid}.photo-grid img{width:100%;height:130px;object-fit:cover;border-radius:8px}.empty{padding:1rem;background:#f3f5f1;border-radius:10px;color:#647067}.publish{display:flex;justify-content:space-between;align-items:center;padding:1rem;background:#e9f3ef;border-radius:12px}.publish div{display:grid;gap:.2rem}.publish small{color:#527168}.actions{display:flex;gap:.5rem;position:sticky;bottom:1rem;background:#fff;padding-top:1rem;border-top:1px solid #e6e2d9}@media(max-width:700px){.heading{flex-direction:column}.fields{grid-template-columns:1fr}.emoji-picker>div{left:0;right:auto;width:min(22rem,75vw)}}
</style>
