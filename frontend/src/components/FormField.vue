<script setup lang="ts">
defineProps<{
  id: string
  label: string
  type?: 'text' | 'password'
  autocomplete?: string
  error?: string
  multiline?: boolean
}>()

const value = defineModel<string>({ default: '' })

const inputClass = 'w-full rounded border border-slate-300 px-3 py-2 focus:border-blue-500 focus:outline-none'
</script>

<template>
  <div>
    <label :for="id" class="mb-1 block text-sm font-medium text-slate-700">{{ label }}</label>
    <textarea v-if="multiline" :id="id" v-model="value" rows="4" :class="inputClass" />
    <input v-else :id="id" v-model="value" :type="type ?? 'text'" :autocomplete="autocomplete" :class="inputClass" />
    <p v-if="error" class="mt-1 text-sm text-red-600">{{ error }}</p>
  </div>
</template>
