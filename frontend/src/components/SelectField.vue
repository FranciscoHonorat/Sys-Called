<script setup lang="ts">
defineProps<{
  id: string
  label: string
  options: ReadonlyArray<{ value: string; label: string }>
}>()

const value = defineModel<string>({ default: '' })
</script>

<template>
  <div>
    <label :for="id" class="mb-1 block text-sm font-medium text-slate-700">{{ label }}</label>
    <select
      :id="id"
      v-model="value"
      class="w-full rounded border border-slate-300 px-3 py-2 focus:border-blue-500 focus:outline-none"
    >
      <option v-for="option in options" :key="option.value" :value="option.value">{{ option.label }}</option>
    </select>
  </div>
</template>
