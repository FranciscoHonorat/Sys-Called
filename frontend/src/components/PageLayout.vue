<script setup lang="ts">
import AppHeader from './AppHeader.vue'
import BackLink from './BackLink.vue'

defineProps<{ back?: boolean; backTo?: string }>()
</script>

<template>
  <main class="min-h-screen bg-slate-100 px-4 py-6 sm:px-6">
    <div class="mx-auto w-full max-w-5xl">
      <AppHeader />
      <BackLink v-if="back" :to="backTo" />
      <slot />
    </div>
  </main>
</template>
