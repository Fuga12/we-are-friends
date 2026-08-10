<script setup lang="ts">
import { onMounted, ref } from 'vue'
import ToggleSwitch from 'primevue/toggleswitch'
import DeleteButton from '@/components/admin/DeleteButton.vue'
import { deleteAdminLead, getAdminLeads, setAdminLeadProcessed, type AdminLead } from '@/api/admin'

const props = defineProps<{ token:string }>()
const emit = defineEmits<{ changed:[] }>()
const items = ref<AdminLead[]>([])
const loading = ref(false)
const error = ref('')

async function load(){loading.value=true;error.value='';try{items.value=await getAdminLeads(props.token)}catch(e){error.value=e instanceof Error?e.message:'Не удалось загрузить заявки'}finally{loading.value=false}}
async function toggle(item:AdminLead,value:boolean){loading.value=true;try{await setAdminLeadProcessed(props.token,item.id,value);await load();emit('changed')}catch(e){error.value=e instanceof Error?e.message:'Не удалось изменить статус'}finally{loading.value=false}}
async function remove(item:AdminLead){loading.value=true;try{await deleteAdminLead(props.token,item.id);await load();emit('changed')}catch(e){error.value=e instanceof Error?e.message:'Не удалось удалить заявку'}finally{loading.value=false}}
function date(value:string){return new Intl.DateTimeFormat('ru-RU',{dateStyle:'medium',timeStyle:'short'}).format(new Date(value))}
onMounted(load)
</script>

<template>
  <section class="leads">
    <div><p class="eyebrow">Обращения с сайта</p><h2>Заявки на запись</h2><p>Новые заявки находятся сверху. После звонка отметьте заявку обработанной.</p></div>
    <p v-if="error" class="form-error">{{error}}</p>
    <div v-if="items.length" class="lead-list">
      <article v-for="item in items" :key="item.id" :class="{processed:item.isProcessed}">
        <div class="lead-main"><strong>{{item.name}}</strong><a :href="`tel:${item.phone}`">{{item.phone}}</a><small>{{date(item.createdAt)}}</small><p v-if="item.message">{{item.message}}</p></div>
        <div class="lead-actions"><label><span>{{item.isProcessed?'Обработана':'Новая'}}</span><ToggleSwitch :model-value="item.isProcessed" @update:model-value="toggle(item,$event)" /></label><DeleteButton :loading="loading" @confirm="remove(item)" /></div>
      </article>
    </div>
    <p v-else-if="!loading" class="empty">Заявок пока нет.</p>
  </section>
</template>

<style scoped>
.leads{display:grid;gap:1.25rem}.leads h2{margin:0}.leads>div>p:not(.eyebrow){color:#647067}.lead-list{display:grid;gap:.75rem}.lead-list article{display:flex;justify-content:space-between;gap:1.5rem;padding:1rem;border:1px solid #e6e2d9;border-left:5px solid #e46129;border-radius:12px}.lead-list article.processed{opacity:.65;border-left-color:#8b948e}.lead-main{display:grid;gap:.25rem}.lead-main a{font-weight:800;color:#1a6757}.lead-main small{color:#6f766e}.lead-main p{margin:.45rem 0 0;white-space:pre-wrap}.lead-actions{display:flex;align-items:center;gap:1rem}.lead-actions label{display:flex;align-items:center;gap:.5rem}.empty{padding:1rem;background:#f3f5f1;border-radius:10px;color:#647067}@media(max-width:650px){.lead-list article,.lead-actions{align-items:flex-start;flex-direction:column}}
</style>
