<script lang="ts" setup>
import { useNuxtApp } from "nuxt/app";
import { onMounted } from "vue";
import type { FormError, FormSubmitEvent } from "#ui/types";

const message = ref<string>("");
const { $socket } = useNuxtApp();

onMounted(() => {
  $socket.onopen = () => {
    console.log("open");
  };

  $socket.onmessage = ({ data }: any) => {
    console.log("data", data);

    const unmarshalledData = JSON.parse(data);

    unmarshalledData["streams"].forEach((element) => {
      rows.value.unshift({
        id: element["id"],
        service: element["service_name"],
        text: element["text"].substr(0, 10),
        flag: "",
        start: element["started_at"],
        end: element["ended_at"],
      });
    });

    message.value = data;
  };

  $socket.onclose = function () {
    console.log("disconnected");
  };
});

const Flag = {
  EMPTY: "EMPTY",
  IN: "IN",
  OUT: "OUT",
};

// Columns
const columns = [
  {
    key: "id",
    label: "ID",
    sortable: true,
  },
  {
    key: "service",
    label: "Service",
    sortable: true,
  },
  {
    key: "text",
    label: "Text",
  },
  {
    key: "flag",
    label: "Flag",
    sortable: true,
  },
  {
    key: "start",
    label: "Start",
    sortable: false,
  },
  {
    key: "end",
    label: "End",
    sortable: false,
  },
];

function getFlagBadgeColor(flag) {
  if (flag === Flag.IN) {
    return "emerald";
  } else if (flag === Flag.OUT) {
    return "orange";
  } else {
    return "black";
  }
}

const rows = ref([
  {
    id: 3,
    service: "ШАПИТУЛИ",
    text: "FLAG",
    flag: Flag.EMPTY,
    start: "start",
    end: "end",
  },
  {
    id: 2,
    service: "ЛАПИТУЛИ",
    text: "FLAG",
    flag: Flag.IN,
    start: "start",
    end: "end",
  },
  {
    id: 1,
    service: "ТУТУТУЛИ",
    text: "FLAG",
    flag: Flag.OUT,
    start: "start",
    end: "end",
  },
]);

const selectedColumns = ref(columns);
const columnsTable = computed(() =>
  columns.filter((column) => selectedColumns.value.includes(column))
);

// Selected Rows
const selectedRows = ref([]);

function select(row) {
  const index = selectedRows.value.findIndex((item) => item.id === row.id);
  if (index === -1) {
    selectedRows.value.push(row);
  } else {
    selectedRows.value.splice(index, 1);
  }
}

const search = ref("");
const selectedStatus = ref([]);
const searchStatus = computed(() => {
  if (selectedStatus.value?.length === 0) {
    return "";
  }

  if (selectedStatus?.value?.length > 1) {
    return `?flag=${selectedStatus.value[0].value}&flag=${selectedStatus.value[1].value}`;
  }

  return `?flag=${selectedStatus.value[0].value}`;
});

// Pagination
const page = ref(1);
const pageCount = ref(10);
const pageTotal = ref(200); // This value should be dynamic coming from the API
const pageFrom = computed(() => (page.value - 1) * pageCount.value + 1);
const pageTo = computed(() =>
  Math.min(page.value * pageCount.value, pageTotal.value)
);

// Data
const { data: todos, pending } = await useLazyAsyncData(
  "todos",
  () =>
    $fetch<
      {
        id: number;
        title: string;
        completed: string;
      }[]
    >(`ws://localhost:2137/api/ws/get-streams`),
  {
    default: () => [],
    watch: [page, search, searchStatus, pageCount],
  }
);

const isOpenServices = ref(false);

const validate = (state: any): FormError[] => {
  const errors = [];
  if (!state.name) errors.push({ path: "name", message: "Required" });
  if (!state.port) errors.push({ path: "port", message: "Required" });
  if (!state.flag_regexp)
    errors.push({ path: "flag_regexp", message: "Required" });
  return errors;
};

async function onSubmit(event: FormSubmitEvent<any>) {
  console.log(event.data);
  console.log(services);
}

const services = [
  {
    key: "service1",
    label: "service1",
    port: 1338,
    flag_regexp: "FLAG1",
    state: reactive({
      name: "service1",
      port: 1338,
      flag_regexp: "FLAG1",
    }),
  },
  {
    key: "service2",
    label: "service2",
    port: 1337,
    flag_regexp: "FLAG2",
    state: reactive({
      name: "service2",
      port: 1337,
      flag_regexp: "FLAG2",
    }),
  },
  {
    key: "add",
    label: "Новый сервис",
    port: undefined,
    flag_regexp: undefined,
    state: reactive({
      name: undefined,
      port: undefined,
      flag_regexp: undefined,
    }),
  },
];
</script>

<template>
  <UCard class="w-full min-h-screen min-w-screen">
    <template #header>
      <h2
        class="font-semibold text-xl text-gray-900 dark:text-white leading-tight"
      >
        <div>
          <UButton label="Services" @click="isOpenServices = true" />
          <UModal v-model="isOpenServices">
            <UTabs
              :items="services"
              orientation="vertical"
              :ui="{
                wrapper: 'flex items-center gap-10',
                list: { width: 'w-60' },
              }"
            >
              <template #item="{ item }">
                <UCard @submit.prevent="onSubmit">
                  <template #header>
                    <p
                      class="text-base font-semibold leading-6 text-gray-900 dark:text-white"
                    >
                      {{ item.label }}
                    </p>
                    <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
                      {{ item.flag }}
                    </p>
                  </template>

                  <UForm
                    :validate="validate"
                    :state="item.state"
                    class="space-y-4"
                    @submit="onSubmit"
                  >
                    <UFormGroup label="Name" name="name">
                      <UInput v-model="item.state.name" />
                    </UFormGroup>

                    <UFormGroup label="Port" name="port">
                      <UInput v-model="item.state.port" />
                    </UFormGroup>

                    <UFormGroup label="Flag regexp" name="flag_regexp">
                      <UInput v-model="item.state.flag_regexp" />
                    </UFormGroup>
                  </UForm>

                  <template #footer>
                    <UButton type="submit"> Save </UButton>
                  </template>
                </UCard>
              </template>
            </UTabs>
            <!--               
              <UCard>
                <template #header>
                  <p
                    class="text-base font-semibold leading-6 text-gray-900 dark:text-white"
                  >
                    Service1
                  </p>
                  <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
                    Описание сервиса
                  </p>
                </template>

                <UForm
                  :validate="validate"
                  :state="state"
                  class="space-y-4"
                  @submit="onSubmit"
                >
                  <UFormGroup label="Email" name="email">
                    <UInput v-model="state.email" />
                  </UFormGroup>

                  <UFormGroup label="Password" name="password">
                    <UInput v-model="state.password" type="password" />
                  </UFormGroup>
                </UForm>

                <template #footer>
                  <UButton type="submit"> Save </UButton>
                </template>
              </UCard>

              <UCard>
                <template #header>
                  <p
                    class="text-base font-semibold leading-6 text-gray-900 dark:text-white"
                  >
                    Service1
                  </p>
                  <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">
                    Описание сервиса
                  </p>
                </template>

                <UForm
                  :validate="validate"
                  :state="state"
                  class="space-y-4"
                  @submit="onSubmit"
                >
                  <UFormGroup label="Email" name="email">
                    <UInput v-model="state.email" />
                  </UFormGroup>

                  <UFormGroup label="Password" name="password">
                    <UInput v-model="state.password" type="password" />
                  </UFormGroup>
                </UForm>

                <template #footer>
                  <UButton type="submit"> Save </UButton>
                </template>
              </UCard> -->
          </UModal>
        </div>
      </h2>
    </template>

    <!-- Filters -->
    <div class="flex items-center justify-between gap-3 px-4 py-3">
      <UInput
        v-model="search"
        icon="i-heroicons-magnifying-glass-20-solid"
        placeholder="Search..."
      />
    </div>

    <!-- Header and Action buttons -->
    <div class="flex justify-between items-center w-full px-4 py-3">
      <div class="flex items-center gap-1.5">
        <span class="text-sm leading-5">Rows per page:</span>

        <USelect
          v-model="pageCount"
          :options="[3, 5, 10, 20, 30, 40]"
          class="me-2 w-20"
          size="xs"
        />
      </div>
    </div>

    <!-- Table -->
    <UTable
      v-model="selectedRows"
      :rows="rows"
      :columns="columnsTable"
      :loading="pending"
      sort-asc-icon="i-heroicons-arrow-up"
      sort-desc-icon="i-heroicons-arrow-down"
      class="w-full min-h-screen min-w-screen"
      :ui="{ td: { base: 'max-w-[0] truncate' } }"
      @select="select"
    >
      <template #flag-data="{ row }">
        <UBadge
          size="xs"
          :label="row.flag"
          :color="getFlagBadgeColor(row.flag)"
          variant="subtle"
        />
      </template>

      <template #actions-data="{ row }">
        <UButton
          v-if="!row.flag"
          icon="i-heroicons-check"
          size="2xs"
          color="emerald"
          variant="outline"
          :ui="{ rounded: 'rounded-full' }"
          square
        />

        <UButton
          v-else
          icon="i-heroicons-arrow-path"
          size="2xs"
          color="orange"
          variant="outline"
          :ui="{ rounded: 'rounded-full' }"
          square
        />
      </template>
    </UTable>

    <!-- Number of rows & Pagination -->
    <template #footer>
      <div class="flex flex-wrap justify-between items-center">
        <div>
          <span class="text-sm leading-5">
            Showing
            <span class="font-medium">{{ pageFrom }}</span>
            to
            <span class="font-medium">{{ pageTo }}</span>
            of
            <span class="font-medium">{{ pageTotal }}</span>
            results
          </span>
        </div>

        <UPagination
          v-model="page"
          :page-count="pageCount"
          :total="pageTotal"
          :ui="{
            wrapper: 'flex items-center gap-1',
            rounded: '!rounded-full min-w-[32px] justify-center',
            default: {
              activeButton: {
                variant: 'outline',
              },
            },
          }"
        />
      </div>
    </template>
  </UCard>
</template>
