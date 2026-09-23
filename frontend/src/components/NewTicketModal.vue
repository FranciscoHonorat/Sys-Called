<script setup lang="ts">
import AppModal from './AppModal.vue'
import NewTicketForm from './NewTicketForm.vue'

defineProps<{ open: boolean }>()
const emit = defineEmits<{ close: []; opened: [ticketId: string] }>()
</script>

<template>
  <AppModal :open="open" title="Abrir novo chamado" @close="emit('close')">
    <NewTicketForm @opened="(id) => emit('opened', id)" />
  </AppModal>
</template>
