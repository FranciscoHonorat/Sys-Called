<script setup lang="ts">
import { inject, ref } from 'vue'

import AppModal from '../components/AppModal.vue'
import PageLayout from '../components/PageLayout.vue'
import { useLoad } from '../composables/useLoad'
import { employeesApiKey, type Employee } from '../employees/employeesApi'
import { needsAttention, roleLabel, situationLabel } from '../employees/roles'

const api = inject(employeesApiKey)!

const { data: employees, failed, reload } = useLoad<Employee[]>(() => api.list(), [])
const issued = ref<{ name: string; password: string } | null>(null)
const actionFailed = ref(false)

async function run(action: () => Promise<void>) {
  actionFailed.value = false
  try {
    await action()
  } catch {
    actionFailed.value = true
  }
}

function approve(employee: Employee) {
  return run(async () => {
    await api.approve(employee.id)
    await reload()
  })
}

function issueTemporaryPassword(employee: Employee) {
  return run(async () => {
    issued.value = { name: employee.name, password: await api.issueTemporaryPassword(employee.id) }
    await reload()
  })
}

const actionClass = 'rounded border border-blue-600 px-2 py-1 text-xs font-medium text-blue-700 hover:bg-blue-50'
</script>

<template>
  <PageLayout back>
    <section class="rounded-lg bg-white p-6 shadow">
      <h2 class="mb-4 text-xl font-semibold text-slate-800">Usuários</h2>
      <p v-if="failed" role="alert" class="rounded bg-red-50 p-2 text-sm text-red-700">Não foi possível carregar os usuários</p>
      <template v-else>
        <p v-if="actionFailed" role="alert" class="mb-3 rounded bg-red-50 p-2 text-sm text-red-700">
          Não foi possível concluir a ação
        </p>
        <div class="overflow-x-auto">
          <table class="w-full text-left text-sm">
            <thead class="text-slate-500">
              <tr>
                <th class="py-2">Nome</th>
                <th>Username</th>
                <th>Perfil</th>
                <th>Situação</th>
                <th><span class="sr-only">Ações</span></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="employee in employees" :key="employee.id" class="border-t border-slate-100">
                <td class="py-2 font-medium text-slate-800">{{ employee.name }}</td>
                <td>{{ employee.username }}</td>
                <td>{{ roleLabel(employee.role) }}</td>
                <td>
                  <span :class="needsAttention(employee) ? 'font-medium text-amber-700' : 'text-slate-600'">{{
                    situationLabel(employee)
                  }}</span>
                </td>
                <td class="py-2 text-right">
                  <button
                    v-if="employee.status === 'pending'"
                    type="button"
                    :aria-label="`Aprovar ${employee.name}`"
                    :class="actionClass"
                    @click="approve(employee)"
                  >Aprovar</button>
                  <button
                    v-else
                    type="button"
                    :aria-label="`Gerar senha temporária para ${employee.name}`"
                    :class="actionClass"
                    @click="issueTemporaryPassword(employee)"
                  >Gerar senha temporária</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </template>
    </section>

    <AppModal :open="issued !== null" title="Senha temporária" @close="issued = null">
      <div v-if="issued" class="space-y-4">
        <p class="rounded bg-slate-100 p-3 text-center font-mono text-lg tracking-wider text-slate-800">{{ issued.password }}</p>
        <p class="text-sm text-slate-600">
          Entregue esta senha pessoalmente: ela só aparece agora. {{ issued.name }} vai precisar trocá-la no próximo acesso.
        </p>
        <button type="button" class="w-full rounded bg-blue-600 py-2 font-medium text-white hover:bg-blue-700" @click="issued = null">
          Fechar
        </button>
      </div>
    </AppModal>
  </PageLayout>
</template>
