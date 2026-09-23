<script setup lang="ts">
import { ref } from 'vue'

import FormField from './FormField.vue'

const props = defineProps<{ title: string; description: string }>()
const emit = defineEmits<{ save: [changes: { title: string; description: string }]; cancel: [] }>()

const draftTitle = ref(props.title)
const draftDescription = ref(props.description)
</script>

<template>
  <form class="space-y-3" @submit.prevent="emit('save', { title: draftTitle, description: draftDescription })">
    <FormField id="edit-title" v-model="draftTitle" label="Título" />
    <FormField id="edit-description" v-model="draftDescription" label="Descrição" multiline />
    <div class="flex gap-2">
      <button type="submit" class="rounded border border-slate-300 px-3 py-2 text-sm text-slate-700 hover:bg-slate-50">
        Salvar
      </button>
      <button
        type="button"
        class="rounded border border-slate-300 px-3 py-2 text-sm text-slate-700 hover:bg-slate-50"
        @click="emit('cancel')"
      >
        Cancelar
      </button>
    </div>
  </form>
</template>
