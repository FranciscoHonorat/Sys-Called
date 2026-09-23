<script setup lang="ts">
import { paths } from '../router/paths'
</script>

<template>
  <div class="space-y-4 text-center">
    <p role="status" class="rounded bg-green-50 p-3 text-sm text-green-800"><slot /></p>
    <RouterLink :to="paths.login" class="text-sm text-blue-600 hover:underline">Voltar para o login</RouterLink>
  </div>
</template>
