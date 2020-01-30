<template>
  <div>
    <h1 class="mb-3">{{taskName}}</h1>
    <div class="mb-3">
      <div v-if="$route.params.platform">
        <strong>Platform:</strong> <span>{{ $route.params.platform }}</span>
      </div>
    </div>
    <b-row>
      <b-col>
        <div class="mb-3">
          <div v-for="(param, key) in params"
            :key="key"
            class="mb-3"
          >
            <p class="mb-1">{{param.label}}</p>
            <input
              v-model="param.value"
              class="mr-2"
            >
          </div>
        </div>
        <b-button
          @click="executeTask"
          class="button mb-3"
          variant="primary"
        >
          <template v-if="executeInProgress">
            <b-spinner label="Spinning"></b-spinner>
          </template>
          <template v-else>
            Execute task
          </template>
        </b-button>
      </b-col>
      <b-col>
        <div
          v-if="success"
          class="success"
        >
          <p>Results: {{success.statusText}}</p>
          <p>{{success.data}}</p>
        </div>
        <div
          v-if="error"
          class="danger"
        >
          {{error}}
        </div>
      </b-col>
    </b-row>
  </div>
</template>

<script lang="ts">
  import { Component, Vue } from 'vue-property-decorator'
  import { AxiosResponse } from 'axios'

  import Api from '@/services/api'

  import { Input } from '@/interfaces/interface'

  @Component
  export default class NewJob extends Vue {
    private executeInProgress: boolean = false
    private params: Input[] = []
    private success: AxiosResponse = { data: '', status: 0, statusText: '', headers: '', config: {} }
    private error: string = ''
    private taskName: string = ''

    /**
     * Methods
     */

    private mounted() {
      this.getTaskConfig()
    }

    private getTaskConfig(): void {
      Api().get(`/task_config?taskName=${this.$route.params.task}`).then(response => {
        this.taskName = response.data.name
        response.data.params.forEach(item => {
          this.params.push({
            label: item.label,
            value: '',
          })
        })
      })
    }

    private executeTask(): void {
      this.resetData()

      const data = {
        task: this.$route.params.task,
        platform: this.$route.params.platform,
        params: this.params.map(param => param.value),
      }

      Api().post('/run_task', data)
      .then((response: AxiosResponse) => {
        this.success = response
      })
      .catch((error: string) => {
        this.error = error
      })
      .finally(() => {
        this.executeInProgress = false
      })
    }

    private resetData(): void {
      this.executeInProgress = true
      this.success = { data: '', status: 0, statusText: '', headers: '', config: {} }
      this.error = ''
    }
  }
</script>

<style lang="scss" scoped>
  .button {
    width: 200px;
  }
</style>
